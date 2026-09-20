import { useNavigate } from "react-router-dom";
import { Disc } from "lucide-react";

import { Text } from "@heroui/react";
import { types } from "@/types/core";

import { Image } from "./Image";

import { albumImageLink } from "@/util/images";
import { formatHRDate } from "../util/datetime";

interface AlbumRecentCardProps {
  album: types.AlbumDetailed;
}

export const AlbumRecentCard: React.FC<AlbumRecentCardProps> = ({ album }) => {
  const navigate = useNavigate();

  return (
    <div
      className="overflow-hidden cursor-pointer max-w-[160px]"
      onClick={() => {
        navigate("/albums/" + album.id);
      }}
    >
        <Image
          height={160}
          width={160}
          fallbackIcon={Disc}
          src={albumImageLink(album.id, 320)}
          radius="none"
          className="shadow-md border-1 border-border"
        />
      <Text type="body-sm" className="pt-1 -mb-1 font-bold truncate">{album.name}</Text>
      <Text type="body-sm" className="truncate">{album.artist_names?.join(", ")}</Text>
      <Text type="body-xs" className="truncate">{`Added ${formatHRDate(album.created_at)}`}</Text>
      </div>
  );
};
