import { useOverlayState } from "@heroui/react";

import { types } from "@/types/core";
import { ArtistRole } from "@/types/artist";

import EditModal from "./EditModal";

import { useApi } from "@/api/ApiContext";
import { Input, NumberInput, Textarea } from "@/primitives/Inputs";
import { ArtistSelect } from "@/primitives/Selects";
import { useEffect, useState } from "react";

type State = ReturnType<typeof useOverlayState>;

export default function EditTrackModal({
  state,
  selectedEditTrack,
  parent_id,
}: {
  state: State;
  selectedEditTrack: types.TrackDetailed | null;
  parent_id: string | undefined;
}) {
  const api = useApi();

  const [artists, setArtists] = useState<types.Artist[]>([]);
  const [availableArtists, setAvailableArtists] = useState<types.Artist[]>([]);

  const loadData= () => {
    api.listArtists().then(setArtists);
  };

  useEffect(() => {
    loadData();
  }, [api]);

  useEffect(() => {
    setAvailableArtists(artists.filter((a => {
      let isPrimary = false

      selectedEditTrack?.artists?.forEach((at) => {
        if (at.id !== a.id) {
          return
        }

        if (at.role === ArtistRole.Primary) {
          isPrimary = true
        }
      })

      console.log(isPrimary)
      return !isPrimary
    })))

  }, [artists, selectedEditTrack]);

  const onTrackEdit = (fd: FormData) => {
    if (selectedEditTrack === null) {
      return;
    }
    const source = {
      musicbrainz_id: fd.get("musicbrainz_id") || undefined,
      title: String(fd.get("title") ?? ""),
      comment: fd.get("comment") || undefined,
      disc_number: Number(fd.get("disc_number")),
      track_number: Number(fd.get("track_number")),
      album_id: parent_id,
      artists: fd.getAll("track_artists"),
    };
console.log(source)
    const trackData = types.TrackUpdate.createFrom(source);

    api.updateTrack(selectedEditTrack.id, trackData);
  };

  return (
    <EditModal header={`Edit Track`} state={state} onSubmit={onTrackEdit}>
      {selectedEditTrack != null && (
        <>
          <Input
            required
            defaultValue={selectedEditTrack.title}
            label="Name"
            name="title"
          />
          <Input
            defaultValue={selectedEditTrack.musicbrainz_id}
            label="MusicBrainz ID"
            name="musicbrainz_id"
          />
          <ArtistSelect
            artists={availableArtists}
            name="track_artists"
            selectedArtistIds={selectedEditTrack.artists?.map((a) => (a.id))}
          />
          <NumberInput
            required
            defaultValue={selectedEditTrack.disc_number}
            label="Disc"
            name="disc_number"
          />
          <NumberInput
            required
            defaultValue={selectedEditTrack.track_number}
            label="Track"
            name="track_number"
          />
          <Textarea
            defaultValue={selectedEditTrack.comment}
            label="Comment"
            name="comment"
          />
        </>
      )}
    </EditModal>
  );
}
