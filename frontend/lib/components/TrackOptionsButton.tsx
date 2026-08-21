import { Dropdown, Button, Header } from "@heroui/react";
import { Ellipsis, SquarePen, Trash2, ListPlus, Plus } from "lucide-react";

import { DropdonwMenuItem } from "@/primitives/DropdownItems/DropdownMenuItem";
import { useSession } from "@/session/SessionContext";

interface TrackOptionsButtonProps {
  trackIdx: number;
  parentType: "album" | "playlist" | "liked_tracks";
  onCallback: (trackIdx: number, callbackType: string) => void;
}

export const TrackOptionsButton: React.FC<TrackOptionsButtonProps> = ({
  trackIdx,
  parentType,
  onCallback,
}) => {
  const session = useSession();

  const callback = (callbackType: React.Key) => {
    onCallback(trackIdx, callbackType.toString());
  };

  const disabledAlbumKeys = (): string[] => {
    let keys = ["delete"];

    if (!session.serverFullyAvailable) {
      keys = keys.concat(["edit", "add_to_playlist"]);
    }

    return keys;
  };

  const disabledPlaylistKeys = (): string[] => {
    let keys: string[] = [];

    if (!session.serverFullyAvailable) {
      keys = keys.concat(["delete"]);
    }

    return keys;
  };

  const disabledLikedTracksKeys = (): string[] => {
    let keys: string[] = [];

    return keys;
  };

  let content;

  if (parentType === "album") {
    content = (
      <Dropdown.Menu
        aria-label="Dropdown menu with description"
        disabledKeys={disabledAlbumKeys()}
        onAction={callback}
      >
        <Dropdown.Section>
          <Header>Actions</Header>
          <DropdonwMenuItem
            description="Edit the selected track's metadata"
            icon={SquarePen}
            id="edit"
            label="Edit Track"
          />
          <DropdonwMenuItem
            description="Add track to the play queue"
            icon={ListPlus}
            id="add_to_queue"
            label="Add to Queue"
          />
          <DropdonwMenuItem
            description="Add track to a selected playlist"
            icon={Plus}
            id="add_to_playlist"
            label="Add to Playlist"
          />
        </Dropdown.Section>
        <Dropdown.Section>
          <Header>Danger zone</Header>
          <DropdonwMenuItem
            description="Permanently delete the track"
            icon={Trash2}
            id="delete"
            label="Delete Track"
            variant="danger"
          />
        </Dropdown.Section>
      </Dropdown.Menu>
    );
  } else if (parentType === "playlist") {
    content = (
      <Dropdown.Menu
        aria-label="Dropdown menu with description"
        disabledKeys={disabledPlaylistKeys()}
        onAction={callback}
      >
        <Dropdown.Section>
          <Header>Actions</Header>
          <DropdonwMenuItem
            description="Add track to the play queue"
            icon={ListPlus}
            id="add_to_queue"
            label="Add to Queue"
          />
        </Dropdown.Section>
        <Dropdown.Section>
          <Header>Danger xone</Header>
          <DropdonwMenuItem
            description="Remove the track from the playlist"
            icon={Trash2}
            id="delete"
            label="Remove Track"
            variant="danger"
          />
        </Dropdown.Section>
      </Dropdown.Menu>
    );
  } else if (parentType === "liked_tracks") {
    content = (
      <Dropdown.Menu
        aria-label="Dropdown menu with description"
        disabledKeys={disabledLikedTracksKeys()}
        onAction={callback}
      >
        <Dropdown.Section>
          <Header>Actions</Header>
          <DropdonwMenuItem
            description="Add track to the play queue"
            icon={ListPlus}
            id="add_to_queue"
            label="Add to Queue"
          />
        </Dropdown.Section>
      </Dropdown.Menu>
    );
  }

  return (
    <div className="flex flex-wrap gap-4">
      <Dropdown>
        <Button isIconOnly size="sm" variant="ghost">
          <Ellipsis />
        </Button>
        <Dropdown.Popover>{content}</Dropdown.Popover>
      </Dropdown>
    </div>
  );
};
