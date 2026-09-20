import { useEffect, useState } from "react";

import { types } from "@/types/core";

import { useApi } from "@/api/ApiContext";
import { Card, Text } from "@heroui/react";
import { SideScrollContainer } from "./SideScrollContainer";
import { SortBy, SortOrder } from "../types/sorting";
import { AlbumRecentCard } from "./AlbumRecentCard";

interface RecentAlbumsProps{}

export const RecentAlbums: React.FC<RecentAlbumsProps> = () => {
    const images = [
        "https://heroui-assets.nyc3.cdn.digitaloceanspaces.com/docs/robot1.jpeg",
        "https://heroui-assets.nyc3.cdn.digitaloceanspaces.com/docs/avocado.jpeg",
        "https://heroui-assets.nyc3.cdn.digitaloceanspaces.com/docs/oranges.jpeg",
    ];

    const [albums, setAlbums] = useState<types.AlbumDetailed[]>([]);

    const api = useApi();

    useEffect(() => {
        api.listAlbums(SortBy.Added, SortOrder.Desc, 10).then(setAlbums);
    }, [api]);

    return (
        <div className="flex flex-col w-full min-w-0">
            <Text type="h4">
            Recently Added Albums
            </Text>
            <SideScrollContainer>
      {albums &&
        albums.map(function (a) {
          return <AlbumRecentCard key={a.id} album={a} />;
        })}
            </SideScrollContainer>
        </div>
    );
};
