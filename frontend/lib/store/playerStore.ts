import { create } from "zustand";

import { PlayerState } from "@/models/PlayerState";
import { types } from "@/types/core";
import { Event } from "@/types/events.ts";
import { Api } from "@/api/Api";

type PlayerStore = {
  player: PlayerState;
  init: (api: Api) => Promise<void>;

  setContext: (ctx: types.PlayContext) => void;
  setPlayback: (pb: types.PlaybackState) => void;
};

export const usePlayerState = create<PlayerStore>((set, get) => ({
  player: new PlayerState(),

  init: async (api: Api) => {
    const { setContext, setPlayback } = get();

    api.eventSubscribe(Event.PlayerContextChange, (ctx: types.PlayContext) => {
      setContext(ctx);
    });

    api.eventSubscribe(Event.PlayerPlaybackChange, (pb: types.PlaybackState) => {
      setPlayback(pb);
    });
  },

  setContext: (ctx) =>
    set((state) => {
      const next = new PlayerState();

      Object.assign(next, state.player);
      next.context = ctx;

      return { player: next };
    }),

  setPlayback: (pb) =>
    set((state) => {
      const next = new PlayerState();

      Object.assign(next, state.player);
      next.playback = pb;

      return { player: next };
    }),
}));
