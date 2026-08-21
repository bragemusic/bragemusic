import { LucideIcon } from "lucide-react";
import { useEffect, useState } from "react";

const radiusClass = (radius: ImageProps["radius"]) => {
  const map = {
    none: "rounded-none",
    xs: "rounded-xs",
    sm: "rounded-sm",
    md: "rounded-md",
    lg: "rounded-lg",
    xl: "rounded-xl",
    "2xl": "rounded-2xl",
    "3xl": "rounded-3xl",
    "4xl": "rounded-4xl",
    full: "rounded-full",
  };

  if (radius == undefined) {
    return map["md"];
  }

  return map[radius];
};

export type ImageProps = {
  src?: string;
  fallbackIcon: LucideIcon;
  height: number;
  width: number;
  onClick?: () => void;
  radius?:
    | "xs"
    | "sm"
    | "md"
    | "lg"
    | "xl"
    | "2xl"
    | "3xl"
    | "4xl"
    | "full"
    | "none";
  className?: string;
  customHeight?: boolean;
};

export function Image({
  src,
  fallbackIcon: Icon,
  height,
  width,
  onClick,
  radius = "md",
  className = "",
  customHeight = false,
}: ImageProps) {
  const [status, setStatus] = useState<"loading" | "ok" | "error">("loading");

  useEffect(() => {
    if (!src) {
      setStatus("error");

      return;
    }

    let cancelled = false;

    setStatus("loading");

    const img = new window.Image();

    img.onload = () => {
      if (!cancelled) setStatus("ok");
    };

    img.onerror = () => {
      if (!cancelled) setStatus("error");
    };

    img.src = src;

    return () => {
      cancelled = true;
    };
  }, [src]);

  if (status === "error" || status === "loading") {
    return (
      <div
        className={
          `bg-default text-foreground flex items-center justify-center ${radiusClass(radius)} ` +
          className
        }
        style={customHeight ? {} : { height, width }}
        onClick={onClick}
      >
        <Icon size={height / 3} />
      </div>
    );
  }

  return (
    <img
      className={
        `object-cover overflow-hidden ${radiusClass(radius)} ` + className
      }
      src={src}
      style={customHeight ? {} : {
        height: `${height}px`,
        width: `${width}px`,
      }}
      onClick={onClick}
    />
  );
}
