package server

import (
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type ReqNoContent struct{}

func (r ReqNoContent) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqArtistsGet struct {
	ArtistID uuid.UUID `path:"artistID" description:"ID of the wanted artist"`
}

func (r ReqArtistsGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqArtistsUpdate struct {
	ArtistID uuid.UUID `path:"artistID" description:"ID of the wanted artist"`
	types.ArtistBase
}

func (r ReqArtistsUpdate) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqAlbumsGet struct {
	AlbumID uuid.UUID `path:"albumID" description:"ID of the wanted album"`
}

func (r ReqAlbumsGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqAlbumsTrackGet struct {
	AlbumID uuid.UUID `path:"albumID" description:"ID of the wanted album"`
	TrackID uuid.UUID `path:"trackID" description:"ID of the wanted track"`
}

func (r ReqAlbumsTrackGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqAlbumsAlbumArtistGet struct {
	AlbumID  uuid.UUID        `path:"albumID" description:"ID of the wanted album"`
	ArtistID uuid.UUID        `path:"artistID" description:"ID of the wanted artist"`
	Role     types.ArtistRole `path:"role" description:"Role of the wanted artist"`
}

func (r ReqAlbumsAlbumArtistGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqAlbumsAlbumTrackGet struct {
	AlbumID     uuid.UUID `path:"albumID" description:"ID of the wanted album track"`
	DiscNumber  int       `path:"discNumber" description:"Disc number of the wanted album track"`
	TrackNumber int       `path:"trackNumber" description:"Track number of the wanted album track"`
}

func (r ReqAlbumsAlbumTrackGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqAlbumsUpdate struct {
	AlbumID uuid.UUID `path:"albumID" description:"ID of the wanted album"`
	types.AlbumUpdate
}

func (r ReqAlbumsUpdate) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqAlbumArtistsGet struct {
	AlbumArtistID uuid.UUID `path:"albumArtistID" description:"ID of the wanted album artist"`
}

func (r ReqAlbumArtistsGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqAlbumTracksGet struct {
	AlbumTrackID uuid.UUID `path:"albumTrackID" description:"ID of the wanted album track"`
}

func (r ReqAlbumTracksGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqTracksGet struct {
	TrackID uuid.UUID `path:"trackID" description:"ID of the wanted track"`
}

func (r ReqTracksGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqTracksUpdate struct {
	TrackID uuid.UUID `path:"trackID" description:"ID of the wanted track"`
	types.TrackUpdate
}

func (r ReqTracksUpdate) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqTracksAddRating struct {
	TrackID uuid.UUID `path:"trackID" description:"ID of the wanted track"`
	types.RatingReq
}

func (r ReqTracksAddRating) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqMediaFilesGet struct {
	MediafileID uuid.UUID `path:"mediafileID" description:"ID of the wanted media file"`
}

func (r ReqMediaFilesGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqMediaFilesGetFile struct {
	MediafileID uuid.UUID `path:"mediafileID" description:"ID of the wanted media file"`
	Range       string    `header:"Range" description:"Optional HTTP Range header for partial content retrieval." example:"bytes=0-1023"`
}

func (r ReqMediaFilesGetFile) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqPlaylistsAdd struct {
	types.PlaylistBase
}

func (r ReqPlaylistsAdd) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqPlaylistsGet struct {
	PlaylistID uuid.UUID `path:"playlistID" description:"ID of the wanted playlist"`
}

func (r ReqPlaylistsGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqPlaylistsUpdate struct {
	PlaylistID uuid.UUID `path:"playlistID" description:"ID of the wanted playlist"`
	types.PlaylistBase
}

func (r ReqPlaylistsUpdate) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqPlaylistsAddTrack struct {
	PlaylistID uuid.UUID `path:"playlistID" description:"ID of the wanted playlist"`
	types.PlaylistTrackReq
}

func (r ReqPlaylistsAddTrack) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqPlaylistTracksGet struct {
	PlaylistTrackID uuid.UUID `path:"playlistTrackID" description:"ID of the wanted playlist track"`
}

func (r ReqPlaylistTracksGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqSyncSync struct {
	types.SyncReq
}

func (r ReqSyncSync) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqSyncPlayHistory struct {
	types.SyncPlayHistoryReq
}

func (r ReqSyncPlayHistory) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqRatingsGet struct {
	RatingID uuid.UUID `path:"ratingID" description:"ID of the wanted rating"`
}

func (r ReqRatingsGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqSearch struct {
	SearchTerm string `query:"q" description:"Free text search string"`
}

func (r ReqSearch) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqImportAlbum struct {
	File     []byte            `form:"file" required:"true" description:"Zip file content" json:"-"`
	Metadata types.ImportAlbum `form:"metadata" required:"true" description:"Additional metadata" json:"-"`
}

func (r ReqImportAlbum) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqLikesGet struct {
	LikeID uuid.UUID `path:"likeID" description:"ID of the wanted like"`
}

func (r ReqLikesGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqDevicesRegister struct {
	types.ReqDevicesRegister
}

func (r ReqDevicesRegister) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqDevicesGet struct {
	DeviceID uuid.UUID `path:"deviceID" description:"ID of the wanted device"`
}

func (r ReqDevicesGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqDevicesPlayerSetRepeat struct {
	DeviceID uuid.UUID `path:"deviceID" description:"ID of the wanted device"`
	types.SetDevicePlayerRepeatReq
}

func (r ReqDevicesPlayerSetRepeat) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqDevicesPlayerSetShuffle struct {
	DeviceID uuid.UUID `path:"deviceID" description:"ID of the wanted device"`
	types.SetDevicePlayerShuffleReq
}

func (r ReqDevicesPlayerSetShuffle) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqDevicesUpdatePlayContext struct {
	DeviceID uuid.UUID `path:"deviceID" description:"ID of the wanted device"`
	types.PlayContextDTO
}

func (r ReqDevicesUpdatePlayContext) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqDevicesUpdatePlaybackState struct {
	DeviceID uuid.UUID `path:"deviceID" description:"ID of the wanted device"`
	types.PlaybackStateDTO
}

func (r ReqDevicesUpdatePlaybackState) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqDevicesPlayerSetState struct {
	DeviceID uuid.UUID `path:"deviceID" description:"ID of the wanted device"`
	types.PlayerState
}

func (r ReqDevicesPlayerSetState) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqDevicesPlayerAddToQueue struct {
	DeviceID uuid.UUID `path:"deviceID" description:"ID of the wanted device"`
	types.TrackDetailed
}

func (r ReqDevicesPlayerAddToQueue) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqList struct {
	Count     bool               `query:"count" description:"Only return the count, not the payload."`
	SortOrder database.SortOrder `query:"sortOrder" description:"Sort ascending or descending."`
	SortBy    database.SortBy    `query:"sortBy" description:"Sort by key."`
}

func (r ReqList) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqListTracksOfAlbum struct {
	AlbumID uuid.UUID `path:"albumID" description:"ID of the wanted album"`
	Count   bool      `query:"count" description:"Only return the count, not the payload."`
}

func (r ReqListTracksOfAlbum) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqListPlaylists struct {
	Count         bool               `query:"count" description:"Only return the count, not the payload."`
	SortOrder     database.SortOrder `query:"sortOrder" description:"Sort ascending or descending."`
	SortBy        database.SortBy    `query:"sortBy" description:"Sort by key."`
	IncludePublic bool               `query:"includePublic" description:"Include public playlists"`
}

func (r ReqListPlaylists) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqListPlaylistTracks struct {
	PlaylistID uuid.UUID          `path:"playlistID" description:"ID of the wanted playlist"`
	Count      bool               `query:"count" description:"Only return the count, not the payload."`
	SortOrder  database.SortOrder `query:"sortOrder" description:"Sort ascending or descending."`
	SortBy     database.SortBy    `query:"sortBy" description:"Sort by key."`
}

func (r ReqListPlaylistTracks) Validate() (validationMessages string, err error) {
	return "", nil
}
