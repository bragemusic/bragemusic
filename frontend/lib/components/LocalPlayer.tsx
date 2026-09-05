import { useEffect, useRef, useState } from "react";

import { useApi} from "@/api/ApiContext";
import { Event } from "@/types/events.ts";
import { types } from "@/types/core";
import { LocalPlayerContext, TrackSource } from "../types/playcontext";
import { timeNow } from "../util/functions";

export function LocalPlayer() {
  const audioRef = useRef<HTMLAudioElement>(null);
  const api = useApi();

    const [ctx, setCtx] = useState(new types.PlayContext)
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
            setCtx((prev) => (new types.PlayContext({
                    ...prev,
                    type: lctx.type,
                    ref_id: lctx.ref_id,
                    tracks: lctx.tracks,
                    track_order: lctx.tracks.map((_, i) => i),
                })));
            updatePb((prev) => ({
                ...prev,
                playing: true,
                track_index: lctx.track_index,
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

        return () => {
            unsubscribePlayerLocalStartContext?.();
            unsubscribePlayPause?.();
            unsubscribePlayerLocalNextTrack?.();
            unsubscribePlayerLocalPreviousTrack?.();
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
    // setIsPlaying(pb.playing)



    // FIXME: Needs to handle shuffle and loop and so on. track_order is ignored atm
      //
    // const track =

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
    }, [currentTrack]);


    useEffect(() => {
        if (pb.playing) {
            audioRef.current?.play();
        } else {
            audioRef.current?.pause();
        }
    }, [pb.playing]);


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

        let ntid = pb.track_index + 1;
        if (ntid >= ctx.track_order.length) {
            stop()
            return {
                ...pb,
                playing: false,
                progress: 0,
                track_sourece: TrackSource.Context,
                track_index: 0,
            }
        }
        //FIXME: Check for shuffle and/or loop
        updatePb((prev) => {
            return {
            ...prev,
            track_index: ntid,
        }});
    }

    const previousTrack = () => {
        const pb = pbRef.current;

        passed70Ref.current = false;

        if (!audioRef.current) {
            return
        }

        let ntid = pb.track_index
        if (audioRef.current?.currentTime < 10) {
           ntid = ntid-1
        }

        if (ntid < 0) {
            stop()
            return {
                ...pb,
                playing: false,
                progress: 0,
                track_sourece: TrackSource.Context,
                track_index: 0,
            }
        }

        audioRef.current.currentTime = 0;
        audioRef.current.play();
        //FIXME: Check for shuffle and/or loop
        updatePb((prev) => {
            return {
            ...prev,
            track_index: ntid,
        }}, 0);
    }


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
