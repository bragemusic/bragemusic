import { useApi } from "@/api/ApiContext";
import { TracksTable, ColumnType } from "@/components/TracksTable";
import { TracksFilterPane } from "@/components/TracksFilterPane";
import { mqMobile } from "@/config/config";
import PageTitleLayout from "@/layouts/page_title";
import { types } from "@/types/core";
import { PlayContextType } from "@/types/playcontext";
import { Button, Text, Drawer, useOverlayState } from "@heroui/react";
import { SlidersHorizontal, Sparkles } from "lucide-react";
import { memo, useEffect, useMemo, useState } from "react";
import { useMediaQuery } from "react-responsive";
import { SmartPlaylistModal } from "@/components/SmartPlaylistModal";

export default function TracksPage() {
    const [filter, setFilter] = useState<types.TrackFilter>({mood: {}} as types.TrackFilter)
    // const [tracks, setTracks] = useState<responses.ListPaginationPayload<types.TrackDetailed>>(new responses.ListPaginationPayload<types.TrackDetailed>)
    const [tracks1, setTracks1] = useState<types.TrackDetailed[]>([])
    const [_, setIsLoading] = useState(false);
    const [totalTracks, setTotalTracks] = useState(0);
    const [artists, setArtists] = useState<types.Artist[]>([]);

    const smartPlaylistState = useOverlayState();

    const api = useApi();
    const isMobile = useMediaQuery({ query: mqMobile })

    useEffect(() => {
        api.listArtists().then(setArtists);
    }, [api])

    useEffect(() => {
        console.log("filter")
        let cancelled = false;

        async function loadAll() {
            setIsLoading(true);
            setTracks1([]);

            // page 1 first
            const first = await api.filterTracks(filter, 1, 100);
            setTotalTracks(first.total_items)

            if (cancelled || first.items === undefined) return;

            // setTracks1(first.items);
            const all = [...first.items];

            const totalPages = first.total_pages;

            // load rest in background
            for (let page = 2; page <= totalPages; page++) {
                const res = await api.filterTracks(filter, page, 100);

                const resitems = res.items
                if (cancelled || !resitems) return;

                all.push(...resitems);
            }

            setTracks1(all)
            setIsLoading(false);

        }

  // const timeout = setTimeout(() => {
      loadAll();
  // }, 400)

  return () => {
      // clearTimeout(timeout)
      cancelled = true
  }
    }, [filter]);


const columns = useMemo<ColumnType[]>(() => {
  return isMobile
    ? [
        "album_cover",
        "title_artist",
        "plist_options",
      ]
    : [
        "album_cover",
        "title",
        "artist",
        "album",
        "length",
        "plays",
      ];
}, [isMobile]);

    const MemoTracksTable = memo(TracksTable);

  return (
    <PageTitleLayout title="Tracks" scroll={false}>
        <SmartPlaylistModal state={smartPlaylistState} artists={artists} initFilter={filter}/>
        <div className="flex flex-col flex-1 w-full min-h-0">
            <div className="flex gap-2">
                <Drawer>
                    <Button variant="secondary" className="mb-4"><SlidersHorizontal/> Filters</Button>
                    <Drawer.Backdrop>
                        <Drawer.Content placement={isMobile ? "top" : "right"}>
                            <Drawer.Dialog>
                                <Drawer.Header>
                                    <Drawer.Heading className="flex gap-2 justify-between items-center">
                                        <div className="flex gap-2 items-center">
                                            <SlidersHorizontal/ >
                                            <Text type="h5">Filters</Text>
                                        </div>
                                        <Text color="muted" type="body-sm">{totalTracks} tracks found</Text>
                                    </Drawer.Heading>
                                </Drawer.Header>
                                <Drawer.Body>
                                    <TracksFilterPane initFilter={filter} onApply={setFilter} artists={artists}/>
                                </Drawer.Body>
                            </Drawer.Dialog>
                        </Drawer.Content>
                    </Drawer.Backdrop>
                </Drawer>
                <Button onClick={() => {smartPlaylistState.setOpen(true)}}><Sparkles/>Create Smart Playlist</Button>
            </div>
              <MemoTracksTable
                columns={columns}
                parent_id="00000000-0000-0000-0000-000000000000"
                parent_type={PlayContextType.Filter}
                tracks={tracks1}
              />
        </div>
    </PageTitleLayout>
  );
}

                // <Drawer.Body>
                //     <Text type="h6"className="pt-4 pb-2">BPM</Text>
                //     <BPMSlider
                //         defaultUpperValue={filter.bpm?.upper}
                //         defaultLowerValue={filter.bpm?.lower}
                //         onChange={(vl, vu) => {
                //             if (vu === undefined) {
                //                 setFilter(prev => new types.TrackFilter({
                //                 ...prev,
                //                 bpm: undefined
                //                 }));
                //                 return
                //             }

                //             setFilter(prev => new types.TrackFilter({
                //             ...prev,
                //             bpm: {
                //                 lower: vl,
                //                 upper: vu,
                //             }
                //             }));
                //             return
                //         }}
                //     />
                //     <Text type="h6"className="pt-4 pb-2">Mood</Text>
                //     <div className="flex flex-col gap-4">
                //         <PercentSlider
                //             label="Aggressive"
                //             defaultValue={filter.mood.aggressive}
                //             onChange={(v) => {
                //                 setFilter(prev => new types.TrackFilter({
                //                     ...prev,
                //                     mood: {
                //                         ...prev.mood,
                //                         aggressive: v ? v / 100 : undefined,
                //                     }
                //                 }));
                //             }}
                //         />
                //         <PercentSlider
                //             label="Calm"
                //             defaultValue={filter.mood.calm}
                //             onChange={(v) => {
                //                 setFilter(prev => new types.TrackFilter({
                //                     ...prev,
                //                     mood: {
                //                         ...prev.mood,
                //                         calm: v ? v / 100 : undefined,
                //                     }
                //                 }));
                //             }}
                //         />
                //         <PercentSlider
                //             label="Happy"
                //             defaultValue={filter.mood.happy}
                //             onChange={(v) => {
                //                 setFilter(prev => new types.TrackFilter({
                //                     ...prev,
                //                     mood: {
                //                         ...prev.mood,
                //                         happy: v ? v / 100 : undefined,
                //                     }
                //                 }));
                //             }}
                //         />
                //         <PercentSlider
                //             label="Sad"
                //             defaultValue={filter.mood.sad}
                //             onChange={(v) => {
                //                 setFilter(prev => new types.TrackFilter({
                //                     ...prev,
                //                     mood: {
                //                         ...prev.mood,
                //                         sad: v ? v / 100 : undefined,
                //                     }
                //                 }));
                //             }}
                //         />
                //     </div>
                // </Drawer.Body>
                // <Drawer.Footer>
                // </Drawer.Footer>
