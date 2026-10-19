import EditModal from "./EditModal";
import { DateInput, Input, Textarea } from "@/primitives/Inputs";
import { types } from "@/types/core";
import { useOverlayState } from "@heroui/react";
import { useApi } from "../api/ApiContext";
import { ArtistSelect } from "@/primitives/Selects";
import { Switch } from "@/primitives/Switches";

type State = ReturnType<typeof useOverlayState>;

interface AlbumEditProps {
    album?: types.AlbumDetailed;
    artists: types.Artist[];
    artistIds?: string[];
    state: State;
}

export const AlbumEditModal: React.FC<AlbumEditProps> = ({
    album,
    artists,
    artistIds,
    state,
}) => {
    const api = useApi();

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

    return (
      <EditModal
        header={`Edit ${album?.name}`}
        state={state}
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
    );
};
