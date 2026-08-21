import {
  Button,
  Modal,
  useOverlayState,
  Select,
  Key,
  ListBox,
} from "@heroui/react";
import { CircleCheck, ListMusic } from "lucide-react";
import { useEffect, useState } from "react";

import { types } from "@/types/core";

import { useApi } from "@/api/ApiContext";

type State = ReturnType<typeof useOverlayState>;

export default function PlaylistSelector({
  state,
  onSubmit,
}: {
  state: State;
  onSubmit: (id: string) => void;
}) {
  const [selectedID, setSelectedID] = useState<Key | null>(null);

  const [playlists, setPlaylists] = useState<types.Playlist[]>([]);

  const api = useApi();

  useEffect(() => {
    api.listPlaylists().then(setPlaylists);
  }, [api]);

  useEffect(() => {
    setSelectedID(null);
  }, []);

  const onClick = () => {
    if (selectedID === null) {
      return;
    }
    onSubmit(selectedID.toString());
  };

  return (
    <Modal>
      <Modal.Backdrop
        isOpen={state.isOpen}
        variant="opaque"
        onOpenChange={state.setOpen}
      >
        <Modal.Container size="lg">
          <Modal.Dialog className="">
            <Modal.CloseTrigger />
            <Modal.Header>
              <Modal.Icon className="bg-accent-soft text-accent-soft-foreground">
                <ListMusic />
              </Modal.Icon>
              <Modal.Heading className="font-semibold">
                Select Playlist
              </Modal.Heading>
            </Modal.Header>
            <Modal.Body className="px-1 pt-4">
              <Select
                placeholder="Select a playlist"
                value={selectedID}
                variant="secondary"
                onChange={(value) => {
                  setSelectedID(value);
                }}
              >
                <Select.Trigger>
                  <Select.Value />
                  <Select.Indicator />
                </Select.Trigger>
                <Select.Popover>
                  <ListBox>
                    {playlists && playlists.map((plist) => (
                      <ListBox.Item
                        key={plist.id}
                        id={plist.id}
                        textValue={plist.name}
                      >
                        {plist.name}
                        <ListBox.ItemIndicator />
                      </ListBox.Item>
                    ))}
                  </ListBox>
                </Select.Popover>
              </Select>
            </Modal.Body>
            <Modal.Footer>
              <Button slot="close" variant="secondary">
                Cancel
              </Button>
              <Button
                isDisabled={selectedID == undefined}
                slot="close"
                variant="primary"
                onPress={onClick}
              >
                <CircleCheck />
                Select Playlist
              </Button>
            </Modal.Footer>
          </Modal.Dialog>
        </Modal.Container>
      </Modal.Backdrop>
    </Modal>
  );
}
