import { useEffect, useRef, useState } from "react";
import { cn } from "@heroui/react";
import { canvasRGB } from "stackblur-canvas";

interface BlurredBackdropProps {
  /** Image URL, e.g. albumImageLink(album.id, 320). Nothing renders without it. */
  src?: string;
  /** Colour boost. ~1.3 suits album art, ~1.9 suits darker artist photos. */
  saturate?: number;
  /** How far down the page the colour reaches before fading into the background. */
  height?: number;
  className?: string;
}

/** Width of the offscreen canvas. The blur removes all detail, so small is fine. */
const CANVAS_WIDTH = 96;
/** StackBlur radius in canvas pixels. ~12 matches the old blur(90px) at xl widths. */
const BLUR_RADIUS = 12;
const BRIGHTNESS = 1.1;

/**
 * Faded, heavily blurred copy of an image behind the top of a page.
 *
 * Place it as the first child of a `relative isolate` wrapper around the page
 * content; it sits behind everything (`-z-10`) and scrolls away with the page.
 *
 * The blur is done once in a small canvas instead of with CSS `filter: blur()`:
 * WebKit (Safari, WebKitGTK) silently drops CSS filters on very large layers,
 * which hid the backdrop at wide screen sizes.
 */
export function BlurredBackdrop({
  src,
  saturate = 1.3,
  height = 620,
  className,
}: BlurredBackdropProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  // Track *which* src is ready/failed, so navigating resets state without an effect.
  const [readySrc, setReadySrc] = useState<string>();
  const [failedSrc, setFailedSrc] = useState<string>();

  useEffect(() => {
    if (!src) return;
    let cancelled = false;

    const img = new Image();
    img.crossOrigin = "anonymous"; // needed to read pixels back from the canvas
    img.decoding = "async";
    img.onload = () => {
      const canvas = canvasRef.current;
      if (cancelled || !canvas) return;
      try {
        drawBlurred(canvas, img, saturate);
        setReadySrc(src);
      } catch {
        setFailedSrc(src); // e.g. canvas tainted by a cross-origin image without CORS
      }
    };
    img.onerror = () => !cancelled && setFailedSrc(src);
    img.src = src;

    return () => {
      cancelled = true;
    };
  }, [src, saturate]);

  if (!src || failedSrc === src) {
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
      <canvas
        ref={canvasRef}
        className={cn(
          "absolute -top-1/2 -left-[10%] h-[170%] w-[120%] max-w-none object-cover",
          "transition-opacity duration-700 ease-out motion-reduce:transition-none",
          readySrc === src ? "opacity-55" : "opacity-0",
        )}
      />
      {/* Fade into the theme background, so it works in light and dark mode */}
      <div className="absolute inset-0 bg-linear-to-b from-background/10 via-background/65 via-60% to-background" />
    </div>
  );
}

/** Draws `img` small, blurs it, then applies CSS-equivalent saturate + brightness. */
function drawBlurred(canvas: HTMLCanvasElement, img: HTMLImageElement, saturate: number) {
  const w = CANVAS_WIDTH;
  const h = Math.max(1, Math.round((w * img.naturalHeight) / img.naturalWidth));
  canvas.width = w;
  canvas.height = h;

  const ctx = canvas.getContext("2d", { willReadFrequently: true });
  if (!ctx) throw new Error("no 2d context");

  ctx.drawImage(img, 0, 0, w, h);
  canvasRGB(canvas, 0, 0, w, h, BLUR_RADIUS);

  // Same colour matrices as CSS saturate() and brightness().
  const s = saturate;
  const data = ctx.getImageData(0, 0, w, h);
  const d = data.data;
  for (let i = 0; i < d.length; i += 4) {
    const r = d[i], g = d[i + 1], b = d[i + 2];
    d[i]     = ((0.213 + 0.787 * s) * r + (0.715 - 0.715 * s) * g + (0.072 - 0.072 * s) * b) * BRIGHTNESS;
    d[i + 1] = ((0.213 - 0.213 * s) * r + (0.715 + 0.285 * s) * g + (0.072 - 0.072 * s) * b) * BRIGHTNESS;
    d[i + 2] = ((0.213 - 0.213 * s) * r + (0.715 - 0.715 * s) * g + (0.072 + 0.928 * s) * b) * BRIGHTNESS;
  }
  ctx.putImageData(data, 0, 0);
}
