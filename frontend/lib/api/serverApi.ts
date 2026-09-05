import { Api } from "./Api";
import { types, serverclient, musicbrainz } from "@/types/core";
import { Event, SSEvent } from "@/types/events";
import { UserDetails } from "@/models/UserDetails";
import ky, { HTTPError } from "ky";
import { requests, responses } from "@/types/server";
import { isBragErr } from "@/util/functions";
import { PlayerApi } from "./PlayerApi";
import { PlayContextType } from "@/types/playcontext";
import { LocalPlayerContext } from "../types/playcontext";

let deviceID: string | null = null;
let shuffle = false
let repeat: "off"|"all"|"one" = "off"

let connectedDeviceID: string | null = null;

export const setDeviceID = (id: string) => {
    deviceID = id;
    localStorage.setItem("brage_device_id", id);
};

export const getDeviceID = () => {
    return deviceID ?? localStorage.getItem("brage_device_id");
};

export const removeDeviceID = () => {
    deviceID = null;
    localStorage.removeItem("brage_device_id");
};


type EventCallback = (data: any) => void;

export class ServerApi implements Api, PlayerApi {
    private listeners: Map<Event, Set<EventCallback>> = new Map();
    private eventSource?: EventSource;

    constructor() {
        this.startSSE()
    }

    private base = ky.extend({
        retry: 0,
        hooks: {
            afterResponse: [
                async ({ response }) => {
                    if (response.status >= 400) {
                        let title = "Error";
                        let msg = "";
                        let code = "";

                        try {
                            const clone = response.clone();
                            const text = await clone.text();

                            msg = text;

                            try {
                                const json = JSON.parse(text);
                                if (isBragErr(json.error)) {
                                    title = json.error.title;
                                    msg = json.error.message;
                                    code = json.error.code;
                                }
                            } catch {}

                        } catch (e) {
                            msg = "could not parse response";
                        }

                        this.emitEvent(Event.MsgErr, {title: title, message: msg})

                        if (response.status == 401 && code != "INVALIDUSERCREDS") {
                            this.emitEvent(Event.UserUpdated, null);
                        }

                    }

                    return response;
                },
            ],
        },
    });

    private auth = this.base.extend({
        prefix: "/auth",
    });

    private api = this.base.extend({
        prefix: "/api",
    });

    private mediaApi = this.api.extend({
        prefix: "/api",
    });


    private pickImage = (): Promise<File | null> => {
        return new Promise((resolve) => {
            const input = document.createElement("input");

            input.type = "file";
            input.accept = "image/*";

            input.onchange = () => {
                resolve(input.files?.[0] ?? null);
            };

            input.click();
        });
    };

    private async uploadImage(
        id: string,
        imgType: "artist" | "album" | "playlist" | "user",
        file: File | serverclient.FileUpload,
    ): Promise<void> {
        const form = new FormData();

        if (file instanceof File) {
            form.append("file", file);
        } else {
            const uint8 = new Uint8Array(file.data);
            const blob = new Blob([uint8], { type: file.type });

            form.append("file", blob, file.name);
        }

        if ( imgType == "user" ) {
            await this.api.post(`/users/me/img`, {
                body: form,
            });
            return
        }

        await this.api.post(`/img/${imgType}s/${id}`, {
            body: form,
        });
    }

    eventSubscribe(eventName: Event, callback: EventCallback): () => void {
        if (!this.listeners.has(eventName)) {
            this.listeners.set(eventName, new Set());
        }

        const set = this.listeners.get(eventName)!;
        set.add(callback);

        // return unsubscribe function
        return () => {
            set.delete(callback);

            // optional cleanup
            if (set.size === 0) {
                this.listeners.delete(eventName);
            }
        };
    }

    emitEvent(eventName: Event, data?: any) {
        const set = this.listeners.get(eventName);
        if (!set) return;

        for (const cb of set) {
            cb(data);
        }
    }

    private emitDeviceEvent(eventName: Event, data: any) {
        if (eventName == Event.PlayerPlaybackState) {
            eventName = Event.DevicePlaybackState;
            if (data.device_id && connectedDeviceID === data.device_id) {
                eventName = Event.PlayerPlaybackChange;
                shuffle = data.shuffle
                repeat = data.repeat
            }
        }

        if (eventName == Event.PlayerPlayContext) {
            eventName = Event.DevicePlayContext;
            if (data.device_id && connectedDeviceID === data.device_id) {
                eventName = Event.PlayerContextChange;
            }
        }
        this.emitEvent(eventName, data);
    }

    async loginServerUser(
        email: string,
        password: string,
        rememberMe: boolean,
    ): Promise<void> {
        const req = {
            email: email,
            password: password,
            rememberMe: rememberMe,
        };
        await this.auth.post("/login", { json: req }).json<responses.Login>();

        const user = await this.getUser();

        this.startSSE();

        this.emitEvent(Event.UserUpdated, user);
    }

    async getUser(): Promise<UserDetails> {
        const user = await this.api.get("/user").json<UserDetails>();
        this.emitEvent(Event.UserUpdated, user);
        return user;
    }

    async logoutLocalUser(): Promise<void> {
        await this.auth.get("/logout");
        this.emitEvent(Event.UserUpdated, null);
    }

    async serverStatus(): Promise<types.ServerApiInfo> {
        const status = await this.api.get("/info").json<types.ServerApiInfo>();
        return status
    }

    async publicServerStatus(): Promise<types.ServerApiInfo> {
        const status = await this.base.get("/info").json<types.ServerApiInfo>();
        return status
    }

    async getAlbum(id: string): Promise<types.AlbumDetailed> {
        return await this.mediaApi
            .get(`/albums/${id}/detailed`)
            .json<types.AlbumDetailed>();
    }

    async listAlbums(): Promise<Array<types.AlbumDetailed>> {
        const resp = await this.mediaApi
            .get("/albums")
            .json<responses.ListPayload<types.AlbumDetailed>>();
        if (resp.items) {
            return resp.items;
        }

        return [];
    }

    async listTracksByAlbum(
        albumID: string,
    ): Promise<Array<types.TrackDetailed>> {
        const resp = await this.mediaApi
            .get(`/albums/${albumID}/tracks-detailed`)
            .json<responses.ListPayload<types.TrackDetailed>>();
        if (resp.items) {
            return resp.items;
        }

        return [];
    }

    async updateAlbum(id: string, data: types.AlbumUpdate): Promise<void> {
        await this.mediaApi.put(`/albums/${id}`, { json: data });
        this.emitEvent(Event.EntitiesUpdated);
    }

    async uploadAlbumImage(
        id: string,
        file: File | serverclient.FileUpload,
    ): Promise<void> {
        return this.uploadImage(id, "album", file);
    }

    async countAlbums(): Promise<number> {
        const resp = await this.mediaApi
            .get("/albums", { searchParams: { count: true } })
            .json<responses.ListPayload<types.AlbumDetailed>>();
        if (resp.count) {
            return resp.count;
        }

        return 0;
    }

    async getTrack(trackID: string, albumID: string): Promise<types.TrackDetailed> {
        return await this.mediaApi.get(`/albums/${albumID}/tracks/${trackID}`).json<types.TrackDetailed>();
    }

    async filterTracks(filter: types.TrackFilter, page:number, limit:number):Promise<responses.ListPaginationPayload<types.TrackDetailed>> {
        const resp = await this.mediaApi
            .post("/tracks/search", {
                searchParams: {
                    page: page,
                    limit: limit,
                },
                body: JSON.stringify(filter),
                headers: {
                    "Content-Type": "application/json",
                },
            })
            .json<responses.ListPaginationPayload<types.TrackDetailed>>();
        return resp
    }

    async likeTrack(id: string): Promise<void> {
        await this.mediaApi.post(`/tracks/${id}/like`);
        this.emitEvent(Event.EntitiesUpdated);
    }

    async listLikedTracks(): Promise<Array<types.TrackDetailed>> {
        const resp = await this.mediaApi
            .get(`/tracks/liked`)
            .json<responses.ListPayload<types.TrackDetailed>>();
        if (resp.items) {
            return resp.items;
        }

        return [];
    }

    async unlikeTrack(id: string): Promise<void> {
        await this.mediaApi.delete(`/tracks/${id}/like`);
        this.emitEvent(Event.EntitiesUpdated);
    }

    async updateTrack(id: string, data: types.TrackUpdate): Promise<void> {
        await this.mediaApi.put(`/tracks/${id}`, { json: data });
        this.emitEvent(Event.EntitiesUpdated);
    }

    async countLikedTracks(): Promise<number> {
        const resp = await this.mediaApi
            .get("/tracks/liked", { searchParams: { count: true } })
            .json<responses.ListPayload<types.ArtistDetailed>>();
        if (resp.count) {
            return resp.count;
        }

        return 0;
    }

    async countTracks(): Promise<number> {
        const resp = await this.mediaApi
            .get("/tracks", { searchParams: { count: true } })
            .json<responses.ListPayload<types.TrackDetailed>>();
        if (resp.count) {
            return resp.count;
        }

        return 0;
    }

    async addPlayCount(trackID: string): Promise<void> {
        await this.mediaApi.post(`/tracks/${trackID}/play-history`);
        this.emitEvent(Event.EntitiesUpdated);
    }

    async listArtists(): Promise<Array<types.ArtistDetailed>> {
        const resp = await this.mediaApi
            .get("/artists")
            .json<responses.ListPayload<types.ArtistDetailed>>();
        if (resp.items) {
            return resp.items;
        }

        return [];
    }

    async countArtists(): Promise<number> {
        const resp = await this.mediaApi
            .get("/artists", { searchParams: { count: true } })
            .json<responses.ListPayload<types.ArtistDetailed>>();
        if (resp.count) {
            return resp.count;
        }

        return 0;
    }

    async createArtist(artist: types.ArtistBase): Promise<void> {
        await this.mediaApi.post(`/artists`, { json: artist});
        this.emitEvent(Event.EntitiesUpdated);
    }

    async getArtist(id: string): Promise<types.Artist> {
        return await this.mediaApi.get(`/artists/${id}`).json<types.Artist>();
    }

    async getArtistTopTracks(id: string): Promise<Array<types.TrackDetailed>> {
        const resp = await this.mediaApi
            .get(`/artists/${id}/top-tracks`)
            .json<responses.ListPayload<types.TrackDetailed>>();
        if (resp.items) {
            return resp.items;
        }

        return [];
    }

    async listAlbumsByArtist(id: string): Promise<Array<types.AlbumDetailed>> {
        const resp = await this.mediaApi
            .get(`/artists/${id}/albums`)
            .json<responses.ListPayload<types.AlbumDetailed>>();
        if (resp.items) {
            return resp.items;
        }

        return [];
    }

    async listFeaturedAlbumsByArtist(id: string): Promise<Array<types.AlbumDetailed>> {
        const resp = await this.mediaApi
            .get(`/artists/${id}/albums/featured`)
            .json<responses.ListPayload<types.AlbumDetailed>>();
        if (resp.items) {
            return resp.items;
        }

        return [];
    }

    async updateArtist(id: string, data: types.Artist): Promise<void> {
        await this.mediaApi.put(`/artists/${id}`, { json: data });
        this.emitEvent(Event.EntitiesUpdated);
    }

    async uploadArtistImage(
        id: string,
        file: File | serverclient.FileUpload,
    ): Promise<void> {
        return this.uploadImage(id, "artist", file);
    }

    // ---- PLAYLISTS ----
    async addPlaylist(playlist: types.Playlist): Promise<void> {
        await this.mediaApi.post(`/playlists`, { json: playlist });
        this.emitEvent(Event.EntitiesUpdated);
    }

    async addSmartPlaylist(playlist: types.PlaylistBase, filter: types.TrackFilter): Promise<void> {
        await this.mediaApi.post(`/playlists`, { json: {...playlist, filter: filter }, searchParams: {type: "smart"} });
        this.emitEvent(Event.EntitiesUpdated);
    }

    async countPlaylists(): Promise<number> {
        const resp = await this.mediaApi
            .get("/playlists", { searchParams: { count: true } })
            .json<responses.ListPayload<types.Playlist>>();
        if (resp.count) {
            return resp.count;
        }

        return 0;
    }

    async deletePlaylist(id: string): Promise<void> {
        await this.mediaApi.delete(`/playlists/${id}`);
    }

    async getPlaylist(id: string): Promise<types.Playlist> {
        return await this.mediaApi.get(`/playlists/${id}`).json<types.Playlist>();
    }

    async listPlaylists(): Promise<Array<types.Playlist>> {
        const resp = await this.mediaApi
            .get("/playlists", {
                searchParams: {
                    includePublic: true,
                    sortBy: "date",
                    sortOrder: "ASC",
                },
            })
            .json<responses.ListPayload<types.Playlist>>();

        if (resp.items) {
            return resp.items;
        }

        return [];
    }

    async updatePlaylist(id: string, data: types.Playlist): Promise<void> {
        await this.mediaApi.put(`/playlists/${id}`, { json: data });
        this.emitEvent(Event.EntitiesUpdated);
    }

    async uploadPlaylistImage(
        id: string,
        file: File | serverclient.FileUpload,
    ): Promise<void> {
        return this.uploadImage(id, "playlist", file);
    }

    async addPlaylistTrack(
        playlistID: string,
        albumID: string,
        trackID: string,
    ): Promise<void> {
        await this.mediaApi.post(`/playlists/${playlistID}/track`, {
            json: { album_id: albumID, track_id: trackID },
        });
        this.emitEvent(Event.EntitiesUpdated);
    }

    async countPlaylistTracks(playlistID: string): Promise<number> {
        const resp = await this.mediaApi
            .get(`/playlists/${playlistID}/tracks`, {
                searchParams: {
                    count: true,
                },
            })
            .json<responses.ListPayload<types.TrackDetailed>>();

        if (resp.count) {
            return resp.count;
        }

        return 0;
    }

    async deletePlaylistTrack(id: string): Promise<void> {
        await this.mediaApi.delete(`/playlist-tracks/${id}`);
    }

    async listPlaylistTracks(
        playlistID: string,
    ): Promise<Array<types.TrackDetailed>> {
        const resp = await this.mediaApi
            .get(`/playlists/${playlistID}/tracks`, {
                searchParams: {
                    sortBy: "date",
                    sortOrder: "ASC",
                },
            })
            .json<responses.ListPayload<types.TrackDetailed>>();

        if (resp.items) {
            return resp.items;
        }

        return [];
    }

    async SearchFull(searchTerm: string): Promise<Array<types.SearchItem>> {
        const resp = await this.mediaApi
            .get(`/search/media`, {
                searchParams: {
                    q: searchTerm,
                },
            })
            .json<responses.ListPayload<types.SearchItem>>();

        if (resp.items) {
            return resp.items;
        }

        return [];
    }

    async importAlbum(file: File, musicbrainzID: string | null): Promise<void> {
        const form = new FormData();

        form.append(
            "metadata",
            JSON.stringify({
                musicbrainz_id: musicbrainzID,
            }),
        );

        form.append("file", file, file.name);
        await this.api.post("/import/album", { body: form, timeout: false });
    }

    async listImportItems(page:number, limit:number):Promise<responses.ListPaginationPayload<types.Import>> {
        const resp = await this.mediaApi
            .get("/import", {
                searchParams: {
                    page: page,
                    limit: limit,
                },
            })
            .json<responses.ListPaginationPayload<types.Import>>();

        return resp;
    }

    async searchMusicBrainz(artist:string, album:string):Promise<responses.ListPayload<musicbrainz.SearchResults>> {
        const resp = await this.mediaApi
            .get("/import/mb/search", {
                searchParams: {
                    artist: artist,
                    album: album,
                },
            })
            .json<responses.ListPayload<musicbrainz.SearchResults>>();

        return resp;
    }

    async rateTrack(trackID: string, value: number): Promise<void> {
        await this.mediaApi.post(`/tracks/${trackID}/ratings`, {
            json: { value: value },
        });
        this.emitEvent(Event.EntitiesUpdated);
    }

    async listEntityEvents(): Promise<Array<types.EntityEvent>> {
        const resp = await this.api
            .get("/admin/entity-events")
            .json<responses.ListPayload<types.EntityEvent>>();
        if (resp.items) {
            return resp.items;
        }

        return [];
    }

    // ==== USER MANAGEMENT ==== //
    async createUser(
        email: string,
        username: string,
        password: string,
        roles: Array<string>,
    ): Promise<void> {
        await this.api.post(`/users`, {
            json: {
                email: email,
                username: username,
                password: password,
                roles: roles,
            },
        });
    }

    async deleteUser(id: string): Promise<void> {
        await this.api.delete(`/users/${id}`);
    }

    async listUsers(includeMachineUsers: boolean): Promise<Array<UserDetails>> {
        const resp = await this.api
            .get("/users", {
                searchParams: { machineUsers: includeMachineUsers },
            })
            .json<responses.ListPayload<UserDetails>>();
        if (resp.items) {
            return resp.items;
        }

        return [];
    }

    async listUserRoles(): Promise<Array<string>> {
        const resp = await this.auth
            .get("/user-roles")
            .json<responses.ListPayload<string>>();
        if (resp.items) {
            return resp.items;
        }

        return [];
    }

    async updateUser(
        id: string,
        email: string,
        username: string,
        password: any,
        roles: Array<string>,
    ): Promise<void> {
        await this.api.post(`/users/${id}`, {
            json: {
                email: email,
                username: username,
                password: password,
                roles: roles,
            },
        });
        await this.getUser()
    }

    public uploadUserImage = async (): Promise<void> => {
        const file = await this.pickImage();

        if (!file) {
            return;
        }

        return this.uploadImage("", "user", file);
    }

    async updateProfile(
        email?: string,
        username?: string,
        password?: string,
        new_password?: string,
        new_password_confirm?: string,
    ): Promise<void> {
        await this.api.put(`/users/me`, {
            json: {
                email: email,
                username: username,
                password: password,
                new_password: new_password,
                new_password_confirm: new_password_confirm
            }
        });
    }

    async listUserTokens(): Promise<Array<types.TokenLimited>> {
        const resp = await this.api
            .get("/users/me/tokens")
            .json<responses.ListPayload<types.TokenLimited>>();
        if (resp.items) {
            return resp.items;
        }

        return [];
    }

    async addApiToken(name: string): Promise<responses.CreateApiToken> {
        return await this.api.post(`/users/me/tokens/api`, {
            json: {
                name: name
            }
        }).json<responses.CreateApiToken>();
    }

    async deleteToken(id: string): Promise<void> {
        await this.api.delete(`/users/me/tokens/` + id)
        return
    }

    // !==== USER MANAGEMENT ==== //

    async supportsSync(): Promise<boolean> {
        return false;
    }

    async syncLibrary(): Promise<void> {}

    // ==== DEVICE MANAGEMENT ==== //
    //
    private async registerDevice(id: string | null): Promise<string> {
        const req: requests.DevicesRegister = {
            id: id === null ? undefined : id,
            name: "Website",
            type: "streaming",
            interface: "web",
            icon: "generic",
            supports_playback: false,
            platform: "web",
            version: import.meta.env.VITE_APP_VERSION,
        };

        const json = await this.api
            .post("/devices", { json: req })
            .json<responses.DeviceRegister>();

        return json.id;
    }

    private subscribeDeviceEvents(deviceID: string) {
        const url = `/api/devices/${deviceID}/events`;

        if (this.eventSource) {
            console.warn("Device event listener already active");
            return;
        }

        this.eventSource = new EventSource(url, {
            withCredentials: true, // important if using cookies
        });

        this.eventSource.onmessage = (event) => {
            try {
                const data: SSEvent = JSON.parse(event.data);

                if (data.type.startsWith("player.")) {
                    this.emitDeviceEvent(data.type, data.data);
                } else {
                    this.emitEvent(data.type, data.data);
                }
            } catch (err) {
                console.error("Failed to parse SSE event", err);
            }
        };

        this.eventSource.onerror = (err) => {
            console.error("SSE error", err);
            // browser will auto-reconnect
        };
    }

    // private unsubscribeDeviceEvents() {
    //     this.eventSource?.close();
    //     this.eventSource = undefined;
    // }

    private async startSSE(): Promise<void> {
        try {
            setDeviceID(await this.registerDevice(getDeviceID()));
        } catch (error) {
            const httpError = error as HTTPError;
            if (httpError.response.status == 404) {
                removeDeviceID()
                this.startSSE()
                return
            }
            console.log(httpError.response.status, error)
            return
        }
        const id = getDeviceID();
        if (id != null) {
            this.subscribeDeviceEvents(id);
        }
    }

    async getConnectedDeviceID(): Promise<string|null> {
        return connectedDeviceID
    }

    async connectDevice(id: string): Promise<void> {
        connectedDeviceID = id;

        const devices = await this.listDevices()
        const device = devices.find((d) => d.id == id)

        this.emitEvent(Event.PlayerContextChange, device?.player_state?.context);
        this.emitEvent(Event.PlayerPlaybackChange, device?.player_state?.playback);
        this.emitEvent(Event.DeviceConnectionID, connectedDeviceID);
    }

    async disconnectDevice(): Promise<void> {
        connectedDeviceID = null;
        this.emitEvent(Event.PlayerContextChange, {
            track_order: [],
        });
        this.emitEvent(Event.PlayerPlaybackChange, {
            shuffle: false,
            repeat: "off",
            track_source: "context",
        });

        this.emitEvent(Event.DeviceConnectionID, null);
    }

    async listDevices(): Promise<Array<types.DeviceDetailed>> {
        const resp = await this.api
            .get("/devices")
            .json<responses.ListPayload<types.DeviceDetailed>>();
        if (resp.items) {
            return resp.items;
        }

        return [];
    }

    async removeDevice(_id: string): Promise<void> {}
    async removeDeviceToken(_id: string): Promise<void> {}
    async removeDeviceAndToken(_id: string): Promise<void> {}

    // !==== DEVICE MANAGEMENT ==== //

    async listThemes(): Promise<Array<types.ThemeDescription>> {
        return [];
    }

    // ==== PLAYER API ==== //
    async addTrackToQueue(trackID: string, albumID: string): Promise<void> {
        if (connectedDeviceID == null) {
            return;
        }

        const track = await this.getTrack(trackID, albumID)
        await this.api.post(`/devices/${connectedDeviceID}/player/queue`, {json: track});
    }

    async nextTrack(): Promise<void> {
        if (connectedDeviceID == null) {
            this.emitEvent(Event.PlayerLocalNextTrack)
            return
        }
        await this.api.post(`/devices/${connectedDeviceID}/player/next`);
    }

    async nextRepeat(): Promise<void> {
        let newRepeat = repeat
        switch (repeat) {
            case "off":
                newRepeat = "all"
                break
            case "all":
                newRepeat = "one"
                break
            case "one":
                newRepeat = "off"
                break
        }

        if (connectedDeviceID == null) {
            this.emitEvent(Event.PlayerLocalRepeat, newRepeat)
            repeat = newRepeat
            return;
        }

        await this.api.post(`/devices/${connectedDeviceID}/player/repeat`, {json: {type: newRepeat}});
    }

    async playPause(): Promise<void> {
        if (connectedDeviceID == null) {
            this.emitEvent(Event.PlayerLocalPlayPause)
            return;
        }
        await this.api.post(`/devices/${connectedDeviceID}/player/play-pause`);
    }

    async previousTrack(): Promise<void> {
        if (connectedDeviceID == null) {
            this.emitEvent(Event.PlayerLocalPreviousTrack)
            return;
        }
        await this.api.post(`/devices/${connectedDeviceID}/player/previous`);
    }

    async start(
        parentType: PlayContextType,
        parentId: string,
        idx: number,
    ): Promise<void> {
        switch (parentType) {
            case PlayContextType.Album:
                this.startPlayerWithAlbum(parentId, idx)
                return
            case PlayContextType.LikedTracks:
                this.startPlayerWithLikedTracks(parentId, idx)
                return
            case PlayContextType.Playlist:
                this.startPlayerWithPlaylist(parentId, idx)
                return
            case PlayContextType.TopTracks:
        }
    }

    private async startPlayerWithAlbum(parentId: string, idx: number) {
        const tracks = await this.listTracksByAlbum(parentId)

        if (connectedDeviceID == null) {
            const ctx: LocalPlayerContext = {
                type: PlayContextType.Album,
                ref_id: parentId,
                tracks: tracks,
                track_index: idx,
            }
            this.emitEvent(Event.PlayerLocalStartContext, ctx)
            return
        }

        const state: types.PlayerState = new types.PlayerState({
            playback: new types.PlaybackState({
                track_index: idx,
            }),
            context: new types.PlayContext({
                type: PlayContextType.Album,
                ref_id: parentId,
                tracks: tracks,
            })
        })

        await this.api.post(`/devices/${connectedDeviceID}/player/state`, {json: state});
    }

    private async startPlayerWithPlaylist(parentId: string, idx: number) {
        if (connectedDeviceID == null) {
            return;
        }

        const tracks = await this.listPlaylistTracks(parentId)
        const state: types.PlayerState = new types.PlayerState({
            playback: new types.PlaybackState({
                track_index: idx,
            }),
            context: new types.PlayContext({
                type: PlayContextType.Playlist,
                ref_id: parentId,
                tracks: tracks,
            })
        })

        await this.api.post(`/devices/${connectedDeviceID}/player/state`, {json: state});
    }

    private async startPlayerWithLikedTracks(parentId: string, idx: number) {
        if (connectedDeviceID == null) {
            return;
        }

        const tracks = await this.listLikedTracks()
        const state: types.PlayerState = new types.PlayerState({
            playback: new types.PlaybackState({
                track_index: idx,
            }),
            context: new types.PlayContext({
                type: PlayContextType.LikedTracks,
                ref_id: parentId,
                tracks: tracks,
            })
        })

        await this.api.post(`/devices/${connectedDeviceID}/player/state`, {json: state});
    }

    async toggleShuffle(): Promise<void> {
        if (connectedDeviceID == null) {
            this.emitEvent(Event.PlayerLocalShuffle, !shuffle)
            shuffle = !shuffle
            return;
        }

        await this.api.post(`/devices/${connectedDeviceID}/player/shuffle`, {json: {active: !shuffle}});
    }

    // !==== PLAYER API ==== //
}
