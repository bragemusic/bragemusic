import { useState, useEffect } from "react";

import { types } from "@/types/core";

import { useApi } from "@/api/ApiContext";
import BodyLayout from "@/layouts/BodyLayout";
import { TracksTable } from "@/components/TracksTable";
import { PlayContextType } from "@/types/playcontext";
import Banner from "@/components/Banner";
import { Event } from "@/types/events";
import { useMediaQuery } from "react-responsive";
import { mqMobile } from "@/config/config";

export default function LikedTracksPage({
  scrollRef,
}: {
  scrollRef: React.RefObject<HTMLElement>;
}) {
  const [tracks, setTracks] = useState<types.TrackDetailed[]>([]);

  const api = useApi();

  const isMobile = useMediaQuery({ query: mqMobile })

  const loadData = () => {
    api.listLikedTracks().then(setTracks);
  };

  useEffect(() => {
    loadData();
  }, [api]);

  useEffect(() => {
    const syncF = api.eventSubscribe(Event.EntitiesUpdated, () => {
      loadData();
    });

    return () => {
      syncF();
    };
  }, []);

  return (
    <>
        <Banner scrollRef={scrollRef} title={"Liked Tracks"}>
          <p className="pt-6">
            Here are all the tracks you like the most. Enjoy!
          </p>
        </Banner>
        <BodyLayout>
          <div className="pt-10">
            <TracksTable

              columns={isMobile ? [
                "album_cover",
                "title_artist",
                "likedtracks_options",
              ] : [
                "play_symbol",
                "title",
                "media",
                "artist",
                "album",
                "length",
                "plays",
                "rating",
                "like",
                "likedtracks_options",
              ]}
              parent_id="00000000-0000-0000-0000-000000000000"
              parent_type={PlayContextType.LikedTracks}
              tracks={tracks}
            />
          </div>
        </BodyLayout>
      </>
  );
}
