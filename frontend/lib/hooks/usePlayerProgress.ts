import { useEffect, useState } from "react";

import { usePlayerState } from "@/store/playerStore";

export function usePlayerProgress() {
  const playback = usePlayerState((s) => s.player.playback);
  const [progressMs, setProgressMs] = useState(0);

  useEffect(() => {
    if (!playback) return;

    const calc = () => {
      const base = playback.progress ?? 0;
      const updatedAt = Date.parse(playback.updated_at ?? "") || Date.now();

      const now = Date.now();
      const delta = playback.playing ? now - updatedAt : 0;

      const value = base + delta;

      setProgressMs(Number.isFinite(value) ? value : 0);
    };

    calc();

    if (!playback.playing) return;

    const interval = setInterval(calc, 250);

    return () => clearInterval(interval);
  }, [playback]);

  return progressMs;
}
