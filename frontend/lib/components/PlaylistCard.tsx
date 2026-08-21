import { useEffect, useState } from "react";
import { ListMusic, Sparkles } from "lucide-react";
import { useNavigate } from "react-router-dom";

import { types } from "@/types/core";

import { Card } from "./Card";

import { useApi } from "@/api/ApiContext";
import { playlistImageLinkFromID } from "@/util/images";

interface PlaylistCardProps {
  playlist: types.Playlist;
}

export const PlaylistCard: React.FC<PlaylistCardProps> = ({ playlist }) => {
  const [trackCount, setTrackCount] = useState(0);

  const api = useApi();
  const navigate = useNavigate();

  useEffect(() => {
    api.countPlaylistTracks(playlist.id).then(setTrackCount);
  }, [api, playlist]);

  return (
    <Card
      description={playlist.type === "smart" ? "Smart playlist" : `${trackCount} track${trackCount != 1 ? "s " : ""}`}
      fallbackIcon={ListMusic}
      imgSrc={playlistImageLinkFromID(playlist.id.toString(), 320)}
      radius="xl"
      title={playlist.name}
      overlayIcon={playlist.type == "smart" ? Sparkles : undefined}
      onClick={() => {
        navigate("/playlists/" + playlist.id);
      }}
    />
  );
};
