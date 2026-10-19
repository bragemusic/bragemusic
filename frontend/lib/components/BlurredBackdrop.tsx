import { useEffect, useState } from "react";
import { cn } from "@heroui/react";

interface BlurredBackdropProps {
  /** Image URL, e.g. albumImageLink(album.id, 320). Nothing renders without it. */
  src?: string;
  /** Colour boost. ~1.3 suits album art, ~1.9 suits darker artist photos. */
  saturate?: number;
  /** How far down the page the colour reaches before fading into the background. */
  height?: number;
  className?: string;
}

/**
 * Faded, heavily blurred copy of an image behind the top of a page.
 *
 * Place it as the first child of a `relative isolate` wrapper around the page
 * content; it sits behind everything (`-z-10`) and scrolls away with the page.
 * The blur destroys all detail, so the 320px image is enough — and it is
 * usually already cached because the cover/avatar uses the same size.
 */
export function BlurredBackdrop({
  src,
  saturate = 1.3,
  height = 620,
  className,
}: BlurredBackdropProps) {
  const [loaded, setLoaded] = useState(false);
  const [failed, setFailed] = useState(false);

  // Reset when navigating between albums/artists (the page stays mounted).
  useEffect(() => {
    setLoaded(false);
    setFailed(false);
  }, [src]);

  if (!src || failed) {
    return null;
  }

  return (
    <div
      aria-hidden
      className={cn(
        "pointer-events-none absolute inset-x-0 top-0 -z-10 overflow-hidden",
        className,
      )}
      style={{ height }}
    >
      <img
        alt=""
        className={cn(
          "absolute -top-1/2 -left-[10%] h-[170%] w-[120%] max-w-none object-cover",
          "transition-opacity duration-700 ease-out motion-reduce:transition-none",
          loaded ? "opacity-55" : "opacity-0",
        )}
        decoding="async"
        src={src}
        style={{ filter: `blur(90px) saturate(${saturate}) brightness(1.1)` }}
        onError={() => setFailed(true)}
        onLoad={() => setLoaded(true)}
      />
      {/* Fade into the theme background, so it works in light and dark mode */}
      <div className="absolute inset-0 bg-linear-to-b from-background/10 via-background/65 via-60% to-background" />
    </div>
  );
}
