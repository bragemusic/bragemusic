export namespace musicbrainz {
	
	export class SearchResults {
	    id: string;
	    artist: string;
	    album: string;
	    score: number;
	
	    static createFrom(source: any = {}) {
	        return new SearchResults(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.artist = source["artist"];
	        this.album = source["album"];
	        this.score = source["score"];
	    }
	}

}

export namespace serverclient {
	
	export class FileUpload {
	    name: string;
	    type: string;
	    data: number[];
	
	    static createFrom(source: any = {}) {
	        return new FileUpload(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.data = source["data"];
	    }
	}

}

export namespace types {
	
	export class AlbumDetailed {
	    id: string;
	    musicbrainz_id?: string;
	    name: string;
	    sort_name: string;
	    artist_ids?: string[];
	    artist_names?: string[];
	    release_date?: string | undefined;
	    track_count: number;
	    disc_count: number;
	    description?: string;
	    owner: string;
	    public?: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new AlbumDetailed(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.musicbrainz_id = source["musicbrainz_id"];
	        this.name = source["name"];
	        this.sort_name = source["sort_name"];
	        this.artist_ids = source["artist_ids"];
	        this.artist_names = source["artist_names"];
	        this.release_date = source["release_date"];
	        this.track_count = source["track_count"];
	        this.disc_count = source["disc_count"];
	        this.description = source["description"];
	        this.owner = source["owner"];
	        this.public = source["public"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AlbumUpdate {
	    musicbrainz_id?: string;
	    name: string;
	    sort_name: string;
	    // Go type: time
	    release_date?: any;
	    tracks?: number;
	    discs?: number;
	    description?: string;
	    owner: string;
	    public?: boolean;
	    artists: number[][];
	
	    static createFrom(source: any = {}) {
	        return new AlbumUpdate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.musicbrainz_id = source["musicbrainz_id"];
	        this.name = source["name"];
	        this.sort_name = source["sort_name"];
	        this.release_date = this.convertValues(source["release_date"], null);
	        this.tracks = source["tracks"];
	        this.discs = source["discs"];
	        this.description = source["description"];
	        this.owner = source["owner"];
	        this.public = source["public"];
	        this.artists = source["artists"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Artist {
	    id: string;
	    musicbrainz_id?: string;
	    name: string;
	    sort_name: string;
	    country?: string;
	    year_started?: number;
	    year_ended?: number;
	    description?: string;
	    role: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Artist(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.musicbrainz_id = source["musicbrainz_id"];
	        this.name = source["name"];
	        this.sort_name = source["sort_name"];
	        this.country = source["country"];
	        this.year_started = source["year_started"];
	        this.year_ended = source["year_ended"];
	        this.description = source["description"];
	        this.role = source["role"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ArtistBase {
	    musicbrainz_id?: string;
	    name: string;
	    sort_name: string;
	    country?: string;
	    year_started?: number;
	    year_ended?: number;
	    description?: string;
	
	    static createFrom(source: any = {}) {
	        return new ArtistBase(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.musicbrainz_id = source["musicbrainz_id"];
	        this.name = source["name"];
	        this.sort_name = source["sort_name"];
	        this.country = source["country"];
	        this.year_started = source["year_started"];
	        this.year_ended = source["year_ended"];
	        this.description = source["description"];
	    }
	}
	export class ArtistDetailed {
	    id: string;
	    musicbrainz_id?: string;
	    name: string;
	    sort_name: string;
	    country?: string;
	    year_started?: number;
	    year_ended?: number;
	    description?: string;
	    role: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    album_count: number;
	    track_count: number;
	
	    static createFrom(source: any = {}) {
	        return new ArtistDetailed(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.musicbrainz_id = source["musicbrainz_id"];
	        this.name = source["name"];
	        this.sort_name = source["sort_name"];
	        this.country = source["country"];
	        this.year_started = source["year_started"];
	        this.year_ended = source["year_ended"];
	        this.description = source["description"];
	        this.role = source["role"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.album_count = source["album_count"];
	        this.track_count = source["track_count"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ArtistMinimal {
	    id: string;
	    name: string;
	    sort_name: string;
	    role: string;
	
	    static createFrom(source: any = {}) {
	        return new ArtistMinimal(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.sort_name = source["sort_name"];
	        this.role = source["role"];
	    }
	}
	export class TrackAnalysis {
	    id: string;
	    bpm: number;
	    key: string;
	    key_scale: string;
	    key_confidence: number;
	    loudness: number;
	    energy: number;
	    danceability: number;
	    mood_happy: number;
	    mood_sad: number;
	    mood_aggresive: number;
	    mood_calm: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new TrackAnalysis(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.bpm = source["bpm"];
	        this.key = source["key"];
	        this.key_scale = source["key_scale"];
	        this.key_confidence = source["key_confidence"];
	        this.loudness = source["loudness"];
	        this.energy = source["energy"];
	        this.danceability = source["danceability"];
	        this.mood_happy = source["mood_happy"];
	        this.mood_sad = source["mood_sad"];
	        this.mood_aggresive = source["mood_aggresive"];
	        this.mood_calm = source["mood_calm"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MediaFile {
	    id: number[];
	    duration_ms: number;
	    bitrate: number;
	    sample_rate: number;
	    file_size: number;
	    codec: string;
	    checksum: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new MediaFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.duration_ms = source["duration_ms"];
	        this.bitrate = source["bitrate"];
	        this.sample_rate = source["sample_rate"];
	        this.file_size = source["file_size"];
	        this.codec = source["codec"];
	        this.checksum = source["checksum"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TrackDetailed {
	    id: string;
	    title: string;
	    album_id?: string;
	    album_name?: string;
	    artists?: ArtistMinimal[];
	    musicbrainz_id?: string;
	    track_number?: number;
	    disc_number?: number;
	    genre?: string;
	    comment?: string;
	    media_file?: MediaFile;
	    play_count: number;
	    context_id?: string;
	    rating?: number;
	    user_rating?: number;
	    liked: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    analysis: TrackAnalysis;
	
	    static createFrom(source: any = {}) {
	        return new TrackDetailed(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.album_id = source["album_id"];
	        this.album_name = source["album_name"];
	        this.artists = this.convertValues(source["artists"], ArtistMinimal);
	        this.musicbrainz_id = source["musicbrainz_id"];
	        this.track_number = source["track_number"];
	        this.disc_number = source["disc_number"];
	        this.genre = source["genre"];
	        this.comment = source["comment"];
	        this.media_file = this.convertValues(source["media_file"], MediaFile);
	        this.play_count = source["play_count"];
	        this.context_id = source["context_id"];
	        this.rating = source["rating"];
	        this.user_rating = source["user_rating"];
	        this.liked = source["liked"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.analysis = this.convertValues(source["analysis"], TrackAnalysis);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PlayContextDTO {
	    device_id: string;
	    type: string;
	    ref_id: string;
	    tracks: TrackDetailed[];
	    track_order: number[];
	    queue: TrackDetailed[];
	
	    static createFrom(source: any = {}) {
	        return new PlayContextDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.device_id = source["device_id"];
	        this.type = source["type"];
	        this.ref_id = source["ref_id"];
	        this.tracks = this.convertValues(source["tracks"], TrackDetailed);
	        this.track_order = source["track_order"];
	        this.queue = this.convertValues(source["queue"], TrackDetailed);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PlaybackStateDTO {
	    device_id: string;
	    playing: boolean;
	    shuffle: boolean;
	    repeat: string;
	    progress: number;
	    updated_at: string;
	    track_source: string;
	    track_index: number;
	
	    static createFrom(source: any = {}) {
	        return new PlaybackStateDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.device_id = source["device_id"];
	        this.playing = source["playing"];
	        this.shuffle = source["shuffle"];
	        this.repeat = source["repeat"];
	        this.progress = source["progress"];
	        this.updated_at = source["updated_at"];
	        this.track_source = source["track_source"];
	        this.track_index = source["track_index"];
	    }
	}
	export class PlayerStateDTO {
	    playback: PlaybackStateDTO;
	    context: PlayContextDTO;
	
	    static createFrom(source: any = {}) {
	        return new PlayerStateDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.playback = this.convertValues(source["playback"], PlaybackStateDTO);
	        this.context = this.convertValues(source["context"], PlayContextDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DeviceDetailed {
	    name: string;
	    type: string;
	    interface: string;
	    icon: string;
	    supports_playback: boolean;
	    platform: string;
	    version: string;
	    id: string;
	    user_id: number[];
	    last_ip: string;
	    // Go type: time
	    last_seen: any;
	    active: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    player_state?: PlayerStateDTO;
	
	    static createFrom(source: any = {}) {
	        return new DeviceDetailed(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.interface = source["interface"];
	        this.icon = source["icon"];
	        this.supports_playback = source["supports_playback"];
	        this.platform = source["platform"];
	        this.version = source["version"];
	        this.id = source["id"];
	        this.user_id = source["user_id"];
	        this.last_ip = source["last_ip"];
	        this.last_seen = this.convertValues(source["last_seen"], null);
	        this.active = source["active"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.player_state = this.convertValues(source["player_state"], PlayerStateDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EntityEvent {
	    id: string;
	    item_id: string;
	    user_id: string;
	    event_type: string;
	    entity_type: string;
	    event_time: string;
	
	    static createFrom(source: any = {}) {
	        return new EntityEvent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.item_id = source["item_id"];
	        this.user_id = source["user_id"];
	        this.event_type = source["event_type"];
	        this.entity_type = source["entity_type"];
	        this.event_time = source["event_time"];
	    }
	}
	export class FilterMood {
	    aggressive?: number;
	    calm?: number;
	    happy?: number;
	    sad?: number;
	
	    static createFrom(source: any = {}) {
	        return new FilterMood(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.aggressive = source["aggressive"];
	        this.calm = source["calm"];
	        this.happy = source["happy"];
	        this.sad = source["sad"];
	    }
	}
	export class FilterUpperLower_int_ {
	    upper: number;
	    lower: number;
	
	    static createFrom(source: any = {}) {
	        return new FilterUpperLower_int_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.upper = source["upper"];
	        this.lower = source["lower"];
	    }
	}
	export class Import {
	    id: string;
	    musicbrainz_id?: string;
	    owner: string;
	    filename: string;
	    type: string;
	    state: string;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Import(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.musicbrainz_id = source["musicbrainz_id"];
	        this.owner = source["owner"];
	        this.filename = source["filename"];
	        this.type = source["type"];
	        this.state = source["state"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	export class ListPaginationPayload_github_com_bragemusic_bragemusic_pkg_types_Import_ {
	    items?: Import[];
	    page: number;
	    limit: number;
	    total_pages: number;
	    total_items: number;
	
	    static createFrom(source: any = {}) {
	        return new ListPaginationPayload_github_com_bragemusic_bragemusic_pkg_types_Import_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], Import);
	        this.page = source["page"];
	        this.limit = source["limit"];
	        this.total_pages = source["total_pages"];
	        this.total_items = source["total_items"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ListPaginationPayload_github_com_bragemusic_bragemusic_pkg_types_TrackDetailed_ {
	    items?: TrackDetailed[];
	    page: number;
	    limit: number;
	    total_pages: number;
	    total_items: number;
	
	    static createFrom(source: any = {}) {
	        return new ListPaginationPayload_github_com_bragemusic_bragemusic_pkg_types_TrackDetailed_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], TrackDetailed);
	        this.page = source["page"];
	        this.limit = source["limit"];
	        this.total_pages = source["total_pages"];
	        this.total_items = source["total_items"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ListPayload_github_com_bragemusic_bragemusic_pkg_musicbrainz_SearchResults_ {
	    items?: musicbrainz.SearchResults[];
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new ListPayload_github_com_bragemusic_bragemusic_pkg_musicbrainz_SearchResults_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], musicbrainz.SearchResults);
	        this.count = source["count"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class PlayContext {
	    type: string;
	    ref_id: string;
	    tracks: TrackDetailed[];
	    track_order: number[];
	    queue: TrackDetailed[];
	
	    static createFrom(source: any = {}) {
	        return new PlayContext(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.ref_id = source["ref_id"];
	        this.tracks = this.convertValues(source["tracks"], TrackDetailed);
	        this.track_order = source["track_order"];
	        this.queue = this.convertValues(source["queue"], TrackDetailed);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class PlaybackState {
	    playing: boolean;
	    shuffle: boolean;
	    repeat: string;
	    progress: number;
	    updated_at: string;
	    track_source: string;
	    track_index: number;
	
	    static createFrom(source: any = {}) {
	        return new PlaybackState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.playing = source["playing"];
	        this.shuffle = source["shuffle"];
	        this.repeat = source["repeat"];
	        this.progress = source["progress"];
	        this.updated_at = source["updated_at"];
	        this.track_source = source["track_source"];
	        this.track_index = source["track_index"];
	    }
	}
	
	export class PlayerState {
	    playback: PlaybackState;
	    context: PlayContext;
	
	    static createFrom(source: any = {}) {
	        return new PlayerState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.playback = this.convertValues(source["playback"], PlaybackState);
	        this.context = this.convertValues(source["context"], PlayContext);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Playlist {
	    id: string;
	    name: string;
	    description?: string;
	    public: boolean;
	    type: string;
	    owner: number[];
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Playlist(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.public = source["public"];
	        this.type = source["type"];
	        this.owner = source["owner"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PlaylistBase {
	    name: string;
	    description?: string;
	    public: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PlaylistBase(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.public = source["public"];
	    }
	}
	export class SearchItem {
	    name: string;
	    html_name: string;
	    id: number[];
	    type: string;
	    link_id: number[];
	    link_type: string;
	
	    static createFrom(source: any = {}) {
	        return new SearchItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.html_name = source["html_name"];
	        this.id = source["id"];
	        this.type = source["type"];
	        this.link_id = source["link_id"];
	        this.link_type = source["link_type"];
	    }
	}
	export class ServerApiInfo {
	    application: string;
	    id: string;
	    status: string;
	    name: string;
	    version: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerApiInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.application = source["application"];
	        this.id = source["id"];
	        this.status = source["status"];
	        this.name = source["name"];
	        this.version = source["version"];
	    }
	}
	export class ThemeDescription {
	    id: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new ThemeDescription(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class TokenLimited {
	    id: string;
	    type: string;
	    name?: string;
	    scopes: string;
	    expires_at?: string;
	    last_used_at?: string;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new TokenLimited(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.name = source["name"];
	        this.scopes = source["scopes"];
	        this.expires_at = source["expires_at"];
	        this.last_used_at = source["last_used_at"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	
	
	export class TrackFilter {
	    bpm?: FilterUpperLower_int_;
	    mood: FilterMood;
	    artists?: string[];
	
	    static createFrom(source: any = {}) {
	        return new TrackFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bpm = this.convertValues(source["bpm"], FilterUpperLower_int_);
	        this.mood = this.convertValues(source["mood"], FilterMood);
	        this.artists = source["artists"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TrackUpdate {
	    id: number[];
	    title: string;
	    musicbrainz_id?: string;
	    genre?: string;
	    comment?: string;
	    media_file?: number[];
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    disc_number: number;
	    track_number: number;
	    album_id: string;
	    artists?: number[][];
	
	    static createFrom(source: any = {}) {
	        return new TrackUpdate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.musicbrainz_id = source["musicbrainz_id"];
	        this.genre = source["genre"];
	        this.comment = source["comment"];
	        this.media_file = source["media_file"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.disc_number = source["disc_number"];
	        this.track_number = source["track_number"];
	        this.album_id = source["album_id"];
	        this.artists = source["artists"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UserDetails {
	    id: string;
	    email: string;
	    username: string;
	    provider: string;
	    role: string[];
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new UserDetails(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.email = source["email"];
	        this.username = source["username"];
	        this.provider = source["provider"];
	        this.role = source["role"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

