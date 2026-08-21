import { useEffect, useState } from "react";
import {
  Chip,
  ChipProps,
  cn,
  useOverlayState,
} from "@heroui/react";
import { Disc, Pause, Play } from "lucide-react";
import React from "react";

import { types } from "@/types/core";

import { TrackOptionsButton } from "./TrackOptionsButton";
import { LikeHeart } from "./LikeHeart";
import { StarRating } from "./StarRating";
import EditTrackModal from "./EditTrackModal";
import PlaylistSelector from "./PlaylistSelector";
import { Image } from "./Image";

import { usePlayerState } from "@/store/playerStore";
import { PlayContextType } from "@/types/playcontext";
import { getArtistsString, getTrackQualityLabel, makeTimestamp, TrackQualityLabel } from "@/util/functions";
import { useApi, usePlayerApi } from "@/api/ApiContext";
import { albumImageLink } from "@/util/images";
import { useMediaQuery } from "react-responsive";
import { mqMobile } from "@/config/config";


type TrackTableCellProps = React.ComponentProps<"div">;
const TrackTableCell: React.FC<TrackTableCellProps> = ({
  className,
  children,
  ...props
}) => {
  return (
    <div
      {...props}
      className={cn(
        "border-none text-inherit whitespace-nowrap truncate",
        className,
      )}
    >
      {children}
    </div>
  );
};

const qualityColors: Record<TrackQualityLabel, ChipProps["color"]> = {
  Lossy: "danger",
  CD: "default",
  "HI-FI": "accent",
  Studio: "success",
} as const;

export type ColumnType =
  | "index"
  | "play_symbol"
  | "album_cover"
  | "title"
  | "title_artist"
  | "media"
  | "artist"
  | "album"
  | "length"
  | "plays"
  | "rating"
  | "like"
  | "album_options"
  | "plist_options"
  | "likedtracks_options";

type ColumnDefinition = {
  name: string;
  width: string;
  renderCell: (track: types.TrackDetailed, ctx: CellContext) => React.ReactNode;
};

interface CellContext {
  numberOfDiscs: number;
  trackIdx: number;
  trackIsPlaying: boolean;
  playerIsPlaying: boolean;
  onAlbumTrackOptionsCallback: (idx: number, type: string) => void;
  onPlaylistTrackOptionsCallback: (idx: number, type: string) => void;
  onLikedTracksOptionsCallback: (idx: number, type: string) => void;
}

const ColMap: Record<ColumnType, ColumnDefinition> = {
  index: {
    name: "",
    width: "50px",
    renderCell: (track, ctx) => (
      <TrackTableCell className="pl-2">
        {ctx.trackIsPlaying ? (
          ctx.playerIsPlaying ? (
            <Play className="fill-inherit stroke-inherit" size={18} />
          ) : (
            <Pause className="fill-inherit stroke-inherit" size={18} />
          )
        ) : ctx.numberOfDiscs > 1 ? (
          `${track.disc_number}:${track.track_number}`
        ) : (
          track.track_number
        )}
      </TrackTableCell>
    ),
  },

  play_symbol: {
    name: "",
    width: "50px",
    renderCell: (_, ctx) => (
      <TrackTableCell className="pl-2">
        {ctx.trackIsPlaying ? (
          ctx.playerIsPlaying ? (
            <Play className="fill-inherit stroke-inherit" size={18} />
          ) : (
            <Pause className="fill-inherit stroke-inherit" size={18} />
          )
        ) : (
          ""
        )}
      </TrackTableCell>
    ),
  },

  album_cover: {
    name: "",
    width: "70px",
    renderCell: (track, ctx) => (
      <TrackTableCell className="py-1.5 pl-2">
        <div className="relative w-[40px] h-[40px]">
          <Image
            fallbackIcon={Disc}
            height={40}
            radius="sm"
            src={albumImageLink(track.album_id ? track.album_id : "", 320)}
            width={40}
          />

          {ctx.trackIsPlaying && (
            <div className="flex absolute inset-0 justify-center items-center rounded-sm bg-black/40">
              {ctx.playerIsPlaying ? (
                <Pause className="fill-white stroke-white" size={18} />
              ) : (
                <Play className="fill-white stroke-white" size={18} />
              )}
            </div>
          )}
        </div>
      </TrackTableCell>
    ),
  },

  title: {
    name: "Name",
    width: "1fr",
    renderCell: (track) => (
      <TrackTableCell className={"truncate"}>
        {track.title}
      </TrackTableCell>
    ),
  },

  title_artist: {
    name: "Name",
    width: "1fr",
    renderCell: (track) => (
      <TrackTableCell className={"truncate"}>
        <div className="flex flex-col">
          <p>{track.title}</p>
          <p className="text-xs text-foreground/50">
            {getArtistsString(track.artists)}
          </p>
        </div>
      </TrackTableCell>
    ),
  },

  media: {
    name: "",
    width: "70px",
    renderCell: (track) => (
      <TrackTableCell className="flex gap-0.5 justify-end items-center pr-4 h-full">
        {track.media_file ? (
            <Chip
              color={qualityColors[getTrackQualityLabel(track.media_file)]}
              size="sm"
              variant="soft"
            >
              {getTrackQualityLabel(track.media_file)}
            </Chip>
        ) : (
          " "
        )}
      </TrackTableCell>
    ),
  },

  artist: {
    name: "Artist",
    width: "250px",
    renderCell: (track) => (
      <TrackTableCell className="pr-2 truncate">{getArtistsString(track.artists)}</TrackTableCell>
    ),
  },

  album: {
    name: "Album",
    width: "250px",
    renderCell: (track) => (
      <TrackTableCell className="pr-2 truncate">
        {track.album_name}
      </TrackTableCell>
    ),
  },

  length: {
    name: "Length",
    width: "70px",
    renderCell: (track) => (
      <TrackTableCell>
        {track.media_file ? makeTimestamp(track.media_file.duration_ms) : " "}
      </TrackTableCell>
    ),
  },

  plays: {
    name: "Plays",
    width: "50px",
    renderCell: (track) => (
      <TrackTableCell>{track.play_count ?? ""}</TrackTableCell>
    ),
  },

  rating: {
    name: "",
    width: "120px",
    renderCell: (track) => (
      <TrackTableCell>
        <StarRating
          rating={track.rating}
          track_id={track.id}
          user_rating={track.user_rating}
        />
      </TrackTableCell>
    ),
  },

  like: {
    name: "",
    width: "30px",
    renderCell: (track) => (
      <TrackTableCell>
        <LikeHeart liked={track.liked} track_id={track.id} />
      </TrackTableCell>
    ),
  },

  album_options: {
    name: "",
    width: "40px",
    renderCell: (_, ctx) => (
      <TrackTableCell className={"py-1"}>
        <TrackOptionsButton
          parentType="album"
          trackIdx={ctx.trackIdx}
          onCallback={ctx.onAlbumTrackOptionsCallback}
        />
      </TrackTableCell>
    ),
  },

  plist_options: {
    name: "",
    width: "40px",
    renderCell: (_, ctx) => (
      <TrackTableCell className={"py-1"}>
        <TrackOptionsButton
          parentType="playlist"
          trackIdx={ctx.trackIdx}
          onCallback={ctx.onPlaylistTrackOptionsCallback}
        />
      </TrackTableCell>
    ),
  },

  likedtracks_options: {
    name: "",
    width: "40px",
    renderCell: (_, ctx) => (
      <TrackTableCell className={"py-1"}>
        <TrackOptionsButton
          parentType="liked_tracks"
          trackIdx={ctx.trackIdx}
          onCallback={ctx.onLikedTracksOptionsCallback}
        />
      </TrackTableCell>
    ),
  },
};

interface TracksTableProps {
  tracks: types.TrackDetailed[];
  columns: ColumnType[];
  parent_id: string | undefined;
  parent_type: PlayContextType;
  alteringRowColor?: boolean;
  visibleCount?: number;
  virtualize?: boolean;
}

export const TracksTable: React.FC<TracksTableProps> = ({
  tracks,
  columns,
  parent_id,
  parent_type,
  alteringRowColor = true,
  visibleCount,
}) => {
  const [numberOfDiscs, setNumberOfDiscs] = useState(1);
  const [selectedEditTrack, setSelectedEditTrack] =
    useState<types.TrackDetailed | null>(null);
  const [selectedPlistTrack, setSelectedPlistTrack] =
    useState<types.TrackDetailed | null>(null);
  // const [isLoading, setIsLoading] = useState(false);

  const currentTrackId = usePlayerState((s) => s.player?.current_track?.id);
  const isPlaying = usePlayerState((s) => s.player?.playback?.playing);
  const playerApi = usePlayerApi();
  const api = useApi();
  const editState = useOverlayState();
  const plistState = useOverlayState();

  const isMobile = useMediaQuery({ query: mqMobile });

  const visibleTracks =
    visibleCount == null ? tracks : tracks.slice(0, visibleCount);

  useEffect(() => {
    const maxDisc = tracks.length
      ? Math.max(...tracks.map((t) => t.disc_number || 1))
      : 1;

    setNumberOfDiscs(maxDisc);
  }, [tracks]);

  const onAddToPlaylist = (plistID: string) => {
    if (parent_id === undefined || selectedPlistTrack === null) {
      return;
    }
    api.addPlaylistTrack(plistID, parent_id.toString(), selectedPlistTrack.id);
  };

  const onAlbumTrackOptionsCallback = (
    trackIdx: number,
    callbackType: string,
  ) => {
    if (callbackType === "edit") {
      setSelectedEditTrack(tracks[trackIdx]);
      editState.setOpen(true);

      return;
    }

    if (callbackType === "add_to_playlist") {
      setSelectedPlistTrack(tracks[trackIdx]);
      plistState.setOpen(true);

      return;
    }

    if (callbackType === "add_to_queue") {
      const t = tracks[trackIdx];

      if (t.album_id === undefined) {
        return;
      }
      playerApi.addTrackToQueue(t.id, t.album_id);

      return;
    }
  };

  const onPlaylistTrackOptionsCallback = (
    trackIdx: number,
    callbackType: string,
  ) => {
    if (callbackType === "add_to_queue") {
      const t = tracks[trackIdx];

      if (t.album_id === undefined) {
        return;
      }
      playerApi.addTrackToQueue(t.id, t.album_id);

      return;
    }

    if (callbackType === "delete") {
      const t = tracks[trackIdx];

      if (t.context_id === undefined) {
        return;
      }
      api.deletePlaylistTrack(t.context_id);

      return;
    }
  };

  const onLikedTracksOptionsCallback = (
    trackIdx: number,
    callbackType: string,
  ) => {
    if (callbackType === "add_to_queue") {
      const t = tracks[trackIdx];

      if (t.album_id === undefined) {
        return;
      }
      playerApi.addTrackToQueue(t.id, t.album_id);

      return;
    }
  };

  const cellContext: CellContext = {
    numberOfDiscs: numberOfDiscs,
    trackIdx: 0,
    trackIsPlaying: false,
    playerIsPlaying: false,
    onAlbumTrackOptionsCallback: onAlbumTrackOptionsCallback,
    onPlaylistTrackOptionsCallback: onPlaylistTrackOptionsCallback,
    onLikedTracksOptionsCallback: onLikedTracksOptionsCallback,
  };

  const gridcols = columns.map((c) => (ColMap[c].width)).join(" ")

  const content = (
    <div className="flex overflow-scroll flex-col flex-1 px-2 min-h-0">
      <div
          className={`grid items-center py-2 font-semibold uppercase text-[10px] text-foreground/50`}
          style={{
            gridTemplateColumns: gridcols,
          }}
      >
        {columns.map((colType) => (
          <div id={`${colType}-head`}>
            {ColMap[colType].name}
          </div>
        ))}
      </div>
      {visibleTracks.map((track, idx) => {
        const rowContext: CellContext = {
          ...cellContext,
          trackIdx: idx,
          trackIsPlaying: currentTrackId === track.id,
          playerIsPlaying: isPlaying,
        };

        return (
          <div
            className={cn(
              `grid items-center py-1 text-sm`,
              track.media_file ? "fill-foreground stroke-foreground" : "text-foreground/40",
              track.media_file ? "hover:bg-accent hover:text-accent-foreground hover:fill-accent-foreground hover:stroke-accent-foreground" : "",
              alteringRowColor && idx % 2 == 1 ? "bg-surface" : "",
            )}
            style={{
              gridTemplateColumns: gridcols,
            }}
            id={track.id}
            onClick={() => {
              if (!track.media_file) {
                return;
              }
              if (!isMobile) {
                return;
              }
              if (track.media_file && parent_id) {
                playerApi.start(parent_type, parent_id, idx);
              }
            }}
            onDoubleClick={() => {
              if (!track.media_file) {
                return;
              }
              if (isMobile) {
                return;
              }
              if (track.media_file && parent_id) {
                playerApi.start(parent_type, parent_id, idx);
              }
            }}
          >
            {columns.map((colType) => (
                ColMap[colType].renderCell(track, rowContext)
            ))}
          </div>
        )})}
    </div>
  )

  return (
    <>
      {parent_type == PlayContextType.Album && (
        <>
          <PlaylistSelector state={plistState} onSubmit={onAddToPlaylist} />
          <EditTrackModal
            parent_id={parent_id}
            selectedEditTrack={selectedEditTrack}
            state={editState}
          />
        </>
      )}
      {content}
    </>
  );
};
