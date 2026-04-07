package client

import (
	"context"
	"log/slog"
	"os"

	"github.com/bragemusic/core/pkg/bragerr"
	"github.com/bragemusic/core/pkg/jobmanager"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/syncer"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (c *clientStreaming) RegisterEventCallback(f func(types.ClientEvent, any)) {
	c.eventCallbacks = append(c.eventCallbacks, f)
}

func (c *clientStreaming) emitEvent(event types.ClientEvent, payload any) {
	for _, f := range c.eventCallbacks {
		f(event, payload)
	}
}

func (c clientStreaming) ImportAlbum(ctx context.Context, filename string, musicbrainzID *string) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()
	return c.ServerClient.ImportAlbum(ctx, r, filename, musicbrainzID)
}

///////

func (c clientStreaming) UpdateArtist(ctx context.Context, artistID uuid.UUID, artistData types.Artist) error {
	if err := c.ServerClient.UpdateArtist(ctx, artistID, artistData); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c clientStreaming) UploadArtistImage(ctx context.Context, artistID uuid.UUID, img serverclient.FileUpload) error {
	if err := c.ServerClient.UploadArtistImage(ctx, artistID, img); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c clientStreaming) UpdateAlbum(ctx context.Context, id uuid.UUID, album types.AlbumUpdate) error {
	if err := c.ServerClient.UpdateAlbum(ctx, id, album); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c clientStreaming) UploadAlbumImage(ctx context.Context, id uuid.UUID, img serverclient.FileUpload) error {
	if err := c.ServerClient.UploadAlbumImage(ctx, id, img); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c clientStreaming) RateTrack(ctx context.Context, trackID uuid.UUID, value int) error {
	if err := c.ServerClient.RateTrack(ctx, trackID, value); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c clientStreaming) LikeTrack(ctx context.Context, trackID uuid.UUID) error {
	if err := c.ServerClient.LikeTrack(ctx, trackID); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c clientStreaming) UnlikeTrack(ctx context.Context, trackID uuid.UUID) error {
	if err := c.ServerClient.UnlikeTrack(ctx, trackID); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c clientStreaming) UpdateTrack(ctx context.Context, id uuid.UUID, track types.TrackUpdate) error {
	if err := c.ServerClient.UpdateTrack(ctx, id, track); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c clientStreaming) AddPlaylist(ctx context.Context, playlist types.Playlist) error {
	if err := c.ServerClient.AddPlaylist(ctx, playlist); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c clientStreaming) AddPlaylistTrack(ctx context.Context, playlistID, albumID, trackID uuid.UUID) error {
	if err := c.ServerClient.AddPlaylistTrack(ctx, playlistID, albumID, trackID); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c clientStreaming) DeletePlaylist(ctx context.Context, id uuid.UUID) error {
	if err := c.ServerClient.DeletePlaylist(ctx, id); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c clientStreaming) DeletePlaylistTrack(ctx context.Context, id uuid.UUID) error {
	if err := c.ServerClient.DeletePlaylistTrack(ctx, id); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c clientStreaming) UpdatePlaylist(ctx context.Context, id uuid.UUID, data types.Playlist) error {
	if err := c.ServerClient.UpdatePlaylist(ctx, id, data); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c clientStreaming) UploadPlaylistImage(ctx context.Context, id uuid.UUID, img serverclient.FileUpload) error {
	if err := c.ServerClient.UploadPlaylistImage(ctx, id, img); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func newStreamingClient(ctx context.Context, config Config, jm *jobmanager.JobManager, sc *serverclient.ServerClient, user types.UserDetails, slogHandler slog.Handler) (clientFace, error) {
	c := &clientStreaming{
		config:       config,
		log:          slog.New(slogHandler).With("service", "client"),
		JobManager:   jm,
		ServerClient: sc,
		user:         user,
		NoSync:       &syncer.NoSync{},
		berr:         bragerr.NewFactory("client"),
	}

	// FIXME: Live somewhere else
	// jm.RegisterJob(ctx, jobmanager.JobDefinition{
	// 	Type:     types.JobAuthClientServerStatus,
	// 	CronExpr: "*/10 * * * * *",
	// 	Run:      c.AuthClient.UpdateServerStatus,
	// })

	return c, nil
}
