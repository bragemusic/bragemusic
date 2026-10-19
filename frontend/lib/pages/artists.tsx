import { useEffect, useState } from "react";

import { types } from "@/types/core";

import { useApi } from "@/api/ApiContext";
import { ArtistCard } from "@/components/ArtistCard";
import CardLayout from "@/layouts/CardLayout";
import { Plus } from "lucide-react";
import { useSession } from "@/session/SessionContext";
import { Button, Key, SearchField, ToggleButton, ToggleButtonGroup, useOverlayState } from "@heroui/react";
import EditModal from "@/components/EditModal";
import { Input, NumberInput, Textarea } from "@/primitives/Inputs";
import { SortBy } from "@/types/sorting";

export default function ArtistsPage() {
  const [artists, setArtists] = useState<types.ArtistDetailed[]>([]);
  const [selectedSorting, setSelectedSorting] = useState(new Set<Key>([SortBy.Name]));

  const sortBy = ([...selectedSorting][0] ?? SortBy.Name) as SortBy;

  const api = useApi();
  const session = useSession();
  const addState = useOverlayState();

  useEffect(() => {
    api.listArtists().then(setArtists);
    console.log(sortBy)
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
      subtitle={`${artists.length} artists in your library`}
      headerContent={
        <div className="flex flex-col gap-4 items-center w-full xl:flex-row">
          <SearchField name="filter" className="w-full">
            <SearchField.Group  className="w-full h-10 rounded-full shadow-none border-border border-1 bg-surface">
              <SearchField.SearchIcon />
              <SearchField.Input className="w-full xl:w-[280px]" placeholder="Filter artists"/>
              <SearchField.ClearButton />
            </SearchField.Group>
          </SearchField>

          <div className="flex gap-4 justify-between items-center w-full">
            <ToggleButtonGroup selectedKeys={selectedSorting} onSelectionChange={(k) => {k.size > 0 && setSelectedSorting(k)}} selectionMode="single" size="lg">
              <ToggleButton id={SortBy.Name}>
                A-Z
              </ToggleButton>
              <ToggleButton id={SortBy.Added}>
                <ToggleButtonGroup.Separator />
                Recently added
              </ToggleButton>
              <ToggleButton id={SortBy.PlayCount}>
                <ToggleButtonGroup.Separator />
                Most played
              </ToggleButton>
            </ToggleButtonGroup>

            <Button size="lg" isDisabled={!session.serverFullyAvailable} onClick={() => {addState.setOpen(true)}} >
              <Plus/>
              Add artist
            </Button>
          </div>
        </div>
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
