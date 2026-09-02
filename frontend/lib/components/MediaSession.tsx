import { usePlayerState } from "@/store/playerStore.ts";
import { useEffect } from "react";
import { albumImageLink } from "@/util/images.tsx";
import { usePlayerApi } from "@/api/ApiContext.tsx";

export const MediaSession = () => {
    const playerApi = usePlayerApi();
    const playCtx = usePlayerState((s) => s.player);


    useEffect(() => {
        if (!("mediaSession" in navigator)) {
            return;
        }

        if (!playCtx?.current_track?.id) {
            return;
        }

        const currentTrack = playCtx.current_track

        const artistNames = currentTrack.artists?.map(a => a.name).join(", ");

        navigator.mediaSession.metadata = new MediaMetadata({
            title: currentTrack.title,
            artist: artistNames,
            album: currentTrack.album_name,
            artwork: [
            {
                src: albumImageLink(currentTrack.album_id as string, 1024),
                sizes: "1024x1024",
                type: "image/jpeg",
            },
            ],
        });

        navigator.mediaSession.setActionHandler("play", () => {
            playerApi.playPause()
        });

        navigator.mediaSession.setActionHandler("pause", () => {
            playerApi.playPause()
        });

        navigator.mediaSession.setActionHandler("previoustrack", () => {
            // api.emitEvent(Event.PlayerLocalPrevious);
        });

        navigator.mediaSession.setActionHandler("nexttrack", () => {
            playerApi.nextTrack()
        });
    }, [playCtx]);



    return (<></>)
}
