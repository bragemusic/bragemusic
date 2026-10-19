import { Link } from "react-router-dom";
import { Disc } from "lucide-react";
import { types } from "@/types/core";
import { Image } from "./Image";
import { albumImageLink } from "@/util/images";
import { cn } from "@heroui/react";

interface AlbumButtonProps {
  album: types.AlbumDetailed;
  size?: "sm" | "base"
}

export const AlbumButton: React.FC<AlbumButtonProps> = ({ album, size="base" }) => {
  return (
    <Link
      style={{ textDecoration: "none", color: "inherit" }}
      to={"/albums/" + album.id}
    >
      <div className={cn("pb-2 sm:pb-0", size == "base" ? "max-w-[180px]" : "max-w-[140px]")}>
        <Image
          className={cn("w-full aspect-square border-border border-1", size == "base" ? "shadow-[0_24px_36px_-26px_rgba(43,44,44,0.6)] max-w-[180px]" : "max-w-[140px]")}
          fallbackIcon={Disc}
          height={180}
          radius="sm"
          src={albumImageLink(album.id, 320)}
          width={180}
          customHeight={true}
        />
        <p className={cn("pt-2 text-sm text-foreground font-display sm:truncate sm:text-md", size == "base" ? "w-[180px]" : "w-[140px]")}>
          {album.name}
        </p>
        <p className="text-xs sm:text-sm text-foreground/80">
          {album?.release_date && new Date(album.release_date).getFullYear()}
        </p>
      </div>
    </Link>
  );
};
