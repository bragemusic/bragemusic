// import { useEffect, useState } from "react";
// import {  types } from "../../wailsjs/go/models";
// import { GetPlayerState, Start } from "../../wailsjs/go/app/App";
// import { useHotkeys } from "react-hotkeys-hook";
// import {
//   Table,
//   TableBody,
//   TableCell,
//   TableColumn,
//   TableHeader,
//   TableRow,
// } from "@heroui/table";
// import { PauseCircle, PlayCircle } from "lucide-react";
// import clsx from "clsx";
// import { Chip, useOverlayState } from "@heroui/react";
// import { makeTimestamp } from "@/util/functions";
// import { StarRating } from "./StarRating";
// import { EventsOn } from "../../wailsjs/runtime/runtime";
// import { Event } from "@/types/events";
// import { TrackOptionsButton } from "./TrackOptionsButton";
// import EditModal from "./EditModal";
// import { Input, NumberInput, Textarea } from "@/primitives/Inputs";
// import { useApi, usePlayerApi } from "@/api/ApiContext";
// import { PlayContextType } from "@/types/playcontext";
// import PlaylistSelector from "./PlaylistSelector";
// import { LikeHeart } from "./LikeHeart";
// import { PlayerState, toPlayerState } from "@/models/PlayerState";

// type Layout = "full" | "playlist" | "likedtracks" | "compact";

// const LAYOUTS: Record<string, Layout> = {
//   FULL: "full",
//   PLAYLIST: "playlist",
//   LIKEDTRACKS: "likedtracks",
//   COMPACT: "compact",
// };

// type ColumnKey =
//   | "index"
//   | "play_symbol"
//   | "title"
//   | "media"
//   | "artist"
//   | "album"
//   | "length"
//   | "plays"
//   | "rating"
//   | "favorite"
//   | "album_options"
//   | "likedtracks_options"
//   | "plist_options";

// const COLUMNS_BY_LAYOUT: Record<Layout, ColumnKey[]> = {
//   full: [
//     "index",
//     "title",
//     "media",
//     "artist",
//     "album",
//     "length",
//     "plays",
//     "rating",
//     "favorite",
//     "album_options",
//   ],
//   playlist: [
//     "play_symbol",
//     "title",
//     "media",
//     "artist",
//     "album",
//     "length",
//     "plays",
//     "rating",
//     "favorite",
//     "plist_options",
//   ],
//   likedtracks: [
//     "play_symbol",
//     "title",
//     "media",
//     "artist",
//     "album",
//     "length",
//     "plays",
//     "rating",
//     "favorite",
//     "likedtracks_options",
//   ],
//   compact: ["index", "title", "length"],
// };

// interface TracksTableOldProps {
//   tracks: types.TrackDetailed[];
//   parent_id: string | undefined;
//   parent_type: PlayContextType;
// }

// export const TracksTableOld: React.FC<TracksTableOldProps> = ({
//   tracks,
//   parent_id,
//   parent_type,
// }) => {
//   const [selectedTrackIdx, setSelectedTrackIdx] = useState<number | null>(null);
//   const [selectedEditTrack, setSelectedEditTrack] =
//     useState<types.TrackDetailed | null>(null);
//   const [selectedPlistTrack, setSelectedPlistTrack] =
//     useState<types.TrackDetailed | null>(null);
//   const [numberOfDiscs, setNumberOfDiscs] = useState(1);
//   const [playCtx, setPlayCtx] = useState<PlayerState | null>(null);
//   const [isPlaying, setIsPlaying] = useState<boolean>();
//   const [layout, setLayout] = useState(LAYOUTS.FULL);

//   const editState = useOverlayState();
//   const plistState = useOverlayState();
//   const api = useApi();
//   const playerApi = usePlayerApi();

//   useEffect(() => {
//     if (parent_type === PlayContextType.Album) {
//       setLayout(LAYOUTS.FULL);
//     } else if (parent_type === PlayContextType.Playlist) {
//       setLayout(LAYOUTS.PLAYLIST);
//     } else if (parent_type === PlayContextType.LikedTracks) {
//       setLayout(LAYOUTS.LIKEDTRACKS);
//     }
//   }, []);

//   useEffect(() => {
//     const fetchPlayCtx = async () => {
//       const pc = await GetPlayerState();
//       setPlayCtx(toPlayerState(pc));
//     };
//     fetchPlayCtx();
//   }, []);

//   useEffect(() => {
//     const maxDisc = Math.max(...tracks.map((t) => t.disc_number || 1));
//     setNumberOfDiscs(maxDisc);
//   }, [tracks]);

//   useEffect(() => {
//     const offPlayCtx = EventsOn(
//       Event.PlayerContextChange,
//       (pc: types.PlayContext) => {
//         setPlayCtx(prev => {
//           const next = new PlayerState();

//           if (prev) {
//             Object.assign(next, prev);
//           }

//           next.context = pc;
//           return next;
//         });
//       },
//     );

//     const offPlaybackState = EventsOn(
//       Event.PlayerPlaybackChange,
//       (pc: types.PlaybackState) => {
//         setPlayCtx(prev => {
//           const next = new PlayerState();

//           if (prev) {
//             Object.assign(next, prev);
//           }

//           next.playback = pc;
//           return next;
//         });
//       },
//     );

//     const offIsPlaying = EventsOn("isPlaying", (p: boolean) => {
//       setIsPlaying(p);
//     });

//     return () => {
//       offPlayCtx();
//       offPlaybackState;
//       offIsPlaying();
//     };
//   }, []);

//   const selectPreviousTrack = () => {
//     if (selectedTrackIdx == null) {
//       setSelectedTrackIdx(0);
//       return;
//     }

//     let newIdx = selectedTrackIdx - 1;
//     if (newIdx < 0) {
//       newIdx = tracks.length - 1;
//     }

//     setSelectedTrackIdx(newIdx);
//   };

//   const selectNextTrack = () => {
//     if (selectedTrackIdx == null) {
//       setSelectedTrackIdx(0);
//       return;
//     }

//     let newIdx = selectedTrackIdx + 1;
//     if (newIdx >= tracks.length) {
//       newIdx = 0;
//     }

//     setSelectedTrackIdx(newIdx);
//   };

//   const startSelectedTrack = () => {
//     if (
//       selectedTrackIdx != null &&
//       tracks[selectedTrackIdx]?.media_file?.id &&
//       parent_id
//     ) {
//       Start(parent_type, parent_id, selectedTrackIdx);
//     }
//   };

//   useHotkeys("up", () => {
//     selectPreviousTrack();
//   });

//   useHotkeys("down", () => {
//     selectNextTrack();
//   });

//   useHotkeys("enter", () => {
//     startSelectedTrack();
//   });

//   const trackIsPlaying = (
//     pc: PlayerState | null,
//     t: types.TrackDetailed,
//   ): boolean => {
//     return (
//       pc?.context.ref_id == parent_id &&
//       pc?.current_track?.id === t.id
//     );
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

//   const onTrackEdit = (fd: FormData) => {
//     if (selectedEditTrack === null) {
//       return;
//     }
//     const source = {
//       musicbrainz_id: fd.get("musicbrainz_id") || undefined,
//       title: String(fd.get("title") ?? ""),
//       comment: fd.get("comment") || undefined,
//       disc_number: Number(fd.get("disc_number")),
//       track_number: Number(fd.get("track_number")),
//       album_id: parent_id,
//     };

//     const trackData = types.TrackUpdate.createFrom(source);
//     api.updateTrack(selectedEditTrack.id, trackData);
//   };

//   const onAddToPlaylist = (plistID: string) => {
//     if (parent_id === undefined || selectedPlistTrack === null) {
//       return;
//     }
//     api.addPlaylistTrack(plistID, parent_id.toString(), selectedPlistTrack.id);
//   };

//   const table = (
//     <>
//       {parent_type == PlayContextType.Album && (
//         <>
//           <PlaylistSelector
//             state={plistState}
//             onSubmit={onAddToPlaylist}
//           />
//           <EditModal
//             state={editState}
//             header={`Edit Track`}
//             onSubmit={onTrackEdit}
//           >
//             {selectedEditTrack != null && (
//               <>
//                 <Input
//                   name="title"
//                   label="Name"
//                   defaultValue={selectedEditTrack.title}
//                   required
//                 />
//                 <Input
//                   name="musicbrainz_id"
//                   label="MusicBrainz ID"
//                   defaultValue={selectedEditTrack.musicbrainz_id}
//                 />
//                 {
//                   // FIXME add track artists here
//                 }
//                 <NumberInput
//                   name="disc_number"
//                   label="Disc"
//                   defaultValue={selectedEditTrack.disc_number}
//                   required
//                 />
//                 <NumberInput
//                   name="track_number"
//                   label="Track"
//                   defaultValue={selectedEditTrack.track_number}
//                   required
//                 />
//                 <Textarea
//                   name="comment"
//                   label="Comment"
//                   defaultValue={selectedEditTrack.comment}
//                 />
//               </>
//             )}
//           </EditModal>
//         </>
//       )}
//       <Table
//         classNames={{ th: "bg-background" }}
//         removeWrapper
//         aria-label="Example static collection table"
//       >
//         <TableHeader>
//           {COLUMNS_BY_LAYOUT[layout].map((col) => {
//             switch (col) {
//               case "index":
//                 return (
//                   <TableColumn key="index" width={1} aria-label="index">
//                     {" "}
//                   </TableColumn>
//                 );

//               case "play_symbol":
//                 return (
//                   <TableColumn key="play_symbol" width={1}>
//                     &nbsp;
//                   </TableColumn>
//                 );
//               case "title":
//                 return <TableColumn key="title">NAME</TableColumn>;
//               case "media":
//                 return (
//                   <TableColumn key="media" width={1}>
//                     &nbsp;
//                   </TableColumn>
//                 );
//               case "artist":
//                 return <TableColumn key="artist">ARTIST</TableColumn>;
//               case "album":
//                 return <TableColumn key="album">ALBUM</TableColumn>;
//               case "length":
//                 return <TableColumn key="length">LENGTH</TableColumn>;
//               case "plays":
//                 return <TableColumn key="plays">PLAYS</TableColumn>;
//               case "rating":
//                 return (
//                   <TableColumn key="rating" width={1}>
//                     &nbsp;
//                   </TableColumn>
//                 );
//               case "favorite":
//                 return (
//                   <TableColumn key="favorite" width={1}>
//                     &nbsp;
//                   </TableColumn>
//                 );
//               case "album_options":
//                 return (
//                   <TableColumn key="options" width={1}>
//                     &nbsp;
//                   </TableColumn>
//                 );
//               case "plist_options":
//                 return (
//                   <TableColumn key="options" width={1}>
//                     &nbsp;
//                   </TableColumn>
//                 );
//               case "likedtracks_options":
//                 return (
//                   <TableColumn key="options" width={1}>
//                     &nbsp;
//                   </TableColumn>
//                 );
//               default:
//                 throw new Error(`Unhandled column: ${col}`);
//             }
//           })}
//         </TableHeader>
//         <TableBody>
//           {tracks.map((t, idx) =>
//             renderTrackRow(
//               t,
//               idx,
//               COLUMNS_BY_LAYOUT[layout],
//               playCtx,
//               onAlbumTrackOptionsCallback,
//               onPlaylistTrackOptionsCallback,
//               onLikedTracksOptionsCallback,
//               parent_type,
//               isPlaying,
//               numberOfDiscs,
//               parent_id,
//               () => setSelectedTrackIdx(idx),
//               trackIsPlaying,
//             ),
//           )}
//         </TableBody>
//       </Table>
//     </>
//   );

//   const noTracks = (
//     <p>No tracks available</p>
//   )

//   return tracks.length > 0 ? table : noTracks
// };

// function renderTrackRow(
//   t: types.TrackDetailed,
//   idx: number,
//   columns: ColumnKey[],
//   playCtx: PlayerState | null,
//   onAlbumOptions: (idx: number, type: string) => void,
//   onPlaylistOptions: (idx: number, type: string) => void,
//   onLikedTracksOptions: (idx: number, type: string) => void,
//   parent_type: PlayContextType,
//   isPlaying?: boolean,
//   numberOfDiscs?: number,
//   parent_id?: string,
//   onSelect?: () => void,
//   trackIsPlaying?: (
//     pc: PlayerState | null,
//     t: types.TrackDetailed,
//   ) => boolean,
// ): React.ReactElement {
//   return (
//     <TableRow
//       key={t.id}
//       onClick={onSelect}
//       onDoubleClick={() => {
//         t.media_file && parent_id && Start(parent_type, parent_id, idx);
//       }}
//       className={clsx(
//         t.media_file
//           ? "cursor-pointer hover:bg-default/60 select-none"
//           : "text-content4",
//       )}
//     >
//       {columns.map((col) => {
//         switch (col) {
//           case "index":
//             return (
//               <TableCell key="index">
//                 {trackIsPlaying && trackIsPlaying(playCtx, t) ? (
//                   isPlaying ? (
//                     <PlayCircle color="#007AFF" size={20} />
//                   ) : (
//                     <PauseCircle color="#007AFF" size={20} />
//                   )
//                 ) : numberOfDiscs && numberOfDiscs > 1 ? (
//                   `${t.disc_number}:${t.track_number}`
//                 ) : (
//                   t.track_number
//                 )}
//               </TableCell>
//             );

//           case "play_symbol":
//             return (
//               <TableCell key="play_symbol">
//                 {trackIsPlaying && trackIsPlaying(playCtx, t) ? (
//                   isPlaying ? (
//                     <PlayCircle color="#007AFF" size={20} />
//                   ) : (
//                     <PauseCircle color="#007AFF" size={20} />
//                   )
//                 ) : (
//                   " "
//                 )}
//               </TableCell>
//             );

//           case "title":
//             return <TableCell key="title">{t.title || " "}</TableCell>;

//           case "media":
//             return (
//               <TableCell key="media" className="flex flex-grow gap-0.5 items-center h-full min-w-20 lg:min-w-40">
//                 {t.media_file ? (
//                   <>
//                     <Chip
//                       className="hidden md:flex"
//                       color="accent"
//                       size="sm"
//                       variant="secondary"
//                     >
//                       {Math.round(t.media_file.bitrate / 1000)} Kbps
//                     </Chip>
//                     <Chip
//                       className="hidden lg:flex"
//                       color="accent"
//                       size="sm"
//                       variant="secondary"
//                     >
//                       {t.media_file.sample_rate / 1000} kHz
//                     </Chip>
//                     <Chip
//                       className="hidden xl:flex"
//                       color="accent"
//                       size="sm"
//                       variant="secondary"
//                     >
//                       {t.media_file.codec}
//                     </Chip>
//                   </>
//                 ) : (
//                   " "
//                 )}
//               </TableCell>
//             );

//           case "artist":
//             return (
//               <TableCell key="artist">
//                 {t.artist_names?.join(", ") || " "}
//               </TableCell>
//             );

//           case "album":
//             return <TableCell key="album">{t.album_name || " "}</TableCell>;

//           case "length":
//             return (
//               <TableCell key="length">
//                 {t.media_file ? makeTimestamp(t.media_file.duration_ms) : " "}
//               </TableCell>
//             );

//           case "plays":
//             return <TableCell key="plays">{t.play_count ?? " "}</TableCell>;

//           case "rating":
//             return (
//               <TableCell key="rating" className="hidden xl:table-cell">
//                 <StarRating
//                   rating={t.rating}
//                   user_rating={t.user_rating}
//                   track_id={t.id}
//                 />
//               </TableCell>
//             );

//           case "favorite":
//             return (
//               <TableCell key="favorite" className="hidden lg:table-cell">
//                 <LikeHeart track_id={t.id} liked={t.liked}/>
//               </TableCell>
//             );

//           case "album_options":
//             return (
//               <TableCell key="options">
//                 <TrackOptionsButton
//                   trackIdx={idx}
//                   parentType="album"
//                   onCallback={onAlbumOptions}
//                 />
//               </TableCell>
//             );

//           case "plist_options":
//             return (
//               <TableCell key="options">
//                 <TrackOptionsButton
//                   trackIdx={idx}
//                   parentType="playlist"
//                   onCallback={onPlaylistOptions}
//                 />
//               </TableCell>
//             );

//           case "likedtracks_options":
//             return (
//               <TableCell key="options">
//                 <TrackOptionsButton
//                   trackIdx={idx}
//                   parentType="liked_tracks"
//                   onCallback={onLikedTracksOptions}
//                 />
//               </TableCell>
//             );

//           default:
//             return <TableCell key={col}> </TableCell>;
//         }
//       })}
//     </TableRow>
//   );
// }
