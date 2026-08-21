import { useState, useEffect } from "react";
import { Link, useParams } from "react-router-dom";
import { ImageUp, SquarePen } from "lucide-react";
import { useOverlayState } from "@heroui/react";

import { types } from "@/types/core";

import { albumImageLink } from "@/util/images";
import Banner from "@/components/Banner";
import EditModal from "@/components/EditModal";
import { DateInput, Input, Textarea } from "@/primitives/Inputs";
import { Switch } from "@/primitives/Switches";
import { useApi } from "@/api/ApiContext";
import { ArtistSelect } from "@/primitives/Selects";
import BodyLayout from "@/layouts/BodyLayout";
import { PlayContextType } from "@/types/playcontext";
import { useSession } from "@/session/SessionContext";
import { TracksTable } from "@/components/TracksTable";
import { Event } from "@/types/events";
import { useMediaQuery } from "react-responsive";
import { mqMobile } from "@/config/config";

export default function AlbumPage({
  scrollRef,
}: {
  scrollRef: React.RefObject<HTMLElement>;
}) {
  const [artists, setArtists] = useState<types.Artist[]>([]);
  const [album, setAlbum] = useState<types.AlbumDetailed>();
  const [tracks, setTracks] = useState<types.TrackDetailed[]>([]);

  const params = useParams();
  const api = useApi();
  const session = useSession();
  const isMobile = useMediaQuery({ query: mqMobile })

  const editState = useOverlayState();

  const loadAlbumData = () => {
    api.getAlbum(params.albumID as string).then(setAlbum);
    api.listTracksByAlbum(params.albumID as string).then(setTracks);
    api.listArtists().then(setArtists);
  };

  useEffect(() => {
    loadAlbumData();
  }, [api]);

  useEffect(() => {
    const syncF = api.eventSubscribe(Event.EntitiesUpdated, () => {
      loadAlbumData();
    });

    return () => {
      syncF();
    };
  }, []);

  const onAlbumEdit = (fd: FormData) => {
    const source = {
      musicbrainz_id: fd.get("musicbrainz_id") || undefined,
      name: String(fd.get("name") ?? ""),
      sort_name: String(fd.get("sort_name") ?? ""),
      artists: fd.getAll("artists"),
      release_date: fd.get("release_date") + "T00:00:00Z" || undefined,
      description: fd.get("description") || undefined,
      public: fd.get("public") != null,
    };

    const albumData = types.AlbumUpdate.createFrom(source);

    api.updateAlbum(album?.id as string, albumData);
  };

  const artistIds = album?.artist_ids;
  const artistNames = album?.artist_names;

  return (
    <>
      {/* --- Album Header --- */}
      {album && (
        <Banner
          bgImageUrl={albumImageLink(album.id, 2400)}
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
          coverImageUrl={albumImageLink(album.id, 320)}
          scrollRef={scrollRef}
          title={album.name}
        >
          <div className="flex gap-2">
            {artistIds?.map((artistId, index) => (
              <Link
                key={artistId}
                style={{ textDecoration: "none", color: "inherit" }}
                to={`/artists/${artistId}`}
              >
                <h2 className="text-lg opacity-80 lg:text-xl xl:text-3xl hover:underline underline-offset-3">
                  {artistNames?.[index]}
                  {index < artistIds.length - 1 && ", "}
                </h2>
              </Link>
            ))}
          </div>
          <p className="text-lg">
            {tracks.length} tracks •{" "}
            {album?.release_date
              ? new Date(album.release_date).getFullYear()
              : "Unknown year"}
          </p>
        </Banner>
      )}

      <EditModal
        header={`Edit ${album?.name}`}
        state={editState}
        onSubmit={onAlbumEdit}
      >
        <Input
          defaultValue={album?.musicbrainz_id}
          label="MusicBrainz ID"
          name="musicbrainz_id"
        />
        <Input
          required
          defaultValue={album?.name}
          description="The name of the album"
          label="Name"
          name="name"
        />
        <Input
          required
          defaultValue={album?.sort_name}
          label="Sort Name"
          name="sort_name"
        />

        <ArtistSelect
          artists={artists}
          name="artists"
          selectedArtistIds={artistIds}
        />

        <DateInput
          defaultValue={album?.release_date}
          label="Release Date"
          name="release_date"
        />
        <Textarea
          defaultValue={album?.description}
          label="Description"
          name="description"
        />
        <Switch
          className="mt-4"
          defaultSelected={album?.public}
          description="Allow other users on the server to listen to the album"
          label="Public"
          name="public"
        />
      </EditModal>

      <BodyLayout>
        {/* --- Tracklist --- */}
        <div className="sm:pt-10">
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
      </BodyLayout>
    </>
  );
}
