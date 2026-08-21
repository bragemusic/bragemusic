import { types as wtypes } from "@/types/core";
import { TrackSource } from "@/types/playcontext";

export class PlayerState extends wtypes.PlayerState {
  constructor(source: any = {}) {
    super(source);

    if (!this.playback) {
      this.playback = new wtypes.PlaybackState({
        track_source: TrackSource.Context,
        track_index: 0,
        playing: false,
        progress_ms: 0,
      });
    }

    if (!this.context) {
      this.context = new wtypes.PlayContext({
        tracks: [],
        queue: [],
        track_order: [],
      });
    }
  }

  get current_track(): wtypes.TrackDetailed | null {
    if (!this.playback) {
      return null
    }

    switch (this.playback.track_source) {
      case TrackSource.Context: {
        // nothing loaded yet
        if (
          !this.context ||
          this.context.track_order === null ||
          this.context.track_order.length === 0
        ) {
          return null;
        }

        if (this.context.track_order.length <= this.playback.track_index) {
          throw new Error("track index oob");
        }

        const idx = this.context.track_order[this.playback.track_index];

        if (this.context.tracks.length <= idx) {
          throw new Error("track index oob");
        }

        return this.context.tracks[idx];
      }

      case TrackSource.Queue: {
        if (this.context.queue === null || this.context.queue.length === 0) {
          return null;
        }

        return this.context.queue[0];
      }

      default:
        throw new Error("unknown track source");
    }
  }
}

export function toPlayerState(raw: wtypes.PlayerState): PlayerState {
  const ps = new PlayerState();

  Object.assign(ps, raw);

  return ps;
}
