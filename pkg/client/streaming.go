package client

import (
	"context"
	"errors"
	"os"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (c *ClientStreaming) RegisterEventCallback(f func(types.ClientEvent, any)) {
	c.eventCallbacks = append(c.eventCallbacks, f)
}

func (c *ClientStreaming) emitEvent(event types.ClientEvent, payload any) {
	for _, f := range c.eventCallbacks {
		f(event, payload)
	}
}

func (c *ClientStreaming) setUser(user *types.UserDetails) {
	c.user = user
}

func (c *ClientStreaming) AddTrackToQueue(ctx context.Context, trackID, albumID uuid.UUID) error {
	return nil
}

func (c *ClientStreaming) StartPlayerWithAlbum(ctx context.Context, albumID uuid.UUID, trackNumber int) error {
	return nil
}

func (c *ClientStreaming) StartPlayerWithPlaylist(ctx context.Context, playlistID uuid.UUID, trackNumber int, sortBy database.SortBy, sortOrder database.SortOrder) error {
	return nil
}

func (c ClientStreaming) GetUser() *types.UserDetails {
	return c.user
}

func (c *ClientStreaming) Login(ctx context.Context, username, password string, longLivedToken bool) (types.UserDetails, error) {
	return c.AuthClient.Login(ctx, username, password, longLivedToken)
}

func (c *ClientStreaming) LoginLocalUser(ctx context.Context, userID uuid.UUID) error {
	c.log.InfoContext(ctx, "logging in local user", "id", userID.String())

	user, err := c.AuthClient.LoginLocalUser(ctx, userID, false)
	if err != nil {
		return err
	}

	c.AuthClient.UserCallback(&user)
	c.ServerStatus(ctx)

	return nil
}

func (c *ClientStreaming) LoginCachedServerUser(ctx context.Context, password string, longLivedToken bool) error {
	if c.user == nil {
		return c.berr.NoUserInContext(errors.New("could not login cached user"))
	}

	return c.AuthClient.LoginCachedServerUser(ctx, *c.user, password, longLivedToken)
}

func (c *ClientStreaming) LogoutLocalUser(ctx context.Context) {
	c.AuthClient.UserCallback(nil)
}

func (c *ClientStreaming) LogoutServerUser(ctx context.Context) error {
	if c.user == nil {
		return c.berr.NoUserInContext(errors.New("could not logout server user"))
	}
	return c.AuthClient.LogoutServerUser(ctx, c.user.ID)
}

func (c ClientStreaming) ImportAlbum(ctx context.Context, filename string, musicbrainzID *string) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()
	return c.ServerClient.ImportAlbum(ctx, r, filename, musicbrainzID)
}

///////

func (c ClientStreaming) UpdateArtist(ctx context.Context, artistID uuid.UUID, artistData types.Artist) error {
	if err := c.ServerClient.UpdateArtist(ctx, artistID, artistData); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c ClientStreaming) UploadArtistImage(ctx context.Context, artistID uuid.UUID, img serverclient.FileUpload) error {
	if err := c.ServerClient.UploadArtistImage(ctx, artistID, img); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c ClientStreaming) UpdateAlbum(ctx context.Context, id uuid.UUID, album types.AlbumUpdate) error {
	if err := c.ServerClient.UpdateAlbum(ctx, id, album); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c ClientStreaming) UploadAlbumImage(ctx context.Context, id uuid.UUID, img serverclient.FileUpload) error {
	if err := c.ServerClient.UploadAlbumImage(ctx, id, img); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c ClientStreaming) RateTrack(ctx context.Context, trackID uuid.UUID, value int) error {
	if err := c.ServerClient.RateTrack(ctx, trackID, value); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c ClientStreaming) UpdateTrack(ctx context.Context, id uuid.UUID, track types.TrackUpdate) error {
	if err := c.ServerClient.UpdateTrack(ctx, id, track); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c ClientStreaming) AddPlaylist(ctx context.Context, playlist types.Playlist) error {
	if err := c.ServerClient.AddPlaylist(ctx, playlist); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c ClientStreaming) AddPlaylistTrack(ctx context.Context, playlistID, albumID, trackID uuid.UUID) error {
	if err := c.ServerClient.AddPlaylistTrack(ctx, playlistID, albumID, trackID); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c ClientStreaming) DeletePlaylist(ctx context.Context, id uuid.UUID) error {
	if err := c.ServerClient.DeletePlaylist(ctx, id); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c ClientStreaming) DeletePlaylistTrack(ctx context.Context, id uuid.UUID) error {
	if err := c.ServerClient.DeletePlaylistTrack(ctx, id); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c ClientStreaming) UpdatePlaylist(ctx context.Context, id uuid.UUID, data types.Playlist) error {
	if err := c.ServerClient.UpdatePlaylist(ctx, id, data); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c ClientStreaming) UploadPlaylistImage(ctx context.Context, id uuid.UUID, img serverclient.FileUpload) error {
	if err := c.ServerClient.UploadPlaylistImage(ctx, id, img); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}
