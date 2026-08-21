import { cn } from "@heroui/react";
import { Heart } from "lucide-react";

import { useApi } from "@/api/ApiContext";

interface LikeHeartProps {
  liked: boolean;
  track_id: string;
}

export const LikeHeart: React.FC<LikeHeartProps> = ({ liked, track_id }) => {
  const api = useApi();

  const Like = () => {
    if (liked) {
      api.unlikeTrack(track_id);
    } else {
      api.likeTrack(track_id);
    }
  };

  return (
    <div onClick={Like}>
      <Heart
        className={cn(
          liked ? "text-danger" : "text-foreground/20",
          "fill-current hover:text-danger",
        )}
        size={20}
        strokeWidth={0}
      />
    </div>
  );
};
