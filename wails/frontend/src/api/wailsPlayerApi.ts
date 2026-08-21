import { PlayerApi, PlayContextType } from "bragemusic-frontend";
import { AddTrackToQueue, NextRepeat, NextTrack, PlayPause, PreviousTrack, Start, ToggleShuffle } from "../../wailsjs/go/app/App";

export class WailsPlayerApi implements PlayerApi {
    async addTrackToQueue(trackID:string, albumID:string): Promise<void> {
        return await AddTrackToQueue(trackID, albumID)
    }
    async nextTrack(): Promise<void> {
        return await NextTrack();
    }

    async nextRepeat(): Promise<void> {
        return await NextRepeat();
    }

    async playPause(): Promise<void> {
        return await PlayPause();
    }

    async previousTrack(): Promise<void> {
        return await PreviousTrack();
    }

    async start(parentType:PlayContextType, parentId:string, idx:number):Promise<void> {
        return await Start(parentType, parentId, idx)
    }

    async toggleShuffle(): Promise<void> {
        return await ToggleShuffle();
    }
}
