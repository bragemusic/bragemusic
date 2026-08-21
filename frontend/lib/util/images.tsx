import { types } from "@/types/core";

export type ImageSize = 320 | 640 | 1024 | 1600 | 2400;

export function albumImageLink(albumID: string, size: ImageSize): string {
  return "/api/img/albums/" + albumID + "/" + size + ".jpg";
}

export function artistImageLink(artist: types.Artist, size: ImageSize): string {
  return "/api/img/artists/" + artist.id + "/" + size + ".jpg";
}

export function artistImageLinkFromID(
  artistID: string,
  size: ImageSize,
): string {
  return "/api/img/artists/" + artistID + "/" + size + ".jpg";
}

export function albumImageLinkFromID(albumID: string, size: ImageSize): string {
  return "/api/img/albums/" + albumID + "/" + size + ".jpg";
}

export function playlistImageLinkFromID(id: string, size: ImageSize): string {
  return "/api/img/playlists/" + id + "/" + size + ".jpg";
}

export function userAvatarLink(id: string, size: ImageSize): string {
  return "/api/img/users/" + id + "/" + size + ".jpg";
}

