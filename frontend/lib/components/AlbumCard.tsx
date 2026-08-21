import { useNavigate } from "react-router-dom";
import { Disc } from "lucide-react";

import { types } from "@/types/core";

import { Card } from "./Card";

import { albumImageLink } from "@/util/images";

interface AlbumCardProps {
  album: types.AlbumDetailed;
}

export const AlbumCard: React.FC<AlbumCardProps> = ({ album }) => {
  const navigate = useNavigate();

  return (
    <Card
      description={`${album.artist_names?.join(", ")}, ${album?.release_date && new Date(album.release_date).getFullYear()}`}
      fallbackIcon={Disc}
      imgSrc={albumImageLink(album.id, 320)}
      radius="none"
      title={album.name}
      onClick={() => {
        navigate("/albums/" + album.id);
      }}
    />
  );
};
