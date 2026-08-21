// import { useEffect, useState } from "react";
// import { Chip, cn, Table, TableLayout, useOverlayState, Virtualizer } from "@heroui/react";
// import { Disc, Pause, Play } from "lucide-react";
// import React from "react";

// import { types } from "@/types/core";

// import { TrackOptionsButton } from "./TrackOptionsButton";
// import { LikeHeart } from "./LikeHeart";
// import { StarRating } from "./StarRating";
// import EditTrackModal from "./EditTrackModal";
// import PlaylistSelector from "./PlaylistSelector";
// import { Image } from "./Image";

// import { usePlayerState } from "@/store/playerStore";
// import { PlayContextType } from "@/types/playcontext";
// import { makeTimestamp } from "@/util/functions";
// import { useApi, usePlayerApi } from "@/api/ApiContext";
// import { albumImageLink } from "@/util/images";
// import { useMediaQuery } from "react-responsive";
// import { mqMobile } from "@/config/config";

// type TrackTableColumnProps = React.ComponentProps<typeof Table.Column>;
// const TrackTableCol: React.FC<TrackTableColumnProps> = ({
//   className,
//   ...props
// }) => {
//   return <Table.Column {...props} className={cn("bg-background", className)} />;
// };

// type TrackTableCellProps = React.ComponentProps<typeof Table.Cell>;
// const TrackTableCell: React.FC<TrackTableCellProps> = ({
//   className,
//   ...props
// }) => {
//   return (
//     <Table.Cell
//       {...props}
//       className={cn(
//         "border-none text-inherit whitespace-nowrap truncate",
//         className,
//       )}
//     />
//   );
// };

// export type ColumnType =
//   | "index"
//   | "play_symbol"
//   | "album_cover"
//   | "title"
//   | "title_artist"
//   | "media"
//   | "artist"
//   | "album"
//   | "length"
//   | "plays"
//   | "rating"
//   | "like"
//   | "album_options"
//   | "plist_options"
//   | "likedtracks_options";

// type ColumnDefinition = {
//   renderHead: (ctx: CellContext) => React.ReactNode;
//   renderCell: (track: types.TrackDetailed, ctx: CellContext) => React.ReactNode;
// };

// interface CellContext {
//   numberOfDiscs: number;
//   trackIdx: number;
//   trackIsPlaying: boolean;
//   playerIsPlaying: boolean;
//   onAlbumTrackOptionsCallback: (idx: number, type: string) => void;
//   onPlaylistTrackOptionsCallback: (idx: number, type: string) => void;
//   onLikedTracksOptionsCallback: (idx: number, type: string) => void;
// }

// const ColMap: Record<ColumnType, ColumnDefinition> = {
//   index: {
//     renderHead: () => <TrackTableCol className="pr-0 w-6 bg-background" />,
//     renderCell: (track, ctx) => (
//       <TrackTableCell className="pl-2 fill-foreground stroke-foreground">
//         {ctx.trackIsPlaying ? (
//           ctx.playerIsPlaying ? (
//             <Play className="fill-inherit stroke-inherit" size={18} />
//           ) : (
//             <Pause className="fill-inherit stroke-inherit" size={18} />
//           )
//         ) : ctx.numberOfDiscs > 1 ? (
//           `${track.disc_number}:${track.track_number}`
//         ) : (
//           track.track_number
//         )}
//       </TrackTableCell>
//     ),
//   },

//   play_symbol: {
//     renderHead: () => (
//       <TrackTableCol className="pr-0 w-11 min-w-11 bg-background" />
//     ),
//     renderCell: (_, ctx) => (
//       <TrackTableCell className="pl-2 fill-foreground stroke-foreground">
//         {ctx.trackIsPlaying ? (
//           ctx.playerIsPlaying ? (
//             <Play className="fill-inherit stroke-inherit" size={18} />
//           ) : (
//             <Pause className="fill-inherit stroke-inherit" size={18} />
//           )
//         ) : (
//           ""
//         )}
//       </TrackTableCell>
//     ),
//   },

//   album_cover: {
//     renderHead: () => <TrackTableCol className="pr-0 min-w-17 bg-background" />,
//     renderCell: (track, ctx) => (
//       <TrackTableCell className="py-1.5 pl-2">
//         <div className="relative w-[40px] h-[40px]">
//           <Image
//             fallbackIcon={Disc}
//             height={40}
//             radius="sm"
//             src={albumImageLink(track.album_id ? track.album_id : "", 320)}
//             width={40}
//           />

//           {ctx.trackIsPlaying && (
//             <div className="flex absolute inset-0 justify-center items-center rounded-sm bg-black/40">
//               {ctx.playerIsPlaying ? (
//                 <Pause className="fill-white stroke-white" size={18} />
//               ) : (
//                 <Play className="fill-white stroke-white" size={18} />
//               )}
//             </div>
//           )}
//         </div>
//       </TrackTableCell>
//     ),
//   },

//   title: {
//     renderHead: () => <TrackTableCol isRowHeader className={"w-full truncate max-w-0"}>Name</TrackTableCol>,
//     renderCell: (track) => <TrackTableCell className={"truncate max-w-20"}>{track.title}</TrackTableCell>,
//   },

//   title_artist: {
//     renderHead: () => <TrackTableCol isRowHeader className={"w-full truncate max-w-0"}>Name</TrackTableCol>,
//     renderCell: (track) => (
//       <TrackTableCell className={"truncate max-w-20"}>
//         <div className="flex flex-col">
//           <p>{track.title}</p>
//           <p className="text-xs text-foreground/50">{track.artist_names?.join(", ") || " "}</p>
//         </div>
//       </TrackTableCell>),
//   },

//   media: {
//     renderHead: () => <TrackTableCol />,
//     renderCell: (track) => (
//       <TrackTableCell className="flex gap-0.5 items-center h-full">
//         {track.media_file ? (
//           <>
//             <Chip color="accent" size="sm" variant="secondary">
//               {Math.round(track.media_file.bitrate / 1000)} Kbps
//             </Chip>
//             <Chip color="accent" size="sm" variant="secondary">
//               {track.media_file.sample_rate / 1000} kHz
//             </Chip>
//           </>
//         ) : (
//           " "
//         )}
//       </TrackTableCell>
//     ),
//   },

//   artist: {
//     renderHead: () => <TrackTableCol>Artist</TrackTableCol>,
//     renderCell: (track) => (
//       <TrackTableCell>{track.artist_names?.join(", ") || " "}</TrackTableCell>
//     ),
//   },

//   album: {
//     renderHead: () => <TrackTableCol>Album</TrackTableCol>,
//     renderCell: (track) => <TrackTableCell className="max-w-80 truncate">{track.album_name}</TrackTableCell>,
//   },

//   length: {
//     renderHead: () => <TrackTableCol>Length</TrackTableCol>,
//     renderCell: (track) => (
//       <TrackTableCell>
//         {track.media_file ? makeTimestamp(track.media_file.duration_ms) : " "}
//       </TrackTableCell>
//     ),
//   },

//   plays: {
//     renderHead: () => <TrackTableCol>Plays</TrackTableCol>,
//     renderCell: (track) => (
//       <TrackTableCell>{track.play_count ?? ""}</TrackTableCell>
//     ),
//   },

//   rating: {
//     renderHead: () => <TrackTableCol />,
//     renderCell: (track) => (
//       <TrackTableCell>
//         <StarRating
//           rating={track.rating}
//           track_id={track.id}
//           user_rating={track.user_rating}
//         />
//       </TrackTableCell>
//     ),
//   },

//   like: {
//     renderHead: () => <TrackTableCol />,
//     renderCell: (track) => (
//       <TrackTableCell>
//         <LikeHeart liked={track.liked} track_id={track.id} />
//       </TrackTableCell>
//     ),
//   },

//   album_options: {
//     renderHead: () => <TrackTableCol />,
//     renderCell: (_, ctx) => (
//       <TrackTableCell className={"py-1"}>
//         <TrackOptionsButton
//           parentType="album"
//           trackIdx={ctx.trackIdx}
//           onCallback={ctx.onAlbumTrackOptionsCallback}
//         />
//       </TrackTableCell>
//     ),
//   },

//   plist_options: {
//     renderHead: () => <TrackTableCol />,
//     renderCell: (_, ctx) => (
//       <TrackTableCell className={"py-1"}>
//         <TrackOptionsButton
//           parentType="playlist"
//           trackIdx={ctx.trackIdx}
//           onCallback={ctx.onPlaylistTrackOptionsCallback}
//         />
//       </TrackTableCell>
//     ),
//   },

//   likedtracks_options: {
//     renderHead: () => <TrackTableCol />,
//     renderCell: (_, ctx) => (
//       <TrackTableCell className={"py-1"}>
//         <TrackOptionsButton
//           parentType="liked_tracks"
//           trackIdx={ctx.trackIdx}
//           onCallback={ctx.onLikedTracksOptionsCallback}
//         />
//       </TrackTableCell>
//     ),
//   },
// };

// interface TracksTableProps {
//   tracks: types.TrackDetailed[];
//   columns: ColumnType[];
//   parent_id: string | undefined;
//   parent_type: PlayContextType;
//   alteringRowColor?: boolean;
//   visibleCount?: number;
//   virtualize?: boolean;
// }

// export const TracksTable: React.FC<TracksTableProps> = ({
//   tracks,
//   columns,
//   parent_id,
//   parent_type,
//   alteringRowColor = true,
//   visibleCount,
//   virtualize = false,
// }) => {
//   const [numberOfDiscs, setNumberOfDiscs] = useState(1);
//   const [selectedEditTrack, setSelectedEditTrack] =
//     useState<types.TrackDetailed | null>(null);
//   const [selectedPlistTrack, setSelectedPlistTrack] =
//     useState<types.TrackDetailed | null>(null);
//   // const [isLoading, setIsLoading] = useState(false);

//   const currentTrackId = usePlayerState((s) => s.player?.current_track?.id);
//   const isPlaying = usePlayerState((s) => s.player?.playback?.playing);
//   const playerApi = usePlayerApi();
//   const api = useApi();
//   const editState = useOverlayState();
//   const plistState = useOverlayState();


//   const isMobile = useMediaQuery({ query: mqMobile })

//   const visibleTracks =
//     visibleCount == null ? tracks : tracks.slice(0, visibleCount);

//   useEffect(() => {
//     const maxDisc = tracks.length
//       ? Math.max(...tracks.map((t) => t.disc_number || 1))
//       : 1;

//     setNumberOfDiscs(maxDisc);
//   }, [tracks]);

//   const onAddToPlaylist = (plistID: string) => {
//     if (parent_id === undefined || selectedPlistTrack === null) {
//       return;
//     }
//     api.addPlaylistTrack(plistID, parent_id.toString(), selectedPlistTrack.id);
//   };

//   const onAlbumTrackOptionsCallback = (
//     trackIdx: number,
//     callbackType: string,
//   ) => {
//     if (callbackType === "edit") {
//       setSelectedEditTrack(tracks[trackIdx]);
//       editState.setOpen(true);

//       return;
//     }

//     if (callbackType === "add_to_playlist") {
//       setSelectedPlistTrack(tracks[trackIdx]);
//       plistState.setOpen(true);

//       return;
//     }

//     if (callbackType === "add_to_queue") {
//       const t = tracks[trackIdx];

//       if (t.album_id === undefined) {
//         return;
//       }
//       playerApi.addTrackToQueue(t.id, t.album_id);

//       return;
//     }
//   };

//   const onPlaylistTrackOptionsCallback = (
//     trackIdx: number,
//     callbackType: string,
//   ) => {
//     if (callbackType === "add_to_queue") {
//       const t = tracks[trackIdx];

//       if (t.album_id === undefined) {
//         return;
//       }
//       playerApi.addTrackToQueue(t.id, t.album_id);

//       return;
//     }

//     if (callbackType === "delete") {
//       const t = tracks[trackIdx];

//       if (t.context_id === undefined) {
//         return;
//       }
//       api.deletePlaylistTrack(t.context_id);

//       return;
//     }
//   };

//   const onLikedTracksOptionsCallback = (
//     trackIdx: number,
//     callbackType: string,
//   ) => {
//     if (callbackType === "add_to_queue") {
//       const t = tracks[trackIdx];

//       if (t.album_id === undefined) {
//         return;
//       }
//       playerApi.addTrackToQueue(t.id, t.album_id);

//       return;
//     }
//   };

//   const cellContext: CellContext = {
//     numberOfDiscs: numberOfDiscs,
//     trackIdx: 0,
//     trackIsPlaying: false,
//     playerIsPlaying: false,
//     onAlbumTrackOptionsCallback: onAlbumTrackOptionsCallback,
//     onPlaylistTrackOptionsCallback: onPlaylistTrackOptionsCallback,
//     onLikedTracksOptionsCallback: onLikedTracksOptionsCallback,
//   };

//   const content = (
//       <Table className="overflow-hidden h-full min-h-0" variant="secondary">
//         <Table.ScrollContainer className="overflow-y-auto h-full">
//     <Table.Content
//       aria-label="Table with custom cells"
//       selectionMode="single"
//     >
//       <Table.Header className={"bg-background"}>
//         {columns.map((colType) => (
//           <React.Fragment key={`${colType}-head`}>
//             {ColMap[colType].renderHead(cellContext)}
//           </React.Fragment>
//         ))}
//       </Table.Header>
//       <Table.Body>
//         {visibleTracks.map((track, idx) => {
//           const rowContext: CellContext = {
//             ...cellContext,
//             trackIdx: idx,
//             trackIsPlaying: currentTrackId === track.id,
//             playerIsPlaying: isPlaying,
//           };

//           return (
//             <Table.Row
//               key={track.id}
//               className={cn(
//                 isMobile? "" : "text-foreground fill-accent",
//                 isMobile? "" : "data-[selected=true]:bg-accent",
//                 isMobile? "" : "data-[selected=true]:text-accent-foreground",
//                 isMobile? "" : "data-[selected=true]:fill-accent-foreground",
//                 isMobile? "" : "data-[selected=true]:stroke-accent-foreground",
//                 isMobile? "" : "[&[data-selected=true]>td]:bg-accent",
//                 alteringRowColor && idx % 2 == 1 ? "bg-surface" : "",
//               )}
//               id={track.id}
//               onClick={() => {
//                 if (!isMobile) {
//                   return
//                 }
//                 if (track.media_file && parent_id) {
//                   playerApi.start(parent_type, parent_id, idx);
//                 }
//               }}
//               onDoubleClick={() => {
//                 if (isMobile) {
//                   return
//                 }
//                 if (track.media_file && parent_id) {
//                   playerApi.start(parent_type, parent_id, idx);
//                 }
//               }}
//             >
//               {columns.map((colType) => (
//                 <React.Fragment key={`${colType}-${track.id}`}>
//                   {ColMap[colType].renderCell(track, rowContext)}
//                 </React.Fragment>
//               ))}
//             </Table.Row>
//           );
//         })}
//       </Table.Body>
//     </Table.Content>
//         </Table.ScrollContainer>
//       </Table>
//   )

//   return (
//     <>
//       {parent_type == PlayContextType.Album && (
//         <>
//           <PlaylistSelector state={plistState} onSubmit={onAddToPlaylist} />
//           <EditTrackModal
//             parent_id={parent_id}
//             selectedEditTrack={selectedEditTrack}
//             state={editState}
//           />
//         </>
//       )}
//     {virtualize ?
//       <Virtualizer
//         layout={TableLayout}
//         layoutOptions={{
//           headingHeight: 42,
//           rowHeight: 52,
//         }}
//       >
//         {content}
//       </Virtualizer>
//       :
//      content
//     }
//     </>
//   );
// };
