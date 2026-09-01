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
    // const [isPlaying, setIsPlaying] = useState(false);

  const [currentTrack, setCurrentTrack] = useState(new types.TrackDetailed)


  // api.eventSubscribe(Event.PlayerLocalContextChange, (ctx: types.PlayContext) => {
  //   setCtx(ctx)
  //   // api.emitEvent(Event.PlayerContextChange, ctx)
  //   console.log("ctx", ctx)
  // });

  // api.eventSubscribe(Event.PlayerLocalPlaybackChange, (pb: types.PlaybackState) => {
  //   setPb(pb)
  //   // api.emitEvent(Event.PlaybackChange, pb)
  //   console.log("pb", pb)
  // });

    useEffect(() => {
        const unsubscribePlayerLocalStartContext = api.eventSubscribe(Event.PlayerLocalStartContext, (lctx: LocalPlayerContext) => {
            console.log(lctx)
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

        return () => {
            unsubscribePlayerLocalStartContext?.();
            unsubscribePlayPause?.();
        };
    }, [api]);

  const evalContextAndState = (ctx: types.PlayContext, pb: types.PlaybackState) => {
      console.log(ctx, pb)
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


  // useEffect(() => {
  //   console.log(playCtx)
  // }, [playCtx]);
  // local playback implementation

  return (
    <audio
      ref={audioRef}
    />
  )
;
}
