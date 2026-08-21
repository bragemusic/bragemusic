import { Star, StarHalf } from "lucide-react";
import { useState } from "react";

import { useApi } from "@/api/ApiContext";

interface StarRatingProps {
  user_rating?: number;
  rating?: number;
  track_id: string;
}

export const StarRating: React.FC<StarRatingProps> = ({
  user_rating,
  rating,
  track_id,
}) => {
  const [hoverIdx, setHoverIdx] = useState<number | null>(null);

  const api = useApi();

  const Rate = (value: number) => {
    api.rateTrack(track_id, value);
  };

  const fillColor = (fallbackColor: string, idx: number): string => {
    if (hoverIdx === null) {
      return fallbackColor;
    }

    if (hoverIdx < idx) {
      return fallbackColor;
    }

    return "fill-[#FFD447]";
  };

  return (
    <div className="relative">
      <div className="flex gap-0">
        {Array.from({ length: 5 }, (_, idx) => (
          <Star
            key={idx}
            className={
              hoverIdx != null && hoverIdx >= idx
                ? "fill-[#FFD447]"
                : "fill-foreground/20"
            }
            size={20}
            strokeWidth={0}
            onClick={() => Rate(idx + 1)}
            onMouseEnter={() => setHoverIdx(idx)}
            onMouseLeave={() => setHoverIdx(null)}
          />
        ))}
      </div>
      {rating != undefined && user_rating == undefined && (
        <div className="flex absolute top-0 gap-0 z-8">
          {Array.from({ length: Math.floor(rating ?? 0) }, (_, idx) => (
            <Star
              key={idx}
              className={fillColor("fill-foreground", idx)}
              size={20}
              strokeWidth={0}
              onClick={() => Rate(idx + 1)}
              onMouseEnter={() => setHoverIdx(idx)}
              onMouseLeave={() => setHoverIdx(null)}
            />
          ))}
          {(rating ?? 0) % 1 > 0 && hoverIdx == null && (
            <StarHalf className="fill-foreground" size={20} strokeWidth={0} />
          )}
        </div>
      )}
      <div className="flex absolute top-0 z-10 gap-0">
        {Array.from({ length: Math.floor(user_rating ?? 0) }, (_, idx) => (
          <Star
            key={idx}
            className={
              hoverIdx != null && hoverIdx < idx
                ? "fill-foreground/20"
                : "fill-[#FFD447]"
            }
            size={20}
            strokeWidth={0}
            onClick={() => Rate(idx + 1)}
            onMouseEnter={() => setHoverIdx(idx)}
            onMouseLeave={() => setHoverIdx(null)}
          />
        ))}
        {(user_rating ?? 0) % 1 > 0 && hoverIdx == null && (
          <StarHalf fill="#FFD447" size={20} strokeWidth={0} />
        )}
      </div>
    </div>
  );
};
