import { useEffect, useState } from "react";

import { types } from "@/types/core";

import { useApi } from "@/api/ApiContext";
import { ArtistCard } from "@/components/ArtistCard";
import CardLayout from "@/layouts/CardLayout";
import { ActionButton } from "@/primitives/Buttons";
import { Plus } from "lucide-react";
import { useSession } from "@/session/SessionContext";
import { useOverlayState } from "@heroui/react";
import EditModal from "@/components/EditModal";
import { Input, NumberInput, Textarea } from "@/primitives/Inputs";

export default function ArtistsPage() {
  const [artists, setArtists] = useState<types.ArtistDetailed[]>([]);

  const api = useApi();
  const session = useSession();
  const addState = useOverlayState();

  useEffect(() => {
    api.listArtists().then(setArtists);
  }, [api]);

  const onAddArtist = (fd: FormData) => {
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

    const artistData = types.ArtistBase.createFrom(source);

    api.createArtist(artistData);
  };

  return (
    <CardLayout
      title="Artists"
      headerContent={
        <>
          <ActionButton
            icon={Plus}
            isDisabled={!session.serverFullyAvailable}
            tooltip="Add new artist"
            onClick={() => {
              addState.setOpen(true);
            }}
          />
        </>
      }
    >
      <EditModal
        header="Add Artist"
        state={addState}
        onSubmit={onAddArtist}
      >
        <Input
          description="Used for collecting information in the MusicBrainz database"
          label="MusicBrainz ID"
          name="musicbrainz_id"
        />
        <Input
          required
          description="The name of the artist"
          label="Name"
          name="name"
        />
        <Input
          required
          description="Name used for sorting"
          label="Sort Name"
          name="sort_name"
        />
        <Input
          description="The country from where the artist originates"
          label="Country"
          name="country"
        />
        <NumberInput
          description="The year the artist first started performing"
          label="Year Started"
          name="year_started"
        />
        <NumberInput
          description="The year when the artist quit performing"
          label="Year Ended"
          name="year_ended"
        />
        <Textarea
          description="Artist description and biography"
          label="Description"
          name="description"
        />
      </EditModal>
      {artists &&
        artists.map((a) => (
          <ArtistCard key={a.id.toString()} artist={a} />
        ))
      }
    </CardLayout>
  );
}

// type ArtistBase struct {
// 	MusicBrainzID *string `db:"musicbrainz_id" json:"musicbrainz_id"`
// 	Name          string  `db:"name" json:"name" required:"true"`
// 	SortName      string  `db:"sort_name" json:"sort_name" required:"true"`
// 	Country       *string `db:"country" json:"country,omitempty"`
// 	YearStarted   *int    `db:"year_started" json:"year_started,omitempty"`
// 	YearEnded     *int    `db:"year_ended" json:"year_ended,omitempty"`
// 	Description   *string `db:"description" json:"description,omitempty"`
// }
