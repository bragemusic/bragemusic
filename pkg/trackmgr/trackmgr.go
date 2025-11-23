package trackmgr

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bragemusic/core/internal/config"
	"github.com/bragemusic/core/pkg/acoustid"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/files"
	"github.com/bragemusic/core/pkg/musicbrainz"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/bragemusic/core/pkg/wiki"
	"github.com/dhowden/tag"
)

type TrackManager struct {
	db            database.Database
	mb            musicbrainz.MusicBrainz
	aid           acoustid.AcoustID
	wiki          wiki.Wiki
	musicDir      string
	artistArtsDir string
	albumArtsDir  string
	log           *slog.Logger
}

func (t TrackManager) AddTrack(ctx context.Context, f *os.File) error {
	metadata, err := tag.ReadFrom(f)
	if err != nil {
		return err
	}

	audioFile, err := files.ParseAudioFile(f, metadata.FileType())
	if err != nil {
		return err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	tx, err := t.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	aIdMatches, err := t.aid.GetMusicBrainzAlbumID(f.Name())
	if err != nil {
		return err
	}

	album, newAlbum, err := t.getOrCreateAlbum(ctx, tx, aIdMatches, metadata)
	if err != nil {
		return err
	}

	artist, newArtist, err := t.getOrCreateArtist(ctx, tx, album, metadata)
	if err != nil {
		return err
	}

	if newArtist {
		artist.ID, err = tx.AddArtist(ctx, artist)
		if err != nil {
			return err
		}
	}

	if newAlbum {
		album.ArtistID = artist.ID
		album.ID, err = tx.AddAlbum(ctx, album)
		if err != nil {
			return err
		}
	}

	tracks, newTracks, err := t.generateTracks(ctx, tx, album, metadata)
	if err != nil {
		return err
	}

	if newTracks {
		ids, err := tx.AddTracks(ctx, tracks)
		if err != nil {
			return err
		}

		for idx := range tracks {
			tracks[idx].ID = ids[idx]
		}
	}

	var track types.Track
	if len(aIdMatches) > 0 && album.MusicBrainzID != nil {
		trackID, err := t.getTrackIDfromAcData(ctx, aIdMatches, *album.MusicBrainzID)
		if err != nil {
			return err
		}
		track, err = tx.GetTrackFromMbID(ctx, trackID)
		if err != nil {
			return err
		}
	} else {
		if len(tracks) > 1 {
			return errors.New("Not supported with multiple non-musicbrainz tracks")
		}
		track = tracks[0]
	}

	folderPath := t.generateAlbumFolderPath(artist.Name, album.Name, t.musicDir)
	filename, err := t.generateTrackPath(*track.DiscNumber, *track.TrackNumber, track.Title, metadata.FileType(), folderPath)
	if err != nil {
		return err
	}

	track = t.updateTrackData(track, audioFile, filename, metadata.FileType())

	err = tx.UpdateTrack(ctx, track)
	if err != nil {
		return err
	}

	// FIXME: Should either download or extract album cover
	// FIXME: Should populate a new artist with metadata and download pics from wiki. Same function can be used for if statement below
	if !newArtist && newAlbum && album.MusicBrainzID != nil && artist.MusicBrainzID == nil {
		mbAlbum, err := t.mb.GetAlbum(ctx, *album.MusicBrainzID)
		if err != nil {
			return err
		}

		if len(mbAlbum.ArtistCredit) > 0 {
			mbArtist, err := t.mb.GetArtist(ctx, mbAlbum.ArtistCredit[0].Artist.ID)
			if err != nil {
				return err
			}

			nArtist := t.generateArtist(mbArtist)
			nArtist.ID = artist.ID

			if nArtist.MusicBrainzID != nil {
				wikiData, err := t.GetArtistMetaData(ctx, *nArtist.MusicBrainzID)
				if err != nil {
					return err
				}
				nArtist.Description = wikiData.Summary

				if wikiData.ImageUrl != nil {
					imgFilename := filepath.Join(t.artistArtsDir, artist.ID+".jpg")
					if err = t.wiki.DownloadFile(ctx, *wikiData.ImageUrl, imgFilename); err != nil {
						return err
					}
					t.log.InfoContext(ctx, "downloaded artist art", "artist", artist.Name)
				}
			}

			err = tx.UpdateArtist(ctx, nArtist)
			if err != nil {
				return err
			}
			t.log.InfoContext(ctx, "artist updated using MusicBrainz ID", "name", artist.Name)
		}
	}

	// NOTE: Now it can happen that an old artist does not get metadata. Maybe
	if newArtist {
		if artist.MusicBrainzID != nil {
			wikiData, err := t.GetArtistMetaData(ctx, *artist.MusicBrainzID)
			if err != nil {
				return err
			}
			artist.Description = wikiData.Summary

			if wikiData.ImageUrl != nil {
				imgFilename := filepath.Join(t.artistArtsDir, artist.ID+".jpg")
				if err = t.wiki.DownloadFile(ctx, *wikiData.ImageUrl, imgFilename); err != nil {
					return err
				}
				t.log.InfoContext(ctx, "downloaded artist art", "artist", artist.Name)
			}

			if err = tx.UpdateArtist(ctx, artist); err != nil {
				return err
			}

		}
	}

	if newAlbum {
		if album.MusicBrainzID != nil {
			t.log.InfoContext(ctx, "downloading album cover from MusicBrainz", "album", album.Name)
			if err = t.mb.DownloadCoverArt(ctx, *album.MusicBrainzID, album.ID, t.albumArtsDir); err != nil {
				t.log.InfoContext(ctx, "could not get album cover from MusicBrainz. Trying from ID3", "album", album.Name)
				imgFilename := filepath.Join(t.albumArtsDir, fmt.Sprintf("%s.%s", album.ID, metadata.Picture().Ext))
				if err = utils.SaveID3Image(ctx, *metadata.Picture(), imgFilename); err != nil {
					t.log.WarnContext(ctx, "could not get image from ID3", "error", err.Error())
				}
			}
		} else if metadata.Picture() != nil {
			imgFilename := filepath.Join(t.albumArtsDir, fmt.Sprintf("%s.%s", album.ID, metadata.Picture().Ext))
			if err = utils.SaveID3Image(ctx, *metadata.Picture(), imgFilename); err != nil {
				t.log.WarnContext(ctx, "could not get image from ID3", "error", err.Error())
			}
		}
	}

	if err = os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return err
	}

	// NOTE: Thinking that we are just overwriting if the file exists
	// if _, err = os.Stat(filename); !errors.Is(err, os.ErrNotExist) {
	// 	return fmt.Errorf("file '%s' already exists", filename)
	// }
	// return tx.Commit()
	// return errors.New("aja")

	dst, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, f); err != nil {
		return err
	}

	return tx.Commit()

	tx.Commit()
	return errors.New("jojo")
}

func (t TrackManager) getTrackIDfromAcData(ctx context.Context, acData []acoustid.AcoustMatch, albumID string) (string, error) {
	for _, ad := range acData {
		if ad.AlbumID == albumID {
			return ad.TrackID, nil
		}
	}
	return "", fmt.Errorf("could not find AcousticID track id for album '%s'", albumID)
}

func (t TrackManager) generateAlbumFolderPath(artist, album, musicDir string) string {
	artist = strings.ReplaceAll(artist, " ", "_")
	album = strings.ReplaceAll(album, " ", "_")

	return path.Join(musicDir, artist, album)
}

func (t TrackManager) generateTrackPath(discNumber, trackNumber int, trackTitle string, format tag.FileType, albumFolder string) (string, error) {
	if format == "" {
		return "", fmt.Errorf("unknown fileformat for track '%s'", trackTitle)
	}

	trackTitle = strings.ReplaceAll(trackTitle, " ", "_")
	trackTitle = strings.ReplaceAll(trackTitle, "/", "_")

	filename := fmt.Sprintf("%02d-%02d-%s.%s", discNumber, trackNumber, trackTitle, strings.ToLower(string(format)))

	return filepath.Join(albumFolder, filename), nil
}

func (t TrackManager) filterAcoustIdMatches(matches []acoustid.AcoustMatch, artistName, albumName string) []acoustid.AcoustMatch {
	fm := []acoustid.AcoustMatch{}
	for _, m := range matches {
		artistMatch := utils.CompareTwoStrings(strings.ToLower(m.ArtistName), strings.ToLower(artistName))
		albumMatch := utils.CompareTwoStrings(strings.ToLower(m.AlbumName), strings.ToLower(albumName))
		if artistMatch > 0.94 && albumMatch > 0.94 {
			fm = append(fm, m)
		}
	}
	return fm
}

func New(cfg config.ServerConfig, db database.Database, aid acoustid.AcoustID, w wiki.Wiki, slogHandler slog.Handler) TrackManager {
	return TrackManager{
		db:            db,
		musicDir:      cfg.MusicDirPath,
		artistArtsDir: cfg.ArtistArtsDirPath,
		albumArtsDir:  cfg.AlbumArtsDirPath,
		aid:           aid,
		wiki:          w,
		log:           slog.New(slogHandler),
	}
}
