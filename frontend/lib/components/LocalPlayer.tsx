import { useEffect, useRef, useState } from "react";

import { useApi} from "@/api/ApiContext";
import { Event } from "@/types/events.ts";
import { types } from "@/types/core";
import { LocalPlayerContext, TrackSource } from "../types/playcontext";

export function LocalPlayer() {
  const audioRef = useRef<HTMLAudioElement>(null);
  const api = useApi();


    // FIXME: NEEDS OWN TYPE FOR ALL OF THIS. THIS COMPONENT IS IN CONTROL
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

    useEffect(() => {
        const unsubscribePlayerLocalStartContext = api.eventSubscribe(Event.PlayerLocalStartContext, (lctx: LocalPlayerContext) => {
            setCtx((prev) => ({
                    ...prev,
                    type: lctx.type,
                    ref_id: lctx.ref_id,
                    tracks: lctx.tracks,
                    track_order: lctx.tracks.map((_, i) => i),
                }));
            setPb((prev) => ({
                ...prev,
                playing: true,
                track_index: lctx.track_index,
            }));
            }
        );

        const unsubscribePlayPause = api.eventSubscribe(
            Event.PlayerLocalPlayPause,
            () => {
            setPb((prev) => ({
                ...prev,
                playing: !prev.playing,
            }));

            }
        );

        const unsubscribePlayerLocalNextTrack = api.eventSubscribe(Event.PlayerLocalNextTrack, () => {
                nextTrack();
            }
        );

        return () => {
            unsubscribePlayerLocalStartContext?.();
            unsubscribePlayPause?.();
            unsubscribePlayerLocalNextTrack?.();
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
        setCtx((prev) => ({
            ...prev,
            ref_id: "",
            tracks: [],
            track_order: [],
            queue: [],
        }))
        setPb((prev) => ({
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
        setPb((prev) => {
            return {
            ...prev,
            track_index: ntid,
        }});
    }


    return (
        <audio
        ref={audioRef}
            onEnded={() => {
                api.emitEvent(Event.PlayerLocalNextTrack);
            }}
        />
    )
;
}
