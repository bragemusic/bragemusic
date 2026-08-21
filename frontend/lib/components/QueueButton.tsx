import { Button, ButtonProps, Popover, Tooltip } from "@heroui/react";
import { List, Music2 } from "lucide-react";
import { useState } from "react";
import { cn } from "@heroui/react";

import { types } from "@/types/core";

import { Image } from "./Image";

import { albumImageLink } from "@/util/images";
import { PlayerState } from "@/models/PlayerState";

interface TrackInfoProps {
  track: types.TrackDetailed;
}

const TrackInfo: React.FC<TrackInfoProps> = ({ track }) => {
  return (
    <div className="flex gap-2 items-center cursor-default select-none">
      <Image
        fallbackIcon={Music2}
        height={35}
        radius="sm"
        src={albumImageLink(track.album_id as string, 320)}
        width={35}
      />
      <div className="flex overflow-hidden flex-col justify-between py-1 h-full text-xs">
        <p className="font-bold truncate">{track.title}</p>
        <div>
          {track.artists &&
            track.artists?.map((artist) => (
              <p key={artist.id} className="text-foreground/90 text-title">
                {artist.name}
              </p>
            ))}
        </div>
      </div>
    </div>
  );
};

interface QueueButtonProps {
  playCtx: PlayerState;
  size?: ButtonProps["size"];
}

export const QueueButton: React.FC<QueueButtonProps> = ({ playCtx, size="sm" }) => {
  const [popoverOpen, setPopoverOpen] = useState(false);

  const isEnabled = (): boolean => {
    if (playCtx === null || !playCtx.context) {
      return false;
    }

    if (!playCtx.context.tracks && !playCtx.context.queue) {
      return false;
    }

    if (
      playCtx.context.tracks.length === 0 &&
      playCtx.context.queue.length === 0
    ) {
      return false;
    }

    return true;
  };

  const typeString = (): string => {
    if (playCtx === null || !playCtx.context) {
      return "";
    }

    if (playCtx.context.type === "liked_tracks") {
      return "liked tracks";
    }

    return playCtx.context.type;
  };

  const popoverContent = (
    <Popover.Content className="overflow-y-scroll w-[240px] max-h-[500px]">
      <Popover.Arrow />
      <Popover.Dialog>
        <Popover.Heading />
        <div className="py-2 px-1 w-full text-foreground">
          {playCtx.context?.queue?.length !== undefined &&
            playCtx?.context.queue.length > 0 && (
              <>
                <p className="pt-4 pb-2 font-bold text-small text-foreground">
                  Next in queue
                </p>
                <div className="flex flex-col gap-4">
                  {playCtx.context?.queue.map(function (t, idx) {
                    return <TrackInfo key={"queue" + idx} track={t} />;
                  })}
                </div>
              </>
            )}
          <p className="pt-4 pb-2 font-bold text-small text-foreground">
            Next in {typeString()}
          </p>
          <div className="flex flex-col gap-1">
            {isEnabled() &&
              playCtx?.context.track_order
                .slice(playCtx.playback.track_index + 1)
                .map(function (idx) {
                  return (
                    <TrackInfo
                      key={"tracks" + idx}
                      track={playCtx.context.tracks[idx]}
                    />
                  );
                })}
          </div>
        </div>
      </Popover.Dialog>
    </Popover.Content>
  );

  return (
    <Tooltip closeDelay={100} delay={300}>
      <span className="inline-flex">
        <Popover
          isOpen={popoverOpen}
          onOpenChange={(open) => {
            if (!open) {
              setPopoverOpen(false);

              return;
            }

            if (isEnabled()) {
              setPopoverOpen(true);
            }
          }}
        >
          <Popover.Trigger>
            <Button
              isDisabled={!isEnabled()}
              isIconOnly={true}
              size={size}
              variant="ghost"
            >
              <List
                className={cn(
                  isEnabled() ? "text-foreground" : "text-content4",
                )}
                size={20}
              />
            </Button>
          </Popover.Trigger>
          {popoverContent}
        </Popover>
      </span>
      <Tooltip.Content>
        <p>View play queue</p>
      </Tooltip.Content>
    </Tooltip>
  );
};
