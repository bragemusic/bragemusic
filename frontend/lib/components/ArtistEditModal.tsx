import EditModal from "./EditModal"
import { Input, Textarea, NumberInput } from "@/primitives/Inputs";
import { types } from "@/types/core";
import { useOverlayState } from "@heroui/react";
import { useApi } from "../api/ApiContext";

type State = ReturnType<typeof useOverlayState>;

interface ArtistEditProps {
    artist?: types.Artist;
    state: State;
}

export const ArtistEditModal: React.FC<ArtistEditProps> = ({
    artist,
    state,
}) => {

  const api = useApi();

  const onArtistEdit = (fd: FormData) => {
    if (!artist) {
      return;
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
      <EditModal
        header={`Edit ${artist?.name}`}
        state={state}
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
    )
}
