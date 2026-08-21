import { useState, useEffect } from "react";
import { useParams } from "react-router-dom";

import { FilePenLine, ImageUp, Trash2 } from "lucide-react";
import { types } from "@/types/core";
import { playlistImageLinkFromID } from "@/util/images";
import Banner from "@/components/Banner";
import { useApi } from "@/api/ApiContext";
import BodyLayout from "@/layouts/BodyLayout";
import { TracksTable } from "@/components/TracksTable";
import { PlayContextType } from "@/types/playcontext";
import EditModal from "@/components/EditModal";
import { Input, Textarea } from "@/primitives/Inputs";
import { Switch } from "@/primitives/Switches";

import { useSession } from "@/session/SessionContext";

import { useOverlayState } from "@heroui/react";
import { Event } from "@/types/events";
import { useMediaQuery } from "react-responsive";
import { mqMobile } from "@/config/config";

export default function PlaylistPage({
  scrollRef,
}: {
  scrollRef: React.RefObject<HTMLElement>;
}) {
  const [playlist, setPlaylist] = useState<types.Playlist | undefined>(
    undefined,
  );
  const [tracks, setTracks] = useState<types.TrackDetailed[]>([]);

  const params = useParams();
  const api = useApi();
  const session = useSession();

  const editState = useOverlayState();

  const isMobile = useMediaQuery({ query: mqMobile })

  const loadData = () => {
    api.getPlaylist(params.playlistID as string).then(setPlaylist);
    api.listPlaylistTracks(params.playlistID as string).then(setTracks);
  };

  useEffect(() => {
    loadData();
  }, [api]);

  useEffect(() => {
    const syncF = api.eventSubscribe(Event.EntitiesUpdated, () => {
      loadData();
    });

    return () => {
      syncF();
    };
  }, []);

  const onPlaylistEdit = (fd: FormData) => {
    if (playlist === undefined) {
      return;
    }

    const source = {
      name: String(fd.get("name") ?? ""),
      description: fd.get("description") || undefined,
      public: fd.get("public") != null,
    };

    const plistData = types.Playlist.createFrom(source);

    api.updatePlaylist(playlist?.id.toString(), plistData);
  };

  const onPlaylistDelete = () => {
    if (playlist === undefined) {
      return;
    }
    api.deletePlaylist(playlist?.id);
  };

  return (
    <>
        {playlist && (
          <Banner
            bgImageUrl={playlistImageLinkFromID(playlist.id.toString(), 2400)}
            coverImageUrl={playlistImageLinkFromID(playlist.id.toString(), 320)}
            scrollRef={scrollRef}
            title={playlist.name}
            buttons={[
              {
                icon: ImageUp,
                onUpload: (files) => {
                  api.uploadPlaylistImage(
                    params.playlistID as string,
                    files[0],
                  );
                },
                isDisabled: !session.serverFullyAvailable,
                tooltip: "Upload Playlist Image",
                type: "fileupload",
              },
              {
                icon: FilePenLine,
                onClick: () => {
                  editState.setOpen(true);
                },
                isDisabled: !session.serverFullyAvailable,
                tooltip: "Edit Playlist",
              },
              {
                icon: Trash2,
                onClick: onPlaylistDelete,
                isDisabled: !session.serverFullyAvailable,
                tooltip: "Remove Playlist",
                variant: "danger",
                confirm: true,
              },
            ]}
          >
            <p className="pt-6">{playlist.description}</p>
          </Banner>
        )}

        {
          <EditModal
            state={editState}
            header={`Edit ${playlist?.name}`}
            onSubmit={onPlaylistEdit}
          >
            <Input
              name="name"
              label="Name"
              defaultValue={playlist?.name}
              required
            />
            <Textarea
              name="description"
              label="Description"
              defaultValue={playlist?.description}
            />
            <Switch
              name="public"
              label="Public"
              defaultSelected={playlist?.public}
              className="mt-4"
            />
          </EditModal>
        }

        <BodyLayout>
          <div className="pt-10">
            {playlist && (
              <TracksTable
                columns={isMobile ? [
                  "album_cover",
                  "title_artist",
                  "plist_options",
                ] : [
                  "play_symbol",
                  "title",
                  "media",
                  "artist",
                  "album",
                  "length",
                  "plays",
                  "rating",
                  "like",
                  "plist_options",
                ]}
                parent_id={playlist.id.toString()}
                parent_type={PlayContextType.Playlist}
                tracks={tracks}
              />
            )}
          </div>
        </BodyLayout>
      </>
  );
}
