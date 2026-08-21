import { Button, cn } from "@heroui/react";
import {
  Play,
  Pause,
  SkipForward,
  SkipBack,
  Shuffle,
  Repeat1,
  Repeat,
  LucideIcon,
} from "lucide-react";

import { usePlayerApi } from "@/api/ApiContext.tsx";
import { PlayerState } from "../models/PlayerState";

interface PlayPauseButtonProps {
  mobile: boolean;
  isPlaying: boolean;
  onClick: () => void;
}

const PlayPauseButton: React.FC<PlayPauseButtonProps> = ({mobile, isPlaying, onClick}) => {
    return (
      <Button
        className={cn("text-white rounded-full bg-accent stroke-accent-foreground fill-accent-foreground", mobile ? "w-18 h-18" : "")}
        isIconOnly={true}
        size="lg"
        onPressEnd={onClick}
      >
        {isPlaying ? (
          <Pause className={cn("fill-inherit stroke-inherit", mobile ? "w-[40px] h-[40px]" : "")} />
        ) : (
          <Play className={cn("fill-inherit stroke-inherit", mobile ? "w-[40px] h-[40px]" : "")} />
        )}
      </Button>
    )
}

interface MdButtonProps {
  mobile: boolean;
  icon: LucideIcon;
  onClick: () => void;
}

const MdButton: React.FC<MdButtonProps> = ({mobile, icon:Icon, onClick}) => {
  return (
      <Button
        className={mobile ? "h-12 w-12 rounded-full" : ""}
        isIconOnly={true}
        size="md"
        variant="ghost"
        onPressEnd={onClick}
      >
        <Icon className={mobile ? "w-[28px] h-[28px]" : ""}/>
      </Button>
  )
}

interface PlayerControlsProps {
  playCtx: PlayerState;
  mobile: boolean;
  className?: string;
}

export const PlayerControls: React.FC<PlayerControlsProps> = ({playCtx, mobile, className=""}) => {
  const playerApi = usePlayerApi();

  return (
    <div className={cn("flex gap-4 items-center", mobile ? "justify-between w-full" : "", className)}>
      <Button
        isIconOnly={true}
        size="sm"
        variant="ghost"
        onPressEnd={() => {playerApi.toggleShuffle()}}
      >
        <Shuffle
          className={cn(
            playCtx?.playback?.shuffle ? "text-accent" : "text-content4",
          )}
          size={17}
        />
      </Button>
      <MdButton mobile={mobile} icon={SkipBack} onClick={() => {playerApi.previousTrack()}}/>
      <PlayPauseButton mobile={mobile} isPlaying={playCtx?.playback?.playing} onClick={() => {playerApi.playPause()}} />
      <MdButton mobile={mobile} icon={SkipForward} onClick={() => {playerApi.nextTrack()}}/>
      <Button
        isIconOnly={true}
        size="sm"
        variant="ghost"
        onPressEnd={() => {playerApi.nextRepeat()}}
      >
        {playCtx?.playback?.repeat === "one" ? (
          <Repeat1 className="text-accent" size={17} />
        ) : (
          <Repeat
            className={cn(
              playCtx?.playback?.repeat === "all"
                ? "text-accent"
                : "text-content4",
            )}
            size={17}
          />
        )}
      </Button>
    </div>
  );
};
