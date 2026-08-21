import { PlayContextType } from "@/types/playcontext";

export interface PlayerApi {
  addTrackToQueue(trackID: string, albumID: string): Promise<void>;
  nextTrack(): Promise<void>;
  nextRepeat(): Promise<void>;
  playPause(): Promise<void>;
  previousTrack(): Promise<void>;
  start(
    parentType: PlayContextType,
    parentId: string,
    idx: number,
  ): Promise<void>;
  toggleShuffle(): Promise<void>;
}
