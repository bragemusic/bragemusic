import { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import { ImageUp, Play, Shuffle, UserPen } from "lucide-react";
import { Button, useOverlayState } from "@heroui/react";

import { types } from "@/types/core";

import { AlbumButton } from "@/components/AlbumButton";
import BannerNew from "@/components/BannerNew";
import { artistImageLink } from "@/util/images";
import { useSession } from "@/session/SessionContext";
import { TracksTable } from "@/components/TracksTable";
import { PlayContextType } from "@/types/playcontext";
import { BannerButtonProps } from "@/primitives/Buttons";
import { useApi } from "@/api/ApiContext";
import { Event } from "@/types/events";
import { useMediaQuery } from "react-responsive";
import { mqMobile } from "@/config/config";
import { About } from "../components/About";
import { ArtistEditModal } from "../components/ArtistEditModal";
import { useBackdrop } from "@/hooks/BackgroundProvider";
import DefaultLayout from "@/layouts/default";

export default function ArtistPageNew({}: {}) {
  const [artist, setArtist] = useState<types.Artist>();
  const [albums, setAlbums] = useState<types.AlbumDetailed[]>([]);
  const [featuredAlbums, setFeaturedAlbums] = useState<types.AlbumDetailed[]>(
    [],
  );
  const [topTracks, setTopTracks] = useState<types.TrackDetailed[]>([]);
  const [topTracksOpen, setTopTracksOpen] = useState(false);

  const api = useApi();

  const params = useParams();
  const session = useSession();
  const isMobile = useMediaQuery({ query: mqMobile });

  const editState = useOverlayState();

  useBackdrop(artist && artistImageLink(artist, 320), { saturate: 1.9 });

  const loadData = () => {
    api.getArtist(params.artistID as string).then(setArtist);
    api.listAlbumsByArtist(params.artistID as string).then(setAlbums);
    api
      .listFeaturedAlbumsByArtist(params.artistID as string)
      .then(setFeaturedAlbums);
    api.getArtistTopTracks(params.artistID as string).then(setTopTracks);
  };

  useEffect(() => {
    loadData();
  }, [api, params.artistID]);

  useEffect(() => {
    const evUpdated = api.eventSubscribe(Event.EntitiesUpdated, () => {
      loadData();
    });

    return () => {
      evUpdated();
    };
  }, []);

  const bannerButtons:BannerButtonProps[] = [
    {
      icon: ImageUp,
      onUpload: (files) => {
        api.uploadArtistImage(params.artistID as string, files[0]);
      },
      isDisabled: !session.serverFullyAvailable,
      tooltip: "Upload Artist Image",
      type: "fileupload",
    },
    {
      icon: UserPen,
      onClick: () => {
        editState.setOpen(true);
      },
      isDisabled: !session.serverFullyAvailable,
      tooltip: "Edit Artist",
    },
  ]

  return (
    <DefaultLayout>
      <ArtistEditModal artist={artist} state={editState}/>
      {artist && (
        <BannerNew
          bannerType="artist"
          imageUrl={artistImageLink(artist, 1024)}
          buttons={bannerButtons}
          title={artist?.name}
          subtitle={`${artist.country}, active since ${artist.year_started}`}
        >
          <div className="flex gap-2 mt-6">
            <Button size="lg">
              <Play/>
              Play top tracks
            </Button>
            <Button size="lg" variant="tertiary" className="bg-background/0 border-1 border-foreground">
              <Shuffle/>
              Shuffle all
            </Button>
          </div>
        </BannerNew>
      )}

        <div className="flex gap-8 mt-8 w-full min-w-0">
          <div className="flex flex-col justify-end w-70">
            {artist?.description &&
              <About text={artist?.description}/>
            }
    {
            // <Surface className="p-2 rounded-xl border-border border-1">
            //   <p>hej</p>
            // </Surface>
    }
          </div>
          <div className="flex-grow">
            <h2 className="mb-2 text-3xl font-black font-display">Top Tracks</h2>
            {artist && (
              <div className="flex flex-col gap-2">
                <TracksTable
                  columns={
                    isMobile
                      ? ["album_cover", "title", "length", "plays"]
                      : ["album_cover", "title", "length", "plays"]
                  }
                  parent_id={artist.id.toString()}
                  parent_type={PlayContextType.TopTracks}
                  tracks={topTracks}
                  visibleCount={topTracksOpen ? 10 : 5}
                />
                <p
                  className="text-sm font-bold cursor-pointer hover:underline text-accent-deep"
                  onClick={() => {
                    setTopTracksOpen(!topTracksOpen);
                  }}
                >
                {topTracksOpen ? "Show 5 tracks" : "Show 10 tracks"}
                </p>
              </div>
            )}
          </div>
        </div>

        {albums.length > 0 && (
          <>
            <h2 className="mt-14 text-4xl font-black font-display">Albums</h2>
            <div className="grid grid-cols-2 gap-2 pt-4 sm:flex sm:flex-wrap sm:gap-8">
              {albums.map(function (a) {
                return <AlbumButton key={a.id} album={a} />;
              })}
            </div>
          </>
        )}

        {featuredAlbums.length > 0 && (
          <>
            <h2 className="mt-14 text-4xl font-black font-display">Appears On</h2>
            <div className="grid grid-cols-2 gap-2 pt-4 sm:flex sm:flex-wrap sm:gap-8">
              {featuredAlbums.map(function (a) {
                return <AlbumButton key={a.id} album={a} />;
              })}
            </div>
          </>
        )}
    </DefaultLayout>
  );
}
