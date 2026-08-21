import { useEffect, useState } from "react";
import { Plus, Sparkles } from "lucide-react";
import { useOverlayState } from "@heroui/react";

import { types } from "@/types/core";

import { useApi } from "@/api/ApiContext";
import CardLayout from "@/layouts/CardLayout";
import { ActionButton } from "@/primitives/Buttons";
import EditModal from "@/components/EditModal";
import { Input, Textarea } from "@/primitives/Inputs";
import { Switch } from "@/primitives/Switches";
import { PlaylistCard } from "@/components/PlaylistCard";
import { useSession } from "@/session/SessionContext";
import { SmartPlaylistModal } from "@/components/SmartPlaylistModal";

export default function PlaylistsPage() {
  const [playlists, setPlaylists] = useState<types.Playlist[]>([]);
  const [artists, setArtists] = useState<types.Artist[]>([]);

  const api = useApi();
  const addState = useOverlayState();
  const addSmartState = useOverlayState();
  const session = useSession();

  useEffect(() => {
    api.listPlaylists().then(setPlaylists);
    api.listArtists().then(setArtists);
  }, [api]);

  const onAddPlaylist = (fd: FormData) => {
    const source = {
      name: String(fd.get("name") ?? ""),
      description: fd.get("description") || undefined,
      public: fd.get("public") != null,
    };

    const plistData = types.Playlist.createFrom(source);

    api.addPlaylist(plistData);
  };

  return (
    <CardLayout
      headerContent={
        <>
          <ActionButton
            icon={Plus}
            isDisabled={!session.serverFullyAvailable}
            tooltip="Add new playlist"
            onClick={() => {
              addState.setOpen(true);
            }}
          />
          <ActionButton
            icon={Sparkles}
            isDisabled={!session.serverFullyAvailable}
            tooltip="Add new smart playlist"
            onClick={() => {
              addSmartState.setOpen(true);
            }}
          />
        </>
      }
      title="Playlists"
    >
      <EditModal
        header="Add Playlist"
        state={addState}
        onSubmit={onAddPlaylist}
      >
        <Input required label="Name" name="name" />
        <Textarea label="Description" name="description" />
        <Switch
          className="mt-4"
          defaultSelected={false}
          label="Public"
          name="public"
        />
      </EditModal>
      <SmartPlaylistModal state={addSmartState} artists={artists} />
      {playlists &&
        playlists.map(function (p) {
          return <PlaylistCard key={p.id} playlist={p} />;
        })}
    </CardLayout>
  );
}
