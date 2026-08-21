import { Button, Chip, Description, Key, Label, ListBox, Modal, Pagination, Select, Table, useOverlayState } from "@heroui/react";
import { CircleQuestionMark, CloudUpload, CirclePause, Search, UploadCloud, CirclePlay, CircleCheck, CircleAlert } from "lucide-react";
import { useEffect, useState } from "react";

import { useApi } from "@/api/ApiContext";
import PageTitleLayout from "@/layouts/page_title";
import { ActionButton } from "@/primitives/Buttons";
import { Input } from "@/primitives/Inputs";
import FileUploadButton from "@/components/FileUploadButton";
import { useSession } from "@/session/SessionContext";
import { Event } from "@/types/events";
import { responses } from "@/types/server";
import { types } from "@/types/core";
import { parseDate, dateHR } from "@/util/functions";
import { UserDetails } from "@/models/UserDetails";
import { UserBadge } from "@/components/UserBadge";
import { MusicbrainzSearch } from "@/components/MusicbrainzSearch";


export default function ImportPage() {
  const [page, setPage] = useState(1);
  const [importItems, setImportItems] = useState<responses.ListPaginationPayload<types.Import>>(new responses.ListPaginationPayload<types.Import>)
  const [users, setUsers] = useState<UserDetails[]>([]);
  const session = useSession();
  const modalState = useOverlayState();

  const api = useApi();

  const loadData = () => {
    api.listImportItems(page, 10).then(setImportItems)
  };

  useEffect(() => {
    loadData();
    api.listUsers(false).then(setUsers);
  }, [api]);

  useEffect(() => {
    loadData();
  }, [page]);

  useEffect(() => {
    const updatedF = api.eventSubscribe(Event.ImporterItemsUpdated, () => {
      loadData()
    });

    return () => {
      updatedF();
    };
  }, []);


  const user = (id: string): UserDetails | undefined => {
    return users.find((u) => u.id === id);
  };


  const chip = (state: string) => {
    let color: "default" | "success" | "danger" | "warning" = "default";
    let icon = <CircleQuestionMark size={18} />;

    switch (state) {
      case "not_started":
        color = "default";
        icon = <CirclePause size={18} />;
        break;
      case "running":
        color = "warning";
        icon = <CirclePlay size={18} />;
        break;
      case "finished":
        color = "success";
        icon = <CircleCheck size={18} />;
        break;
      case "error":
        color = "danger";
        icon = <CircleAlert size={18} />;
        break;
    }

    return (
      <Chip color={color} className="gap-1 px-1" size="sm" variant="secondary">
        {icon}
        {state.replace("_", " ")}
      </Chip>
    );
  };

  const count = importItems?.items?.length || 0;
  const currentPageItemStart = (page-1) * importItems.limit + 1
  const currentPageItemEnd = (page-1) * importItems.limit + count

  return (
    <PageTitleLayout
      title="Import Media"
      headerContent={
        <>
          <ActionButton
            icon={CloudUpload}
            isDisabled={!session.serverFullyAvailable}
            tooltip="Upload Files"
            onClick={() => {
              modalState.setOpen(true);
            }}
          />
        </>
      }
    >
      <UploadModal
        state={modalState}
      />
      <div className="flex flex-col w-full">
        <h2 className="pb-4 text-xl font-semibold">History</h2>
        <Table className="w-full">
          <Table.ScrollContainer>
            <Table.Content aria-label="Table with pagination" className="min-w-[600px]">
              <Table.Header>
                  <Table.Column>State</Table.Column>
                  <Table.Column>User</Table.Column>
                  <Table.Column>Updated At</Table.Column>
                  <Table.Column isRowHeader>Filename</Table.Column>
              </Table.Header>
              <Table.Body items={importItems.items}>
                {(item) => (
                  <Table.Row>

                <Table.Cell className="items-center">
                  {chip(item.state)}
                </Table.Cell>
                <Table.Cell>
                  <UserBadge size="sm" user={user(item.owner)} />
                </Table.Cell>
                    <Table.Cell className="truncate">{dateHR(parseDate(item.updated_at))}</Table.Cell>
                    <Table.Cell className="truncate width-full">{item.filename}</Table.Cell>
                  </Table.Row>
                )}
              </Table.Body>
            </Table.Content>
          </Table.ScrollContainer>
          <Table.Footer>
            <Pagination size="sm">
              <Pagination.Summary>
                {currentPageItemStart} to {currentPageItemEnd} of {importItems.total_items} results
              </Pagination.Summary>
              <Pagination.Content>
                <Pagination.Item>
                  <Pagination.Previous
                    isDisabled={page === 1}
                    onPress={() => setPage((p) => Math.max(1, p - 1))}
                  >
                    <Pagination.PreviousIcon />
                    Prev
                  </Pagination.Previous>
                </Pagination.Item>
                  {Array.from({ length: importItems.total_pages}, (_, i) => i + 1).map((p) => (
                    <Pagination.Item key={p}>
                      <Pagination.Link
                        isActive={p === page}
                        onPress={() => setPage(p)}
                      >
                        {p}
                      </Pagination.Link>
                    </Pagination.Item>
                  ))}
                <Pagination.Item>
                  <Pagination.Next
                    isDisabled={page === importItems.total_pages}
                    onPress={() => setPage((p) => Math.min(importItems.total_pages, p + 1))}
                  >
                    Next
                    <Pagination.NextIcon />
                  </Pagination.Next>
                </Pagination.Item>
              </Pagination.Content>
            </Pagination>
          </Table.Footer>
        </Table>
      </div>
    </PageTitleLayout>
  )
}


type State = ReturnType<typeof useOverlayState>;

function UploadModal({
  state,
}: {
  state: State;
}) {

  const [uploadType, setUploadType] = useState<Key | null>("album")
  const [mbID, setMbID] = useState<string | undefined>(undefined);
  const [showMbSearch, setShowMbSearch] = useState(false);

  const api = useApi();


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
                <CloudUpload/>
              </Modal.Icon>
              <Modal.Heading className="font-semibold">Upload Mediafile</Modal.Heading>
            </Modal.Header>
              {showMbSearch ?
                <MusicbrainzSearch setMbID={setMbID} setShowMbSearch={setShowMbSearch}/>
                :
               <>
                <Modal.Body className="px-1 pt-4">
                  <Select className="" value={uploadType} onChange={(value) => setUploadType(value)} >
                    <Label>Media Type</Label>
                    <Select.Trigger>
                      <Select.Value />
                      <Select.Indicator />
                    </Select.Trigger>
                    <Select.Popover>
                      <ListBox>
                        <ListBox.Item id="album" textValue="Album">
                          Album
                          <ListBox.ItemIndicator />
                        </ListBox.Item>
                        <ListBox.Item id="track" textValue="Track" isDisabled>
                          Track
                          <ListBox.ItemIndicator />
                        </ListBox.Item>
                      </ListBox>
                    </Select.Popover>
                    <Description>Select the type of upload you want to do</Description>
                  </Select>
                  <div className="flex gap-4 items-center py-4 w-full">
                    <Input
                      className="flex-grow"
                      color="secondary"
                      label="MusicBrainz ID"
                      name="musicbrainz_id"
                      onChange={(e) => {
                        setMbID(e);
                      }}
                      value={mbID}
                      description="Enter a MusicBrainz ID to make it easier for the importer. (optional)"
                    />
                    <ActionButton
                      icon={Search}
                      size="lg"
                      tooltip="Search Online"
                      variant="primary"
                      onClick={() => {setShowMbSearch(true)}}
                    />
                  </div>
              </Modal.Body>
              <Modal.Footer>
                <Button slot="close" variant="secondary">
                  Cancel
                </Button>

                <FileUploadButton
                  accept=".zip"
                  multiple={false}
                  variant="primary"
                  onUpload= {(files) => {
                    switch (uploadType) {
                        case "album":
                          api.importAlbum(
                            files[0],
                            mbID === undefined ? null : mbID,
                          );
                          state.setOpen(false)
                          return
                    }
                  }}
                >
                  <UploadCloud size={20} />
                  Upload
                </FileUploadButton>
              </Modal.Footer>
              </>
          }
          </Modal.Dialog>
        </Modal.Container>
      </Modal.Backdrop>
    </Modal>
  );
}
