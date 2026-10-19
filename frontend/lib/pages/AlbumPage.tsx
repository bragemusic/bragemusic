import { useState, useEffect } from "react";
import { Link, useParams } from "react-router-dom";
import { ImageUp, Play, Shuffle, SquarePen } from "lucide-react";
import { Button, useOverlayState } from "@heroui/react";

import { types } from "@/types/core";

import { albumImageLink } from "@/util/images";
import { useApi } from "@/api/ApiContext";
import { PlayContextType } from "@/types/playcontext";
import { useSession } from "@/session/SessionContext";
import { TracksTable } from "@/components/TracksTable";
import { Event } from "@/types/events";
import { useMediaQuery } from "react-responsive";
import { mqMobile } from "@/config/config";
import DefaultLayout from "@/layouts/default";
import BannerNew from "@/components/BannerNew";
import { useBackdrop } from "@/hooks/BackgroundProvider";
import { About } from "../components/About";
import { AlbumEditModal } from "@/components/AlbumEditModal";
import { AlbumButton } from "@/components/AlbumButton";


type MoreBy = { artistId: string; name: string; albums: types.AlbumDetailed[] };

export default function AlbumPage({}: {}) {
  const [artists, setArtists] = useState<types.Artist[]>([]);
  const [album, setAlbum] = useState<types.AlbumDetailed>();
  const [tracks, setTracks] = useState<types.TrackDetailed[]>([]);
  const [moreBy, setMoreBy] = useState<MoreBy[]>([]);

  const params = useParams();
  const api = useApi();
  const session = useSession();
  const isMobile = useMediaQuery({ query: mqMobile })

  const editState = useOverlayState();


  useBackdrop(album && albumImageLink(album.id, 320));

  const loadAlbumData = () => {
    api.getAlbum(params.albumID as string).then(setAlbum);
    api.listTracksByAlbum(params.albumID as string).then(setTracks);
    api.listArtists().then(setArtists);
  };

  useEffect(() => {
    loadAlbumData();
  }, [api, params.albumID]);


  useEffect(() => {
    const ids = album?.artist_ids ?? [];
    if (ids.length === 0) {
      setMoreBy([]);
      return;
    }

    const perArtist = ids.length === 1 ? 4 : 2;

    Promise.all(ids.map((id) => api.listAlbumsByArtist(id))).then((lists) => {
      const seen = new Set([album!.id]);
      const sections = lists
        .map((albums, i) => ({
          artistId: ids[i],
          name: album!.artist_names![i],
          albums: albums.filter((a) => !seen.has(a.id) && seen.add(a.id)).slice(0, perArtist),
        }))
        .filter((s) => s.albums.length > 0);
      setMoreBy(sections);
    });
  }, [album?.id]);


  useEffect(() => {
    const syncF = api.eventSubscribe(Event.EntitiesUpdated, () => {
      loadAlbumData();
    });

    return () => {
      syncF();
    };
  }, []);


  const artistIds = album?.artist_ids;
  const artistNames = album?.artist_names;

  const year = album?.release_date ? new Date(album.release_date).getFullYear() : "Unknown year"
  const totalPlayCount = tracks.reduce((sum, t) => sum + t.play_count, 0);

  if (album) {
    return (
      <DefaultLayout>
        <AlbumEditModal album={album} artists={artists} artistIds={artistIds} state={editState} />
        <BannerNew
          bannerType="album"
          buttons={[
            {
              icon: ImageUp,
              onUpload: (files) => {
                api.uploadAlbumImage(
                  params.albumID as string,
                  files[0],
                );
              },
              isDisabled: !session.serverFullyAvailable,
              tooltip: "Upload Album Cover",
              type: "fileupload",
            },
            {
              icon: SquarePen,
              onClick: () => {
                editState.setOpen(true);
              },
              isDisabled: !session.serverFullyAvailable,
              tooltip: "Edit Album",
            },
          ]}
          imageUrl={albumImageLink(album.id, 320)}
          title={album.name}
        >
          <div className="mt-2">
            <div className="flex gap-4 items-center">
              <div className="flex gap-2">
                {artistIds?.map((artistId, index) => (
                  <Link
                    key={artistId}
                    style={{ textDecoration: "none", color: "inherit" }}
                    to={`/artists/${artistId}`}
                  >
                    <h2 className="text-lg font-black opacity-80 lg:text-xl xl:text-3xl hover:underline font-display underline-offset-3">
                      {artistNames?.[index]}
                      {index < artistIds.length - 1 && ", "}
                    </h2>
                  </Link>
                ))}
              </div>
              <p className="text-foreground/80">{`${year}, ${tracks.length} tracks`}</p>
            </div>
            <div className="flex gap-2 mt-6">
              <Button size="lg">
                <Play/>
                Play
              </Button>
              <Button size="lg" variant="tertiary" className="bg-background/0 border-1 border-foreground">
                <Shuffle/>
                Shuffle
              </Button>
            </div>
          </div>
        </BannerNew>

        <div className="flex flex-col-reverse gap-12 mt-8 xl:flex-row xl:gap-8">
          <div className="flex flex-col gap-10 justify-start w-full xl:w-80 xl:pt-15">
            <About
              text={album?.description}
              tableData={{
                "Released": album?.release_date ? new Date(album.release_date).toLocaleDateString("en-US", {
                    month: "short",
                    day: "numeric",
                    year: "numeric",
                    timeZone: "UTC",
                  }) : "Unknown",
                "Your Plays": totalPlayCount.toString(),
                "Visibility": album.public ? "Public" : "Private"
              }}
            />
            {moreBy.map((s) => (
              <section key={s.artistId} className="flex flex-col gap-3.5">
                <div className="flex items-baseline">
                  <h2 className="flex-1 font-bold font-display text-[22px] tracking-[-0.02em]">More by {s.name}</h2>
                  <Link to={`/artists/${s.artistId}`} className="text-sm font-semibold hover:underline text-accent-deep">All</Link>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  {s.albums.map((a) => <AlbumButton key={a.id} album={a} />)}
                </div>
              </section>
            ))}
          </div>
          <TracksTable
            columns={isMobile ? [
              "index",
              "title_artist",
              "album_options",
            ] : [
              "index",
              "title_artist",
              "media",
              "length",
              "plays",
              "rating",
              "like",
              "album_options",
            ]}
            parent_id={album?.id}
            parent_type={PlayContextType.Album}
            tracks={tracks}
          />
        </div>
      </DefaultLayout>
    )
  } else {
    <p>bej</p>
  }
}
