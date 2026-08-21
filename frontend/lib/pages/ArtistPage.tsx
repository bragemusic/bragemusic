import { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import { ChevronDown, ChevronUp, ImageUp, UserPen } from "lucide-react";
import { useOverlayState } from "@heroui/react";

import { types } from "@/types/core";

import { AlbumButton } from "@/components/AlbumButton";
import Banner from "@/components/Banner";
import { artistImageLink } from "@/util/images";
import EditModal from "@/components/EditModal";
import { Input, Textarea, NumberInput } from "@/primitives/Inputs";
import BodyLayout from "@/layouts/BodyLayout";
import { useSession } from "@/session/SessionContext";
import { TracksTable } from "@/components/TracksTable";
import { PlayContextType } from "@/types/playcontext";
import { ActionButton } from "@/primitives/Buttons";
import { useApi } from "@/api/ApiContext";
import { Event } from "@/types/events";
import { useMediaQuery } from "react-responsive";
import { mqMobile } from "@/config/config";

export default function ArtistPage({
  scrollRef,
}: {
  scrollRef: React.RefObject<HTMLElement>;
}) {
  const [artist, setArtist] = useState<types.Artist>();
  const [albums, setAlbums] = useState<types.AlbumDetailed[]>([]);
  const [featuredAlbums, setFeaturedAlbums] = useState<types.AlbumDetailed[]>([]);
  const [topTracks, setTopTracks] = useState<types.TrackDetailed[]>([]);
  const [topTracksOpen, setTopTracksOpen] = useState(false);

  const api = useApi();

  const params = useParams();
  const session = useSession();
  const isMobile = useMediaQuery({ query: mqMobile })


  const editState = useOverlayState();

  const loadData = () => {
    api.getArtist(params.artistID as string).then(setArtist);
    api.listAlbumsByArtist(params.artistID as string).then(setAlbums);
    api.listFeaturedAlbumsByArtist(params.artistID as string).then(setFeaturedAlbums);
    api.getArtistTopTracks(params.artistID as string).then(setTopTracks);
  };

  useEffect(() => {
    loadData();
  }, [api]);

  useEffect(() => {
    const evUpdated = api.eventSubscribe(Event.EntitiesUpdated, () => {
      loadData();
    });

    return () => {
      evUpdated();
    };
  }, []);

  const onArtistEdit = (fd: FormData) => {
    if (!artist) {
      return
    }

    const source = {
      musicbrainz_id: fd.get("musicbrainz_id") || undefined,
      name: String(fd.get("name") ?? ""),
      sort_name: String(fd.get("sort_name") ?? ""),
      country: fd.get("country") || undefined,

      year_started: fd.get("year_started")
        ? Number(fd.get("year_started"))
        : undefined,

      year_ended: fd.get("year_ended")
        ? Number(fd.get("year_ended"))
        : undefined,

      description: fd.get("description") || undefined,
    };

    const artistData = types.Artist.createFrom(source);

    api.updateArtist(artist.id.toString(), artistData);
  };

  return (
    <>
      {artist && (
        <Banner
          bgImageUrl={artistImageLink(artist, 2400)}
          buttons={[
            {
              icon: ImageUp,
              onUpload: (files) => {
                api.uploadArtistImage(
                  params.artistID as string,
                  files[0],
                );
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
          ]}
          scrollRef={scrollRef}
          title={artist?.name}
        >
          <p className="pt-6">{artist?.description}</p>
        </Banner>
      )}

      <EditModal
        header={`Edit ${artist?.name}`}
        state={editState}
        onSubmit={onArtistEdit}
      >
        <Input
          defaultValue={artist?.musicbrainz_id}
          description="Used for collecting information in the MusicBrainz database"
          label="MusicBrainz ID"
          name="musicbrainz_id"
        />
        <Input
          required
          defaultValue={artist?.name}
          description="The name of the artist"
          label="Name"
          name="name"
        />
        <Input
          required
          defaultValue={artist?.sort_name}
          description="Name used for sorting"
          label="Sort Name"
          name="sort_name"
        />
        <Input
          defaultValue={artist?.country}
          description="The country from where the artist originates"
          label="Country"
          name="country"
        />
        <NumberInput
          defaultValue={artist?.year_started}
          description="The year the artist first started performing"
          label="Year Started"
          name="year_started"
        />
        <NumberInput
          defaultValue={artist?.year_ended}
          description="The year when the artist quit performing"
          label="Year Ended"
          name="year_ended"
        />
        <Textarea
          defaultValue={artist?.description}
          description="Artist description and biography"
          label="Description"
          name="description"
        />
      </EditModal>

      <BodyLayout>
        <div>
          <h2 className="mb-2 text-xl font-semibold">Top Tracks</h2>
          {artist && (
            <div className="flex flex-col gap-2 mb-10">
              <TracksTable
                alteringRowColor={false}
                columns={isMobile ? ["album_cover", "title", "length", "plays"] : ["album_cover", "title", "length", "plays"]}
                parent_id={artist.id.toString()}
                parent_type={PlayContextType.TopTracks}
                tracks={topTracks}
                visibleCount={topTracksOpen ? 10 : 5}
              />
              <ActionButton
                icon={topTracksOpen ? ChevronUp : ChevronDown}
                size="sm"
                tooltip={topTracksOpen ? "Show less" : "Show more"}
                onClick={() => {
                  setTopTracksOpen(!topTracksOpen);
                }}
              />
            </div>
          )}
        </div>

        {albums.length > 0 &&
          <>
          <h2 className="pt-8 text-xl font-semibold">Albums</h2>
          <div className="grid grid-cols-2 gap-2 pt-4 sm:flex sm:flex-wrap sm:gap-8 sm:pl-4">
            {albums.map(function (a) {
              return <AlbumButton key={a.id} album={a} />;
            })}
          </div>
          </>
        }

        {featuredAlbums.length > 0 &&
          <>
          <h2 className="pt-8 text-xl font-semibold">Appears on</h2>
          <div className="grid grid-cols-2 gap-2 pt-4 sm:flex sm:flex-wrap sm:gap-8 sm:pl-4">
            {featuredAlbums.map(function (a) {
              return <AlbumButton key={a.id} album={a} />;
            })}
          </div>
          </>
        }
      </BodyLayout>
    </>
  );
}
