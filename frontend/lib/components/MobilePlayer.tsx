import { Music2, Pause, Play } from "lucide-react";
import { useRef, useState } from "react";
import { Image } from "./Image";
import { usePlayerState } from "@/store/playerStore.ts";
import { albumImageLink } from "@/util/images.tsx";
import { usePlayerApi } from "@/api/ApiContext";
import { DeviceButton } from "./DeviceButton";
import { Link } from "react-router-dom";
import { PlayContextType } from "@/types/playcontext";
import { getArtistsString, makeTimestamp } from "@/util/functions";
import { usePlayerProgress } from "@/hooks/usePlayerProgress";
import { PlayerControls } from "./PlayerControls";
import { QueueButton } from "./QueueButton";
import { AnimatePresence, motion } from "framer-motion";


export const MobilePlayer = () => {
    const [maximized, setMaximized] = useState(false);
    const playCtx = usePlayerState((s) => s.player);
    const progressMs = usePlayerProgress();
    const startX = useRef<number | null>(null);

    const playerApi = usePlayerApi();

    const handleTouchStart = (e: React.TouchEvent<HTMLDivElement>) => {
        startX.current = e.changedTouches[0].clientX;
    };

    const handleTouchEnd = (e: React.TouchEvent<HTMLDivElement>) => {
        if (startX.current === null) return;

        const endX = e.changedTouches[0].clientX;
        const diff = endX - startX.current;

        if (diff > 50) {
            playerApi.previousTrack()
        }

        if (diff < -50) {
            playerApi.nextTrack()
        }

        startX.current = null;
    };

    const trackLink = (): string => {
        if (!playCtx.context) {
        return ""
        }

        if (playCtx?.context.type == "album") {
            return `/albums/${playCtx?.current_track?.album_id}`;
        } else if (playCtx?.context.type == "top_tracks") {
            // return `/artists/${playCtx?.current_track?.artists?.length > 0 ? playCtx?.current_track?.artists[0].id : ""}`;
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
            <>
                <div
                    onTouchStart={handleTouchStart}
                    onTouchEnd={handleTouchEnd}
                    className="absolute z-20 px-4 pb-4 w-full bottom-17"
                >
                    <div className="flex gap-2 items-center p-1.5 px-4 w-full rounded-full shadow-md backdrop-blur-sm bg-surface/60 border-1 border-border">
                        <Image
                            fallbackIcon={Music2}
                            height={30}
                            src={albumImageLink(
                            playCtx?.current_track?.album_id as string,
                            320,
                            )}
                            width={30}
                        />
                        <div
                            onClick={() => {setMaximized(true)}}
                            className="flex overflow-hidden flex-col flex-grow justify-between h-full text-[12px]"
                        >
                            <p className="font-semibold leading-tight">
                                {playCtx?.current_track?.title}
                            </p>
                            <p className="leading-tight text-foreground/80 text-title">
                                {getArtistsString(playCtx?.current_track?.artists)}
                            </p>
                        </div>
                        <div
                            className="pr-2 stroke-default-foreground fill-default-foreground"
                            onClick={() => {playerApi.playPause()}}
                        >
                            {playCtx?.playback?.playing ? (
                            <Pause size={20} className="fill-inherit stroke-inherit" />
                            ) : (
                            <Play size={20} className="fill-inherit stroke-inherit" />
                            )}
                        </div>
                        <DeviceButton minimized={true}/>
                    </div>
                </div>

                <AnimatePresence>
                {maximized &&
                    <motion.div
                    initial={{ y: "100%" }}
                    animate={{ y: 0 }}
                    exit={{ y: "100%" }}
                    transition={{ type: "spring", stiffness: 240, damping: 26, opacity: {duration: 0.18}  }}
                    drag="y"
                    dragDirectionLock
                    dragConstraints={{ top: 0, bottom: 500 }}
                    dragElastic={0.12}
                    onDragEnd={(_e, info) => {
                        if (info.offset.y > 120 || info.velocity.y > 100) {
                        setMaximized(false);
                        }
                    }}
                        className="flex overflow-hidden overscroll-none absolute top-0 left-0 flex-col items-center pt-8 w-full h-full bg-background z-60 touch-none"
                    >
                        <div className="flex flex-col justify-end px-8 pb-8 w-full h-1/2 max-h-1/2">
                            <Image
                                fallbackIcon={Music2}
                                height={320}
                                src={albumImageLink(
                                playCtx?.current_track?.album_id as string,
                                640,
                                )}
                                width={320}
                                radius="xl"
                                customHeight={true}
                                className="shadow-xl aspect-square"
                            />
                        </div>
                        <div className="flex flex-col gap-2 justify-between px-8 pb-8 w-full h-1/2 max-h-1/2">
                            <div className="flex overflow-hidden flex-col justify-between py-1 text-lg">
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
                            <div className="flex gap-2 items-center w-full">
                                <span className="text-xs text-foreground">
                                    {makeTimestamp(progressMs)}
                                </span>
                                <div className="overflow-hidden w-full rounded-full border-border/70 border-1 bg-background h-[8px]">
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
                            <PlayerControls playCtx={playCtx} mobile={true} className="pt-8"/>
                            <div className="flex justify-between items-center">
                                <QueueButton playCtx={playCtx} size="lg" />
                                <DeviceButton size="lg" />
                            </div>
                        </div>
                    </motion.div>
                }
            </AnimatePresence>
         </>
        )
};
