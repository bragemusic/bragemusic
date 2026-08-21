import {
  Music2,
} from "lucide-react";
import { useHotkeys } from "react-hotkeys-hook";
import { Link } from "react-router-dom";

import { Image } from "./Image";
import { QueueButton } from "./QueueButton.tsx";
import { DeviceButton } from "./DeviceButton.tsx";
import { PlayerControls } from "./PlayerControls.tsx";

import { albumImageLink } from "@/util/images.tsx";
import { makeTimestamp } from "@/util/functions.ts";
import { PlayContextType } from "@/types/playcontext.ts";
import { usePlayerApi } from "@/api/ApiContext.tsx";
import { usePlayerState } from "@/store/playerStore.ts";
import { usePlayerProgress } from "@/hooks/usePlayerProgress.ts";

export const Player = () => {
  // FIXME: Needs to add a function to call on mount to determine if isPlaying

  const playerApi = usePlayerApi();

  const playCtx = usePlayerState((s) => s.player);
  const progressMs = usePlayerProgress();

  useHotkeys("space", () => {
    playerApi.playPause();
  });

  const trackLink = (): string => {
    if (!playCtx.context) {
      return ""
    }

    if (playCtx?.context.type == "album") {
      return `/albums/${playCtx?.current_track?.album_id}`;
    } else if (playCtx?.context.type == "top_tracks") {
      // return `/artists/${playCtx?.current_track?.artist_ids}`;
      return `/artists/`;
    } else if (playCtx?.context.type == PlayContextType.Playlist) {
      return `/playlists/${playCtx.context.ref_id}`;
    } else if (playCtx?.context.type == PlayContextType.LikedTracks) {
      return `/liked-tracks`;
    }

    return "";
  };

  const artistIds = playCtx?.current_track?.artists?.map(a => a.id);
  const artistNames = playCtx?.current_track?.artists?.map(a => a.name);

  const duration = playCtx?.current_track?.media_file?.duration_ms ?? 0;

  const progressPercent =
    duration > 0 ? Math.min(1, Math.max(0, progressMs / duration)) : 0;

  return (
    <div className="flex justify-between select-none h-25 bg-surface border-t-1 border-border">
      <div className="flex gap-4 items-center p-4 w-1/3 text-foreground">
        <Link to={trackLink()}>
          <Image
            fallbackIcon={Music2}
            height={65}
            src={albumImageLink(
              playCtx?.current_track?.album_id as string,
              320,
            )}
            width={65}
          />
        </Link>
        <div className="flex overflow-hidden flex-col justify-between py-1 h-full text-sm">
          <Link className="font-bold hover:underline truncate" to={trackLink()}>
            {playCtx?.current_track?.title}
          </Link>
          <div>
            {artistIds?.map((artistId, index) => (
              <Link
                key={artistId}
                className="hover:underline text-foreground/80 text-title"
                to={`/artists/${artistId}`}
              >
                {artistNames?.[index]}
                {index < artistIds.length - 1 && ", "}
              </Link>
            ))}
          </div>
          <Link
            className="hover:underline text-foreground/80 truncate"
            to={`/albums/${playCtx?.current_track?.album_id}`}
          >
            {playCtx?.current_track?.album_name}
          </Link>
        </div>
      </div>

      <div className="flex flex-col justify-between items-center p-3 w-1/3 h-full">
        <PlayerControls playCtx={playCtx} mobile={false}/>
        <div className="flex gap-2 items-center w-full">
          <span className="text-xs text-foreground">
            {makeTimestamp(progressMs)}
          </span>
          <div className="overflow-hidden w-full rounded-full border-border/70 border-1 bg-background h-[6px]">
            <div
              className="h-full origin-left bg-accent"
              style={{ transform: `scaleX(${progressPercent})` }}
            />
          </div>
          <span className="text-xs text-foreground">
            {makeTimestamp(
              playCtx?.current_track?.media_file?.duration_ms
                ? playCtx?.current_track?.media_file?.duration_ms
                : 0,
            )}
          </span>
        </div>
      </div>

      <div className="flex gap-2 justify-end items-center p-4 w-1/3">
        <QueueButton playCtx={playCtx} />
        <DeviceButton />
      </div>
    </div>
  );
};
