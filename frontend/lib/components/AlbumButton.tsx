import { Link } from "react-router-dom";
import { Disc } from "lucide-react";
import { types } from "@/types/core";
import { Image } from "./Image";
import { albumImageLink } from "@/util/images";

interface AlbumButtonProps {
  album: types.AlbumDetailed;
}

export const AlbumButton: React.FC<AlbumButtonProps> = ({ album }) => {
  return (
    <Link
      style={{ textDecoration: "none", color: "inherit" }}
      to={"/albums/" + album.id}
    >
      <div className="pb-2 sm:pb-0 max-w-[180px]">
        <Image
          className="w-full max-w-[180px] aspect-square"
          fallbackIcon={Disc}
          height={180}
          radius="none"
          src={albumImageLink(album.id, 320)}
          width={180}
          customHeight={true}
        />
        <p className="pt-2 text-sm text-foreground w-[180px] sm:truncate sm:w-[180px] sm:text-md">
          {album.name}
        </p>
        <p className="text-xs sm:text-sm text-foreground/80">
          {album?.release_date && new Date(album.release_date).getFullYear()}
        </p>
      </div>
    </Link>
  );
};
