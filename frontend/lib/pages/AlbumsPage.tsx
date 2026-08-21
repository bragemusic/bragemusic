import { useEffect, useState } from "react";

import { types } from "@/types/core";

import { useApi } from "@/api/ApiContext";
import CardLayout from "@/layouts/CardLayout";
import { AlbumCard } from "@/components/AlbumCard";

export default function AlbumsPage() {
  const [albums, setAlbums] = useState<types.AlbumDetailed[]>([]);

  const api = useApi();

  useEffect(() => {
    api.listAlbums().then(setAlbums);
  }, [api]);

  return (
    <CardLayout title="Albums">
      {albums &&
        albums.map(function (a) {
          return <AlbumCard key={a.id} album={a} />;
        })}
    </CardLayout>
  );
}
