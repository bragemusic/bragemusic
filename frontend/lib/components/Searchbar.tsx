import { Music2, Search } from "lucide-react";
import { useEffect, useState } from "react";
import { Button, Input, Modal, useOverlayState } from "@heroui/react";
import { Link } from "react-router-dom";

import { types } from "@/types/core";

import { Image } from "./Image";

import { albumImageLinkFromID, artistImageLinkFromID } from "@/util/images";
import { useApi } from "@/api/ApiContext";

const generateImageLink = (si: types.SearchItem): string => {
  if (si.link_type === "artist") {
    return artistImageLinkFromID(si.link_id.toString(), 320);
  } else if (si.link_type === "album") {
    return albumImageLinkFromID(si.link_id.toString(), 320);
  }

  return "";
};

const generateLink = (si: types.SearchItem): string => {
  if (si.link_type === "artist") {
    return "/artists/" + si.link_id;
  } else if (si.link_type === "album") {
    return "/albums/" + si.link_id;
  }

  return "";
};

interface SearchResultsProps {
  results: types.SearchItem[];
}

interface SearchbarProps {
  trigger?: React.ComponentType<{ onClick: () => void }>;
}

export const Searchbar: React.FC<SearchbarProps> = ({ trigger: Trigger }) => {
  const [results, setResults] = useState<types.SearchItem[]>([]);
  const [value, setValue] = useState("");
  const [debouncedValue, setDebouncedValue] = useState("");

  const state = useOverlayState();
  const api = useApi();

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, 400);

    return () => {
      clearTimeout(handler);
    };
  }, [value]);

  useEffect(() => {
    if (!debouncedValue) {
      setResults([]);

      return;
    }

    api.SearchFull(debouncedValue).then(setResults);
  }, [debouncedValue]);

  const onClose = () => {
    state.setOpen(false);
    setResults([]);
  };

  const SearchResults: React.FC<SearchResultsProps> = ({ results }) => (
    <div className="flex flex-col gap-2 justify-start mt-4">
      {results.map((si) => (
        <Link
          style={{ textDecoration: "none", color: "inherit" }}
          to={generateLink(si)}
          onClick={() => requestAnimationFrame(onClose)}
        >
          <div className="flex flex-row gap-3 p-2 rounded-lg bg-background hover:bg-accent hover:text-accent-foreground">
            <Image
              fallbackIcon={Music2}
              height={40}
              radius="sm"
              src={generateImageLink(si)}
              width={40}
            />
            <div className="flex flex-col">
              <p
                dangerouslySetInnerHTML={{ __html: si.html_name }}
                className="text-md"
              />
              <p className="capitalize text-small">{si.type}</p>
            </div>
          </div>
        </Link>
      ))}
    </div>
  );

  return (
    <div className="">
      {Trigger ? (
        <Trigger
          onClick={() => {
            state.setOpen(true);
          }}
        />
      ) : (
        <Button
          isIconOnly
          className="w-full text-md bg-background border-1 border-border text-foreground"
          size="sm"
          onPress={() => {
            state.setOpen(true);
          }}
        >
          <Search size={18} />
          <p className="hidden px-2 lg:block">Search</p>
        </Button>
      )}

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
                  <Search />
                </Modal.Icon>
                <Input
                  // classNames={{
                  //   base: "w-full border-b-1 border-default-700",
                  //   input: "text-xl my-5 ml-3",
                  //   inputWrapper: "bg-background py-8",
                  // }}
                  placeholder="Type to search..."
                  // startContent={<Search size={25} />}
                  // radius="none"
                  type="search"
                  variant="secondary"
                  onChange={(e) => setValue(e.target.value)}
                />
              </Modal.Header>
              <Modal.Body className="px-1 pt-4">
                <SearchResults results={results} />
              </Modal.Body>
            </Modal.Dialog>
          </Modal.Container>
        </Modal.Backdrop>
      </Modal>
    </div>
  );
};
