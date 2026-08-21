import { BragErr } from "@/types/bragerr";
import { types } from "@/types/core";

export const makeTimestamp = (ms: number): string => {
  const s = Math.floor(ms / 1000);
  const minutes = Math.floor((s % 3600) / 60);
  const remainingSeconds = s % 60;

  return minutes + ":" + String(remainingSeconds).padStart(2, "0");
};

export function isBragErr(obj: any): obj is BragErr {
  return (
    obj &&
    typeof obj === "object" &&
    typeof obj.code === "string" &&
    typeof obj.title === "string" &&
    typeof obj.message === "string"
  );
}

export const parseDate = (ds: string | undefined): Date | undefined => {
  if ( ds === undefined ) {
    return undefined
  }

  return new Date(ds)
}

export const dateHR = (date: Date | undefined): string => {
  if ( date === undefined ) {
    return ""
  }

  const day = String(date.getDate()).padStart(2, "0");
  const month = date.toLocaleString(undefined, { month: "short" }); // browser locale tz
  const year = date.getFullYear();

  const hours = String(date.getHours()).padStart(2, "0");
  const mins = String(date.getMinutes()).padStart(2, "0");
  const secs = String(date.getSeconds()).padStart(2, "0");

  return `${day} ${month} ${year} ${hours}:${mins}:${secs}`;
}

export type TrackQualityLabel = "Lossy" | "CD" | "HI-FI" | "Studio"

export const getTrackQualityLabel = (mediafile: types.MediaFile): TrackQualityLabel => {
  // This must update when other formats are supported
  if (mediafile.codec != "flac") {
    return "Lossy"
  } else if (mediafile.sample_rate < 48000 ) {
    return "CD"
  } else if ( mediafile.sample_rate < 192000 ) {
    return "HI-FI"
  } else {
    return "Studio"
  }
}

export const getArtistsString = (artists?: types.ArtistMinimal[]): string => {
  if (!artists) {
    return ""
  }
  return artists.map((a) => (a.name)).join(", ")
}
