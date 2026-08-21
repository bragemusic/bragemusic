import { User } from "lucide-react";
import { useNavigate } from "react-router-dom";

import { types } from "@/types/core";

import { Card } from "./Card";

import { artistImageLink } from "@/util/images";

interface ArtistCardProps {
  artist: types.ArtistDetailed;
}

export const ArtistCard: React.FC<ArtistCardProps> = ({ artist }) => {
  const navigate = useNavigate();

  return (
    <Card
      description={`${artist.album_count} album${artist.album_count > 1 ? "s " : ""} • ${artist.track_count} track${artist.track_count > 1 ? "s " : ""}`}
      fallbackIcon={User}
      radius="full"
      imgSrc={artistImageLink(artist, 320)}
      title={artist.name}
      onClick={() => {
        navigate("/artists/" + artist.id);
      }}
    />
  );
};
