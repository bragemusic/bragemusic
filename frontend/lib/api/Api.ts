import { UserDetails } from "@/models/UserDetails";
import { musicbrainz, types } from "@/types/core";
import { Event } from "@/types/events.ts";
import { responses } from "@/types/server";

export interface Api {
  eventSubscribe(eventName: Event, callback: (...data: any) => void): () => void;

  // getCachedUsers(): Promise<Array<types.UserDetails>>;
  // selectLocalUser(userID: string): Promise<void>;
  getUser(): Promise<UserDetails>
  loginServerUser(
    email: string,
    password: string,
    rememberMe: boolean,
  ): Promise<void>;
  logoutLocalUser(): Promise<void>;
  publicServerStatus(): Promise<types.ServerApiInfo>;
  serverStatus(): Promise<types.ServerApiInfo>;

  createArtist(artist:types.ArtistBase):Promise<void>;
  getArtist(id:string):Promise<types.Artist>;
  getArtistTopTracks(id:string):Promise<Array<types.TrackDetailed>>;
  listAlbumsByArtist(id:string):Promise<Array<types.AlbumDetailed>>;
  listFeaturedAlbumsByArtist(id:string):Promise<Array<types.AlbumDetailed>>;
  updateArtist(id:string, data:types.Artist):Promise<void>;

  getAlbum(id: string): Promise<types.AlbumDetailed>;
  listAlbums(): Promise<Array<types.AlbumDetailed>>;
  listTracksByAlbum(albumID: string): Promise<Array<types.TrackDetailed>>;
  updateAlbum(id: string, data: types.AlbumUpdate): Promise<void>;
  countAlbums(): Promise<number>;

  filterTracks(filter: types.TrackFilter, page:number, limit:number):Promise<responses.ListPaginationPayload<types.TrackDetailed>>
  likeTrack(id: string): Promise<void>;
  listLikedTracks(): Promise<Array<types.TrackDetailed>>;
  unlikeTrack(id: string): Promise<void>;
  updateTrack(id: string, data: types.TrackUpdate): Promise<void>;
  countLikedTracks(): Promise<number>;
  countTracks(): Promise<number>;
  addPlayCount(trackID: string): Promise<void>;

  listArtists(): Promise<Array<types.ArtistDetailed>>;
  countArtists(): Promise<number>;

  addPlaylist(playlist: types.Playlist): Promise<void>;
  addSmartPlaylist(playlist: types.PlaylistBase, filter: types.TrackFilter): Promise<void>;
  countPlaylists(): Promise<number>;
  deletePlaylist(id: string): Promise<void>;
  getPlaylist(id: string): Promise<types.Playlist>;
  listPlaylists(): Promise<Array<types.Playlist>>;
  updatePlaylist(id: string, data: types.Playlist): Promise<void>;
  uploadPlaylistImage(id: string, file: File): Promise<void>;

  addPlaylistTrack(
    playlistID: string,
    albumID: string,
    trackID: string,
  ): Promise<void>;
  countPlaylistTracks(playlistID: string): Promise<number>;
  deletePlaylistTrack(id: string): Promise<void>;
  listPlaylistTracks(playlistID: string): Promise<Array<types.TrackDetailed>>;

  uploadAlbumImage(id: string, file: File): Promise<void>;
  uploadArtistImage(id: string, file: File): Promise<void>;

  SearchFull(searchTerm: string): Promise<Array<types.SearchItem>>;

  importAlbum(file: File, musicbrainzID: string | null): Promise<void>;
  listImportItems(page:number, limit:number):Promise<responses.ListPaginationPayload<types.Import>>
  searchMusicBrainz(artist:string, album:string):Promise<responses.ListPayload<musicbrainz.SearchResults>>

  rateTrack(trackID: string, value: number): Promise<void>;

  listEntityEvents(): Promise<Array<types.EntityEvent>>;

  createUser(
    email: string,
    username: string,
    password: string,
    roles: Array<string>,
  ): Promise<void>;
  deleteUser(id: string): Promise<void>;
  listUsers(includeMachineUsers: boolean): Promise<Array<UserDetails>>;
  listUserRoles(): Promise<Array<string>>;
  listUserTokens(): Promise<Array<types.TokenLimited>>;
  addApiToken(name: string): Promise<responses.CreateApiToken>;
  deleteToken(id: string): Promise<void>;
  updateUser(
    id: string,
    email: string,
    username: string,
    password: any,
    roles: Array<string>,
  ): Promise<void>;

  uploadUserImage(): Promise<void>;
  updateProfile(
    email?: string,
    username?: string,
    password?: string,
    new_password?: string,
    new_password_confirm?: string,
  ): Promise<void>;

  supportsSync(): Promise<boolean>;
  syncLibrary():Promise<void>;

  getConnectedDeviceID(): Promise<string|null>
  connectDevice(id: string): Promise<void>;
  disconnectDevice(): Promise<void>;
  listDevices(): Promise<Array<types.DeviceDetailed>>;
  removeDevice(id: string): Promise<void>;
  removeDeviceToken(id: string): Promise<void>;
  removeDeviceAndToken(id: string): Promise<void>;

  listThemes(): Promise<Array<types.ThemeDescription>>;


  emitEvent(eventName: Event, data?: any): void;
  //
}
