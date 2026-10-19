import { useEffect, useState } from "react";

import { types } from "@/types/core";

import { useApi } from "@/api/ApiContext";
import CardLayout from "@/layouts/CardLayout";
import { AlbumCard } from "@/components/AlbumCard";
import { SortBy, SortOrder } from "@/types/sorting";
import { Button, Key, SearchField, ToggleButton, ToggleButtonGroup } from "@heroui/react";
import { Plus } from "lucide-react";

export default function AlbumsPage() {
  const [albums, setAlbums] = useState<types.AlbumDetailed[]>([]);
  const [selectedSorting, setSelectedSorting] = useState(new Set<Key>([SortBy.Name]));

  const sortBy = ([...selectedSorting][0] ?? SortBy.Name) as SortBy;

  const api = useApi();

  useEffect(() => {
    api.listAlbums(sortBy, SortOrder.Asc).then(setAlbums);
  }, [api, sortBy]);

  return (
    <CardLayout
      title="Albums"
      subtitle={`${albums.length} albums in your library`}
      headerContent={
        <div className="flex flex-col gap-4 items-center w-full xl:flex-row">
          <SearchField name="filter" className="w-full">
            <SearchField.Group  className="w-full h-10 rounded-full shadow-none border-border border-1 bg-surface">
              <SearchField.SearchIcon />
              <SearchField.Input className="w-full xl:w-[280px]" placeholder="Filter albums"/>
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

            <Button size="lg" isDisabled >
              <Plus/>
              Add album
            </Button>
          </div>
        </div>
      }
    >
      {albums &&
        albums.map(function (a) {
          return <AlbumCard key={a.id} album={a} />;
        })}
    </CardLayout>
  );
}
