import { useEffect, useRef, useState } from "react";

import { useApi} from "@/api/ApiContext";
import { Event } from "@/types/events.ts";
import { types } from "@/types/core";
import { LocalPlayerContext, TrackSource } from "../types/playcontext";
import { timeNow } from "../util/functions";

function shuffle<T>(array: T[]): T[] {
    const result = [...array];

    for (let i = result.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));

        [result[i], result[j]] = [result[j], result[i]];
    }

    return result;
}

export function LocalPlayer() {
  const audioRef = useRef<HTMLAudioElement>(null);
  const api = useApi();

    const [ctx, setCtx] = useState<types.PlayContext>({
        type: "",
        ref_id: "",
        tracks: [],
        track_order: [],
        queue: [],
    })
    const [pb, setPb] = useState<types.PlaybackState>({
        playing: false,
        shuffle: false,
        repeat: "off",
        progress: 0,
        updated_at: "",
        track_source: TrackSource.Context,
        track_index: 0,
    })
    const ctxRef = useRef(ctx);
    const pbRef = useRef(pb);

    const [currentTrack, setCurrentTrack] = useState(new types.TrackDetailed)

    useEffect(() => {
        ctxRef.current = ctx;
    }, [ctx]);

    useEffect(() => {
        pbRef.current = pb;
    }, [pb]);

useEffect(() => {
    const currentTrack =
        ctxRef.current.track_order[pbRef.current.track_index];

    const { trackOrder, trackIndex } = shuffleTrackOrder(
        ctxRef.current.tracks,
        pb.shuffle,
        currentTrack,
    );

    setCtx(prev => ({
        ...prev,
        track_order: trackOrder,
    }));

    updatePb(prev => ({
        ...prev,
        track_index: trackIndex,
    }));
}, [pb.shuffle]);

    const passed70Ref = useRef(false);

    const handleTimeUpdate = () => {
        const audio = audioRef.current;
        if (!audio || !audio.duration) {
            return;
        }

        const progress = audio.currentTime / audio.duration;

        if (progress >= 0.75 && !passed70Ref.current) {
            passed70Ref.current = true;
            api.addPlayCount(currentTrack.id)
        }
    };

    const updatePb = (update: (prev: types.PlaybackState) => types.PlaybackState, progress?: number) => {
        setPb(prev => ({
            ...update(prev),
            progress: progress ?? Math.round(audioRef.current?.currentTime ? audioRef.current?.currentTime * 1000 : 0),
            updated_at: timeNow(),
        }));
    };

    useEffect(() => {
        const unsubscribePlayerLocalStartContext = api.eventSubscribe(Event.PlayerLocalStartContext, (lctx: LocalPlayerContext) => {
            const { trackOrder, trackIndex } = shuffleTrackOrder(
                lctx.tracks,
                pbRef.current.shuffle,
                lctx.track_index,
            );

            updatePb(pb => ({
                ...pb,
                playing: true,
                track_index: trackIndex,
            }));

            setCtx(prev => new types.PlayContext({
                ...prev,
                type: lctx.type,
                ref_id: lctx.ref_id,
                tracks: lctx.tracks,
                track_order: trackOrder,
            }));
            }
        );

        const unsubscribePlayPause = api.eventSubscribe(
            Event.PlayerLocalPlayPause,
            () => {
            updatePb((prev) => ({
                ...prev,
                playing: !prev.playing,
            }));

            }
        );

        const unsubscribePlayerLocalNextTrack = api.eventSubscribe(Event.PlayerLocalNextTrack, () => {
                nextTrack();
            }
        );

        const unsubscribePlayerLocalPreviousTrack = api.eventSubscribe(Event.PlayerLocalPreviousTrack, () => {
                previousTrack();
            }
        );

        const unsubscribePlayerLocalRepeat = api.eventSubscribe(Event.PlayerLocalRepeat, (repeat: string) => {
            updatePb((prev) => ({
                ...prev,
                repeat: repeat,
            }));
            }
        );

        const unsubscribePlayerLocalShuffle= api.eventSubscribe(Event.PlayerLocalShuffle, (shuffle: boolean) => {
            updatePb((prev) => ({
                ...prev,
                shuffle: shuffle,
            }));
            }
        );

        return () => {
            unsubscribePlayerLocalStartContext?.();
            unsubscribePlayPause?.();
            unsubscribePlayerLocalNextTrack?.();
            unsubscribePlayerLocalPreviousTrack?.();
            unsubscribePlayerLocalRepeat?.();
            unsubscribePlayerLocalShuffle?.();
        };
    }, [api]);

    const evalContextAndState = (ctx: types.PlayContext, pb: types.PlaybackState) => {
        if (ctx.type != "album") {
        console.error("only album implemented for playcontext")
        return
        }

        if (ctx.tracks.length == 0 || pb.track_index >= ctx.tracks.length) {
        return
        }


        setCurrentTrack(ctx.tracks[ctx.track_order[pb.track_index]])
    }

    useEffect(() => {
        api.emitEvent(Event.PlayerContextChange, ctx)
    }, [ctx]);

    useEffect(() => {
        if (!pb.track_source) {
            return
        }
        api.emitEvent(Event.PlayerPlaybackChange, pb)
    }, [pb]);

    useEffect(() => {
        evalContextAndState(ctx, pb)
        // audioRef.current?.play();
    }, [pb, ctx]);

    useEffect(() => {
        startCurrentTrack()
    }, [currentTrack]);


    useEffect(() => {
        if (pb.playing) {
            audioRef.current?.play();
        } else {
            audioRef.current?.pause();
        }
    }, [pb.playing]);


    const startCurrentTrack = () => {
        passed70Ref.current = false;
        setPb((prev) => ({
            ...prev,
            progress: 0,
        }))
        if (!currentTrack.media_file) {
            return
        }
        audioRef.current!.src = "/api/mediafiles/" + currentTrack.media_file.id + "/file";
        audioRef.current?.play();
    }

    const stop = () => {
        audioRef.current?.pause()
        audioRef.current?.removeAttribute("src");
        setCtx((prev) => (new types.PlayContext({
            ...prev,
            ref_id: "",
            tracks: [],
            track_order: [],
            queue: [],
        })))
        updatePb((prev) => ({
            ...prev,
            playing: false,
            progress: 0,
            track_sourece: TrackSource.Context,
            track_index: 0,
        }))
    }

    const nextTrack = () => {
        const ctx = ctxRef.current;
        const pb = pbRef.current;

        passed70Ref.current = false;

        if (!audioRef.current) {
            return
        }

        if (pb.repeat == "one") {
            audioRef.current.currentTime = 0;
            audioRef.current.play();
            updatePb((prev) => {
                return {
                ...prev,
            }}, 0);
            return
        }

        let ntid = pb.track_index + 1;
        if (ntid >= ctx.track_order.length) {
            if (pb.repeat == "all") {
                ntid = 0
            } else {
                stop()
                return
            }
        }
        //FIXME: Check for shuffle
        updatePb((prev) => {
            return {
            ...prev,
            track_index: ntid,
        }});
    }

    const previousTrack = () => {
        const ctx = ctxRef.current;
        const pb = pbRef.current;

        passed70Ref.current = false;

        if (!audioRef.current) {
            return
        }

        if (pb.repeat == "one") {
            audioRef.current.currentTime = 0;
            audioRef.current.play();
            updatePb((prev) => {
                return {
                ...prev,
            }}, 0);
            return
        }

        let ntid = pb.track_index
        if (audioRef.current?.currentTime < 10) {
           ntid = ntid-1
        }

        if (ntid < 0) {
            if (pb.repeat == "all") {
                ntid = ctx.track_order.length - 1
            } else {
                stop()
                return
            }
        }

        audioRef.current.currentTime = 0;
        audioRef.current.play();
        //FIXME: Check for shuffle
        updatePb((prev) => {
            return {
            ...prev,
            track_index: ntid,
        }}, 0);
    }

const shuffleTrackOrder = (
    tracks: types.TrackDetailed[],
    shuffle: boolean,
    currentIndex: number,
): { trackOrder: number[]; trackIndex: number } => {
    const numberOfTracks = tracks.length;

    if (numberOfTracks === 0) {
        return {
            trackOrder: [],
            trackIndex: 0,
        };
    }

    if (!shuffle) {
        return {
            trackOrder: tracks.map((_, i) => i),
            trackIndex: currentIndex,
        };
    }

    const nums = tracks.map((_, i) => i);

    // Put selected track first
    [nums[0], nums[currentIndex]] = [nums[currentIndex], nums[0]];

    // Shuffle everything after the selected track
    for (let i = nums.length - 1; i > 1; i--) {
        const j = 1 + Math.floor(Math.random() * i);
        [nums[i], nums[j]] = [nums[j], nums[i]];
    }

    return {
        trackOrder: nums,
        trackIndex: 0,
    };
};


    return (
        <audio
        ref={audioRef}
            onTimeUpdate={handleTimeUpdate}
            onEnded={() => {
                api.emitEvent(Event.PlayerLocalNextTrack);
            }}
        />
    )
;
}
