package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/bragemusic/core/pkg/bragerr"
	"github.com/bragemusic/core/pkg/client"
	"github.com/bragemusic/core/pkg/config"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/musicbrainz"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/internal/assethandler"
	"github.com/bragemusic/core/internal/utils"
	"github.com/gofrs/uuid/v5"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	ctxCancelFunc context.CancelFunc
	wg            sync.WaitGroup
	client        client.ClientFace
	authClient    client.Auth
	slogHandler   slog.Handler
	clientConfig  config.ClientConfig
	log           *slog.Logger
}

// NewApp creates a new App application struct
func New(authClient client.Auth, clientConfig config.ClientConfig, slogHandler slog.Handler) *App {
	a := &App{
		client:       nil,
		authClient:   authClient,
		log:          slog.New(slogHandler),
		clientConfig: clientConfig,
		slogHandler:  slogHandler,
	}

	// a.client.RegisterPlaybackStateCallback(callbackFuncerCtx[types.PlaybackState](a, types.ClientEventPlayerPlaybackChange))
	// a.client.RegisterPlayContextCallback(callbackFuncerCtx[types.PlayContext](a, types.ClientEventPlayerContextChange))
	// // a.client.RegisterServerAvailabilityCallback(callbackFuncer[types.ServerApiInfo](a, types.ClientEventServerOnline))
	// a.client.RegisterSyncInProgressCallback(callbackFuncer[bool](a, types.ClientEventSyncInProgress))
	// // a.client.RegisterUserCallback(a.userCallback)
	// a.client.RegisterEventCallback(a.sendEvent)

	// // a.client.SubscribeToEventTypes(a.sendSSEvent, types.SSEventTypeClientConnected, types.SSEventTypeClientDisconnected)
	// // a.client.SubscribeToEventTypes(a.sendRemoteDevicePlayerState, types.SSEventTypePlayerPlayContext, types.SSEventTypePlayerPlaybackState)

	// a.client.SubscribeToClientEvents(a.sendSSEvent)

	return a
}

func (a *App) createClient() (err error) {
	a.client, err = a.authClient.NewClient(a.ctx, a.clientConfig.ClientConfig(), a.slogHandler)
	if err != nil {
		return err
	}

	a.client.RegisterPlaybackStateCallback(callbackFuncerCtx[types.PlaybackState](a, types.ClientEventPlayerPlaybackChange))
	a.client.RegisterPlayContextCallback(callbackFuncerCtx[types.PlayContext](a, types.ClientEventPlayerContextChange))
	a.client.RegisterSyncInProgressCallback(callbackFuncer[bool](a, types.ClientEventSyncInProgress))
	a.client.RegisterServerAvailabilityCallback(callbackFuncer[types.ServerApiInfo](a, types.ClientEventServerOnline))
	a.client.RegisterEventCallback(a.sendEvent)
	a.client.SubscribeToClientEvents(a.sendSSEvent)

	go func() {
		a.client.StartScheduler(a.ctx)
	}()

	tokenType := types.TokenFrontendShort
	if a.clientConfig.General.ClientType == types.DeviceTypeSync {
		tokenType = types.TokenAPI
	}

	a.sendEvent(types.ClientEventAuthLoggedIn, true)
	a.sendEvent(types.ClientEventUserUpdated, a.authClient.GetUser(a.ctx, tokenType))

	return nil
}

func (a *App) destroyClient() (err error) {
	if a.client == nil {
		return errors.New("no client to destroy")
	}

	err = a.client.Close(a.ctx)
	if err != nil {
		return err
	}

	a.sendEvent(types.ClientEventAuthLoggedIn, false)
	a.sendEvent(types.ClientEventUserUpdated, nil)
	a.client = nil

	return nil
}

func (a *App) startDevHandler() {
	f := assethandler.New(a.clientConfig, a.authClient.ServerClient())
	if err := http.ListenAndServe(":31145", &f); err != nil {
		a.log.Error("dev handler crashed", "error", err.Error())
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx, a.ctxCancelFunc = context.WithCancel(ctx)

	err := utils.SetupFolderStructure(ctx, a.clientConfig)
	if err != nil {
		a.log.ErrorContext(ctx, "could not create folder structure", "error", err.Error())
		a.ctxCancelFunc()
		return
	}

	if runtime.Environment(ctx).BuildType == "dev" {
		go a.startDevHandler()
	}

	tokenType := types.TokenFrontendShort
	if a.clientConfig.General.ClientType == types.DeviceTypeSync {
		tokenType = types.TokenAPI
	}

	user := a.authClient.GetUser(a.ctx, tokenType)
	if user != nil {
		err := a.createClient()
		if err != nil {
			a.log.ErrorContext(ctx, "could not create client", "error", err.Error())
			a.ctxCancelFunc()
			return
		}
	}
	// if err := a.client.LoadCachedUserAndToken(ctx, uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111")); err != nil {
	// 	a.log.ErrorContext(ctx, "token not loaded", "error", err.Error())
	// }

	// a.cachedUsers, err = a.client.GetCachedUsers(ctx)
	// if err != nil {
	// 	a.log.ErrorContext(ctx, "could not load cached users", "error", err.Error())
	// }

	// a.wg.Add(1)
	// a.client.StartSyncDaemon(a.ctx, a.wg.Done)

	// go func() {
	// 	a.client.StartScheduler(ctx)
	// }()

	// a.wg.Add(1)
	// a.client.StartStatusDaemon(a.ctx, a.wg.Done)
}

func (a *App) Shutdown(ctx context.Context) {
	a.log.InfoContext(ctx, "shutdown started")
	a.ctxCancelFunc()

	a.log.InfoContext(ctx, "waiting for workers")
	a.wg.Wait()

	a.log.InfoContext(ctx, "all workers finished. Quitting app")
}

func (a *App) Themes() []types.ThemeDescription {
	td := []types.ThemeDescription{}

	td = append(td, types.ThemeDescription{ID: "light", Name: "Brage Light"})
	td = append(td, types.ThemeDescription{ID: "dark", Name: "Brage Dark"})

	for id, theme := range a.clientConfig.Themes {
		td = append(td, types.ThemeDescription{
			ID:   id,
			Name: theme.Name,
		})
	}

	return td
}

func (a *App) userCallback(user *types.UserDetails) {
	// a.user = user
	// if user != nil {
	// 	a.ctx = auth.UpgradeContextWithUser(a.ctx, *user)
	// }
	a.sendEvent(types.ClientEventUserUpdated, user)
}

func (a *App) SendMessage(message string) {
	runtime.EventsEmit(a.ctx, "message", message)
}

func (a *App) sendEvent(e types.ClientEvent, payload any) {
	runtime.EventsEmit(a.ctx, string(e), payload)
}

func (a *App) sendSSEvent(ctx context.Context, e types.SSEvent) {
	runtime.EventsEmit(a.ctx, string(e.Type), e.Data)
}

func (a *App) sendRemoteDevicePlayerState(ctx context.Context, e types.SSEvent) {
	runtime.EventsEmit(a.ctx, string(e.Type)+".remote", e.Data)
}

func (a *App) SendMessage3(isPlaying bool) {
	runtime.EventsEmit(a.ctx, "isPlaying", isPlaying)
}

func (a *App) SendMessage4(progressMs int64) {
	runtime.EventsEmit(a.ctx, "progress", progressMs)
}

func (a *App) GetCachedUsers() []types.UserDetails {
	cu, err := a.authClient.GetCachedUsers(a.ctx)
	if err != nil {
		a.handleError(err)
		return []types.UserDetails{}
	}
	return cu
}

// func (a *App) SelectLocalUser(userID string) {
// 	uID, err := uuid.FromString(userID)
// 	if err != nil {
// 		a.handleError(err)
// 		return
// 	}

// 	err = a.authClient.LoginLocalUser(a.ctx, uID)
// 	if err != nil {
// 		a.handleError(err)
// 		return
// 	}

// 	err = a.createClient()
// 	if err != nil {
// 		a.handleError(err)
// 		return
// 	}
// }

func (a *App) LogoutLocalUser() {
	err := a.authClient.Logout(a.ctx)
	if err != nil {
		a.handleError(err)
		return
	}
	err = a.destroyClient()
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) ListUserTokens() []types.TokenLimited {
	tokens, err := a.client.ListUserTokens(a.ctx)
	if err != nil {
		a.handleError(err)
		return []types.TokenLimited{}
	}

	return tokens
}

func (a *App) DeleteUserToken(id string) {
	uid, err := uuid.FromString(id)
	if err != nil {
		a.handleError(err)
		return
	}

	err = a.client.DeleteUserToken(a.ctx, uid)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) CreateAPIToken(name string) string {
	token, err := a.client.CreateAPIToken(a.ctx, name)
	if err != nil {
		a.handleError(err)
		return ""
	}

	return token
}

func (a *App) ListArtists() []types.ArtistDetailed {
	artists, err := a.client.ListArtists(a.ctx, database.SortByName, database.SortAsc)
	if err != nil {
		a.handleError(err)
		return nil
	}
	return artists
}

func (a *App) CreateArtist(artist types.ArtistBase) {
	err := a.client.CreateArtist(a.ctx, artist)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) GetArtist(artistID string) types.Artist {
	uID, err := uuid.FromString(artistID)
	if err != nil {
		a.handleError(err)
		return types.Artist{}
	}

	artist, err := a.client.GetArtist(a.ctx, uID)
	if err != nil {
		a.handleError(err)
		return types.Artist{}
	}

	return artist
}

func (a *App) ListAlbumsByArtist(artistID string) []types.AlbumDetailed {
	uID, err := uuid.FromString(artistID)
	if err != nil {
		a.handleError(err)
		return []types.AlbumDetailed{}
	}

	albums, err := a.client.ListAlbumsByArtist(a.ctx, uID, database.SortByDate, database.SortAsc)
	if err != nil {
		a.handleError(err)
		return nil
	}

	if albums == nil {
		return []types.AlbumDetailed{}
	}

	return albums
}

func (a *App) ListFeaturedAlbumsByArtist(artistID string) []types.AlbumDetailed {
	uID, err := uuid.FromString(artistID)
	if err != nil {
		a.handleError(err)
		return []types.AlbumDetailed{}
	}

	albums, err := a.client.ListFeaturedAlbumsByArtist(a.ctx, uID, database.SortByDate, database.SortAsc)
	if err != nil {
		a.handleError(err)
		return nil
	}

	if albums == nil {
		return []types.AlbumDetailed{}
	}

	return albums
}

func (a *App) ListAlbums() []types.AlbumDetailed {
	albums, err := a.client.ListAlbums(a.ctx, database.SortByName, database.SortAsc)
	if err != nil {
		a.handleError(err)
		return nil
	}
	return albums
}

func (a *App) GetAlbum(albumID string) types.AlbumDetailed {
	uID, err := uuid.FromString(albumID)
	if err != nil {
		a.handleError(err)
		return types.AlbumDetailed{}
	}

	album, err := a.client.GetAlbumDetailed(a.ctx, uID)
	if err != nil {
		a.handleError(err)
		return types.AlbumDetailed{}
	}

	return album
}

func (a *App) ListTracksByAlbum(albumID string) []types.TrackDetailed {
	uID, err := uuid.FromString(albumID)
	if err != nil {
		a.handleError(err)
		return []types.TrackDetailed{}
	}

	tracks, err := a.client.ListTracksDetailedByAlbum(a.ctx, uID)
	if err != nil {
		a.handleError(err)
		return nil
	}

	return tracks
}

func (a *App) UploadArtistImage(artistID string, img serverclient.FileUpload) {
	uID, err := uuid.FromString(artistID)
	if err != nil {
		a.handleError(err)
		return
	}

	err = a.client.UploadArtistImage(a.ctx, uID, img)
	if err != nil {
		a.handleError(err)
		return
	}

	a.sendEvent(types.ClientEventMsgSuccess, "Artist Image Uploaded")
}

func (a *App) GetArtistTopTracks(artistID string) []types.TrackDetailed {
	id, err := uuid.FromString(artistID)
	if err != nil {
		a.handleError(err)
		return []types.TrackDetailed{}
	}

	tracks, err := a.client.GetArtistTopTracks(a.ctx, id)
	if err != nil {
		a.handleError(err)
		return []types.TrackDetailed{}
	}

	if tracks == nil {
		return []types.TrackDetailed{}
	}

	return tracks
}

func (a *App) UpdateArtist(artistID uuid.UUID, artistData types.Artist) {
	err := a.client.UpdateArtist(a.ctx, artistID, artistData)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) CountArtists() int {
	count, err := a.client.CountArtists(a.ctx)
	if err != nil {
		a.handleError(err)
		return 0
	}
	return count
}

func (a *App) UpdateAlbum(id string, album types.AlbumUpdate) {
	uid, err := uuid.FromString(id)
	if err != nil {
		a.handleError(err)
		return
	}

	err = a.client.UpdateAlbum(a.ctx, uid, album)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) CountAlbums() int {
	count, err := a.client.CountAlbums(a.ctx)
	if err != nil {
		a.handleError(err)
		return 0
	}
	return count
}

func (a *App) UpdateTrack(id string, track types.TrackUpdate) {
	uid, err := uuid.FromString(id)
	if err != nil {
		a.handleError(err)
		return
	}

	err = a.client.UpdateTrack(a.ctx, uid, track)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) RateTrack(trackID string, value int) {
	uid, err := uuid.FromString(trackID)
	if err != nil {
		a.handleError(err)
		return
	}

	err = a.client.RateTrack(a.ctx, uid, value)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) LikeTrack(trackID string) {
	uid, err := uuid.FromString(trackID)
	if err != nil {
		a.handleError(err)
		return
	}

	err = a.client.LikeTrack(a.ctx, uid)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) UnlikeTrack(trackID string) {
	uid, err := uuid.FromString(trackID)
	if err != nil {
		a.handleError(err)
		return
	}

	err = a.client.UnlikeTrack(a.ctx, uid)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) CountLikedTracks() int {
	count, err := a.client.CountLikedTracks(a.ctx)
	if err != nil {
		a.handleError(err)
		return 0
	}
	return count
}

func (a *App) FilterTracks(filter types.TrackFilter, page, limit int) types.ListPaginationPayload[types.TrackDetailed] {
	results, err := a.client.FilterTracks(a.ctx, filter, page, limit)
	if err != nil {
		a.handleError(err)
		return types.ListPaginationPayload[types.TrackDetailed]{}
	}

	return results
}

func (a *App) ListLikedTracks() []types.TrackDetailed {
	tracks, err := a.client.ListLikedTracks(a.ctx)
	if err != nil {
		a.handleError(err)
		return []types.TrackDetailed{}
	}

	if tracks == nil {
		return []types.TrackDetailed{}
	}

	return tracks
}

func (a *App) CountTracks() int {
	count, err := a.client.CountTracks(a.ctx)
	if err != nil {
		a.handleError(err)
		return 0
	}
	return count
}

func (a *App) UploadAlbumImage(id string, img serverclient.FileUpload) {
	uid, err := uuid.FromString(id)
	if err != nil {
		a.handleError(err)
		return
	}

	err = a.client.UploadAlbumImage(a.ctx, uid, img)
	if err != nil {
		a.handleError(err)
		return
	}

	a.sendEvent(types.ClientEventMsgSuccess, "Album Image Uploaded")
}

func (a *App) GetPlayerState() types.PlayerState {
	return a.client.PlayerState()
}

func (a *App) GetUser() *types.UserDetails {
	tokenType := types.TokenFrontendShort
	if a.clientConfig.General.ClientType == types.DeviceTypeSync {
		tokenType = types.TokenAPI
	}

	user := a.authClient.GetUser(a.ctx, tokenType)
	a.sendEvent(types.ClientEventUserUpdated, user)
	return user
}

func (a *App) handleError(err error) {
	if err == nil {
		return
	}

	title := "Unknown Error"
	code := "UNKNOWN"
	msg := "unknown error"
	status := 0
	service := "unknown"
	logErr := err

	berr, ok := err.(*bragerr.BragErr)
	if ok {
		title = berr.Title
		msg = berr.Message
		code = berr.Code
		status = berr.Status
		service = berr.Service
		if berr.Err != nil {
			logErr = berr.Err
		}
	}

	if status == 0 {
		serr, ok := err.(serverclient.ErrStatus)
		if ok {
			status = serr.Status
		}
	}

	a.log.ErrorContext(a.ctx, msg, "code", code, "service", service, "error", logErr.Error())

	a.sendEvent(types.ClientEventMsgErr, Message{
		Title:   title,
		Message: msg,
	})

	if status == http.StatusUnauthorized && code != bragerr.ErrInvalidUserCreds.Code {
		err = a.destroyClient()
		if err != nil {
			a.log.ErrorContext(a.ctx, "could not destroy client", "error", err.Error())
			return
		}
	}
}

func (a *App) SupportsSync() bool {
	if a.clientConfig.ClientConfig().ClientType == types.DeviceTypeSync {
		return true
	} else {
		return false
	}
}

func (a *App) SyncLibrary() {
	err := a.sync()
	if err != nil {
		return
	}
	a.sendEvent(types.ClientEventMsgSuccess, "Sync finished")
}

func (a *App) sync() error {
	err := a.client.Sync(a.ctx)
	if err != nil {
		a.handleError(err)
		return err
	}
	return nil
}

func (a *App) LoginServerUser(username, password string, longLivedToken bool) {
	tokenType := types.TokenFrontendShort

	if a.clientConfig.General.ClientType == types.DeviceTypeSync {
		tokenType = types.TokenAPI
	} else if longLivedToken {
		tokenType = types.TokenFrontendLong
	}

	err := a.authClient.LoginServerUser(a.ctx, username, password, tokenType)
	if err != nil {
		a.handleError(err)
		return
	}

	err = a.createClient()
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) ServerStatus() types.ServerApiInfo {
	if a.client == nil {
		status, err := a.authClient.ServerStatus(a.ctx)
		if err != nil {
			a.handleError(err)
			return types.ServerApiInfo{}
		}
		return status
	}

	status, err := a.client.ServerStatus(a.ctx)
	if err != nil {
		a.handleError(err)
		return types.ServerApiInfo{}
	}
	return status
}

func (a *App) AddPlaylist(playlist types.Playlist) {
	err := a.client.AddPlaylist(a.ctx, playlist)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) AddSmartPlaylist(playlist types.PlaylistBase, filter types.TrackFilter) {
	err := a.client.AddSmartPlaylist(a.ctx, playlist, filter)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) AddPlaylistTrack(playlistID, albumID, trackID string) {
	pID, err := uuid.FromString(playlistID)
	if err != nil {
		a.handleError(err)
		return
	}

	aID, err := uuid.FromString(albumID)
	if err != nil {
		a.handleError(err)
		return
	}

	tID, err := uuid.FromString(trackID)
	if err != nil {
		a.handleError(err)
		return
	}

	err = a.client.AddPlaylistTrack(a.ctx, pID, aID, tID)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) CountPlaylists() int {
	count, err := a.client.CountPlaylists(a.ctx)
	if err != nil {
		a.handleError(err)
		return 0
	}
	return count
}

func (a *App) CountPlaylistTracks(playlistID string) int {
	pID, err := uuid.FromString(playlistID)
	if err != nil {
		a.handleError(err)
		return 0
	}

	count, err := a.client.CountPlaylistTracks(a.ctx, pID)
	if err != nil {
		a.handleError(err)
		return 0
	}
	return count
}

func (a *App) DeletePlaylist(id string) {
	uid, err := uuid.FromString(id)
	if err != nil {
		a.handleError(err)
		return
	}

	err = a.client.DeletePlaylist(a.ctx, uid)
	if err != nil {
		a.handleError(err)
		return
	}

	a.sendEvent(types.ClientEventMsgSuccess, "Playlist removed")
}

func (a *App) DeletePlaylistTrack(id string) {
	uid, err := uuid.FromString(id)
	if err != nil {
		a.handleError(err)
		return
	}

	err = a.client.DeletePlaylistTrack(a.ctx, uid)
	if err != nil {
		a.handleError(err)
		return
	}

	a.sendEvent(types.ClientEventMsgSuccess, "Track removed from playlist")
}

func (a *App) GetPlaylist(id string) types.Playlist {
	uid, err := uuid.FromString(id)
	if err != nil {
		a.handleError(err)
		return types.Playlist{}
	}

	plist, err := a.client.GetPlaylist(a.ctx, uid)
	if err != nil {
		a.handleError(err)
		return types.Playlist{}
	}

	return plist
}

func (a *App) ListPlaylists() []types.Playlist {
	plists, err := a.client.ListPlaylists(a.ctx, false, database.SortByDate, database.SortAsc)
	if err != nil {
		a.handleError(err)
		return []types.Playlist{}
	}

	return plists
}

func (a *App) ListPlaylistTracks(playlistID string) []types.TrackDetailed {
	pUID, err := uuid.FromString(playlistID)
	if err != nil {
		a.handleError(err)
		return []types.TrackDetailed{}
	}

	tracks, err := a.client.ListPlaylistTracks(a.ctx, pUID, database.SortByDate, database.SortAsc)
	if err != nil {
		a.handleError(err)
		return []types.TrackDetailed{}
	}

	if tracks == nil {
		return []types.TrackDetailed{}
	}

	return tracks
}

func (a *App) UpdatePlaylist(id string, data types.Playlist) {
	uid, err := uuid.FromString(id)
	if err != nil {
		a.handleError(err)
		return
	}

	err = a.client.UpdatePlaylist(a.ctx, uid, data)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) UploadPlaylistImage(id string, img serverclient.FileUpload) {
	uid, err := uuid.FromString(id)
	if err != nil {
		a.handleError(err)
		return
	}

	err = a.client.UploadPlaylistImage(a.ctx, uid, img)
	if err != nil {
		a.handleError(err)
		return
	}

	a.sendEvent(types.ClientEventMsgSuccess, "Playlist Image Uploaded")
}

func (a *App) ListEntityEvents() []types.EntityEvent {
	events, err := a.client.ListEntityEvents(a.ctx)
	if err != nil {
		a.handleError(err)
		return []types.EntityEvent{}
	}

	return events
}

func (a *App) SearchFull(searchTerm string) []types.SearchItem {
	res, err := a.client.SearchFull(a.ctx, searchTerm)
	if err != nil {
		a.handleError(err)
		return nil
	}

	if res == nil {
		return []types.SearchItem{}
	}

	return res
}

func (a *App) ImportAlbum(musicbrainzID *string) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select album zip",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Zip files",
				Pattern:     "*.zip",
			},
		},
	})
	if err != nil {
		a.handleError(err)
	}

	err = a.client.ImportAlbum(a.ctx, path, musicbrainzID)
	if err != nil {
		a.handleError(err)
		return
	}

	a.sendEvent(types.ClientEventMsgSuccess, "Album sent to server for processing")
}

func (a *App) ListImportItems(page, limit int) types.ListPaginationPayload[types.Import] {
	items, err := a.client.ListImportItems(a.ctx, page, limit)
	if err != nil {
		a.handleError(err)
		return types.ListPaginationPayload[types.Import]{}
	}

	return items
}

func (a *App) SearchMusicBrainz(artist, album string) types.ListPayload[musicbrainz.SearchResults] {
	res, err := a.client.SearchMusicBrainz(a.ctx, artist, album)
	if err != nil {
		a.handleError(err)
		return types.ListPayload[musicbrainz.SearchResults]{}
	}

	return res
}

func callbackFuncer[T any](a *App, e types.ClientEvent) func(T) {
	return func(payload T) {
		a.sendEvent(e, payload)
	}
}

func callbackFuncerCtx[T any](a *App, e types.ClientEvent) func(context.Context, T) {
	return func(ctx context.Context, payload T) {
		a.sendEvent(e, payload)
	}
}

// func callbackSSEFuncer[T any](a *App, e types.SSEventType) func(context.Context, T) {
// 	return func(ctx context.Context, payload T) {
// 		a.sendSSEvent(e, payload)
// 	}
// }
