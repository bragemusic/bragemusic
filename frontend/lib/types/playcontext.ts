import { types } from "./core";

export enum PlayContextType {
  Album = "album",
  LikedTracks = "liked_tracks",
  Playlist = "playlist",
  TopTracks = "top_tracks",
  Filter = "filter",
}

export enum TrackSource {
  Context = "context",
  Queue = "queue",
}

export class LocalPlayerContext {
    type: PlayContextType = PlayContextType.Album;
    ref_id = "";
    tracks: types.TrackDetailed[] = [];
    track_index = 0;
}
