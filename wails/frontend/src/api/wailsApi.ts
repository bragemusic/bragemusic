import { Api, UserDetails, responses } from "bragemusic-frontend";
import { types, musicbrainz } from "../../wailsjs/go/models";
import * as App from "../../wailsjs/go/app/App";
import { EventsOn } from "../../wailsjs/runtime/runtime";

export class WailsApi implements Api {
    eventSubscribe = EventsOn;

    async getUser(): Promise<UserDetails> {
        const user = await App.GetUser()
        return new UserDetails(user)
    }

    async loginServerUser(email:string,password:string,rememberMe:boolean):Promise<void> {
        return await App.LoginServerUser(email, password, rememberMe)
    }

    async logoutLocalUser():Promise<void> {
        return await App.LogoutLocalUser()
    }

    async publicServerStatus():Promise<types.ServerApiInfo> {
        return await App.ServerStatus()
    }

    async serverStatus():Promise<types.ServerApiInfo> {
        return await App.ServerStatus()
    }


    async getAlbum(id:string):Promise<types.AlbumDetailed> {
       return await App.GetAlbum(id)
    }

    async listAlbums():Promise<Array<types.AlbumDetailed>> {
        return await App.ListAlbums()
    }

    async listTracksByAlbum(albumID:string):Promise<Array<types.TrackDetailed>> {
        return await App.ListTracksByAlbum(albumID)
    }

    async updateAlbum(id:string, data:types.AlbumUpdate):Promise<void> {
        return await App.UpdateAlbum(id, data)
    }

    async countAlbums():Promise<number> {
        return await App.CountAlbums()
    }


    async likeTrack(id:string):Promise<void> {
        return await App.LikeTrack(id)
    }

    async listLikedTracks():Promise<Array<types.TrackDetailed>> {
        return await App.ListLikedTracks()
    }

    async unlikeTrack(id:string):Promise<void> {
        return await App.UnlikeTrack(id)
    }

    async updateTrack(id:string, data:types.TrackUpdate):Promise<void> {
        return await App.UpdateTrack(id, data)
    }

    async countLikedTracks():Promise<number> {
        return await App.CountLikedTracks()
    }

    async countTracks():Promise<number> {
        return await App.CountTracks()
    }

    async filterTracks(filter: types.TrackFilter, page:number, limit:number):Promise<responses.ListPaginationPayload<types.TrackDetailed>> {
        const raw = await App.FilterTracks(filter, page, limit);
        return raw as unknown as responses.ListPaginationPayload<types.TrackDetailed>;
    }


    async createArtist(artist:types.ArtistBase):Promise<void> {
        return await App.CreateArtist(artist)
    }

    async getArtist(id:string):Promise<types.Artist> {
        return await App.GetArtist(id)
    }

    async getArtistTopTracks(id:string):Promise<Array<types.TrackDetailed>> {
        return await App.GetArtistTopTracks(id)
    }

    async listAlbumsByArtist(id:string):Promise<Array<types.AlbumDetailed>> {
        return await App.ListAlbumsByArtist(id)
    }

    async listFeaturedAlbumsByArtist(id:string):Promise<Array<types.AlbumDetailed>> {
        return await App.ListFeaturedAlbumsByArtist(id)
    }

    async listArtists():Promise<Array<types.ArtistDetailed>> {
        return await App.ListArtists()
    }

    async updateArtist(id:string, data:types.Artist):Promise<void> {
        return await App.UpdateArtist(id, data)
    }

    async countArtists():Promise<number> {
        return await App.CountArtists()
    }


    async addPlaylist(playlist:types.Playlist):Promise<void> {
        return await App.AddPlaylist(playlist)
    }

    async addSmartPlaylist(playlist: types.PlaylistBase, filter: types.TrackFilter): Promise<void> {
        return await App.AddSmartPlaylist(playlist, filter)
    }

    async countPlaylists():Promise<number> {
        return await App.CountPlaylists()
    }

    async deletePlaylist(id:string):Promise<void> {
        return await App.DeletePlaylist(id)
    }

    async getPlaylist(id:string):Promise<types.Playlist> {
        return await App.GetPlaylist(id)
    }

    async listPlaylists():Promise<Array<types.Playlist>> {
        return await App.ListPlaylists()
    }

    async updatePlaylist(id:string, data:types.Playlist):Promise<void> {
        return await App.UpdatePlaylist(id, data)
    }

    async uploadPlaylistImage(
      id: string,
      file: File,
    ): Promise<void> {
      const buffer = await file.arrayBuffer();

      await App.UploadPlaylistImage(id, {
        name: file.name,
        type: file.type,
        data: Array.from(new Uint8Array(buffer)),
      });
      return
    }


    async addPlaylistTrack(playlistID:string, albumID:string, trackID:string):Promise<void> {
        return await App.AddPlaylistTrack(playlistID, albumID, trackID)
    }

    async countPlaylistTracks(playlistID:string):Promise<number> {
        return await App.CountPlaylistTracks(playlistID)
    }

    async deletePlaylistTrack(id:string):Promise<void> {
        return await App.DeletePlaylistTrack(id)
    }

    async listPlaylistTracks(playlistID:string):Promise<Array<types.TrackDetailed>> {
        return await App.ListPlaylistTracks(playlistID)
    }


    async SearchFull(searchTerm:string):Promise<Array<types.SearchItem>> {
        return await App.SearchFull(searchTerm)
    }


    async importAlbum(_file: File, musicbrainzID: string | null): Promise<void> {
        await App.ImportAlbum(musicbrainzID);
        return
    }

    async listImportItems(page:number, limit:number):Promise<responses.ListPaginationPayload<types.Import>> {
        const raw = await App.ListImportItems(page, limit);
        return raw as unknown as responses.ListPaginationPayload<types.Import>;
    }

    async searchMusicBrainz(artist:string, album:string):Promise<responses.ListPayload<musicbrainz.SearchResults>> {
        const raw = await App.SearchMusicBrainz(artist, album);
        return raw as unknown as responses.ListPayload<musicbrainz.SearchResults>;
    }


    async rateTrack(trackID:string, value:number):Promise<void> {
        return await App.RateTrack(trackID, value)
    }


    async listEntityEvents():Promise<Array<types.EntityEvent>> {
        return await App.ListEntityEvents()
    }


    async createUser(email:string,username:string,password:string,roles:Array<string>):Promise<void> {
        return await App.CreateUser(email, username, password, roles)
    }

    async deleteUser(id:string):Promise<void> {
       return await App.DeleteUser(id)
    }

    async listUsers(includeMachineUsers: boolean): Promise<Array<UserDetails>> {
        const users = await App.ListUsers(includeMachineUsers)
        return users.map(u => new UserDetails(u))
    }

    async listUserRoles(): Promise<Array<string>> {
        return await App.ListUserRoles()
    }

    async listUserTokens(): Promise<Array<types.TokenLimited>> {
        return await App.ListUserTokens()
    }

    async deleteToken(id: string): Promise<void> {
        return await App.DeleteUserToken(id)
    }

    async addApiToken(name: string): Promise<responses.CreateApiToken> {
        const token = await App.CreateAPIToken(name)
        return new responses.CreateApiToken({token: token})
    }

    async updateUser(id:string,email:string,username:string,password:any,roles:Array<string>):Promise<void> {
        return await App.UpdateUser(id, email, username, password, roles)
    }

    async updateProfile(
        email?: string,
        username?: string,
        password?: string,
        new_password?: string,
        new_password_confirm?: string,
    ): Promise<void> {
        return await App.UpdateProfile(email, username, password, new_password, new_password_confirm)
    }

    async uploadUserImage(): Promise<void> {
        return await App.UploadUserImage()
    }

    async supportsSync():Promise<boolean> {
        return await App.SupportsSync()
    }

    async syncLibrary():Promise<void> {
        return await App.SyncLibrary()
    }

    async connectDevice(id:string):Promise<void> {
        return await App.ConnectDevice(id)
    }

    async disconnectDevice():Promise<void> {
        return await App.DisconnectDevice()
    }

    async getConnectedDeviceID(): Promise<string|null> {
        return await App.GetConnectedDeviceID()
    }

    async listDevices():Promise<Array<types.DeviceDetailed>> {
        return await App.ListDevices()
    }

    async removeDevice(id:string):Promise<void> {
        return await App.RemoveDevice(id)
    }

    async removeDeviceToken(id:string):Promise<void> {
        return await App.RemoveDeviceToken(id)
    }

    async removeDeviceAndToken(id:string):Promise<void> {
        return await App.RemoveDeviceAndToken(id)
    }


    async listThemes():Promise<Array<types.ThemeDescription>> {
        return await App.Themes()
    }
}
