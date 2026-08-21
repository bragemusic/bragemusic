import { create } from "zustand";

import { types } from "@/types/core";

import { Event } from "@/types/events";
import { Api } from "@/api/Api";

type DeviceMap = Record<string, types.DeviceDetailed>;

type DeviceStore = {
  devices: DeviceMap;

  init: (api: Api) => Promise<void>;
  loadDevices: (api: Api) => Promise<void>;
  setDevices: (devices: types.DeviceDetailed[]) => void;

  updatePlayback: (deviceId: string, pb: types.PlaybackStateDTO) => void;
  updateContext: (deviceId: string, ctx: types.PlayContextDTO) => void;
};

export const useDeviceStore = create<DeviceStore>((set, get) => ({
  devices: {},

  setDevices: (devices) => {
    console.log("Fetched devices:", devices);
    const map: DeviceMap = {};

    for (const d of devices) {
      map[d.id] = d;
    }

    set({ devices: map });
  },

  updatePlayback: (deviceId, pb) =>
    set((state) => {
      const device = state.devices[deviceId];

      if (!device || !device.player_state) return state;

      // Mutate the existing PlayerStateDTO instance
      device.player_state.playback = pb;

      return {
        devices: {
          ...state.devices,
          [deviceId]: device,
        },
      };
    }),

  updateContext: (deviceId, ctx) =>
    set((state) => {
      const device = state.devices[deviceId];

      if (!device || !device.player_state) return state;

      // Mutate the existing PlayerStateDTO instance
      device.player_state.context = ctx;

      return {
        devices: {
          ...state.devices,
          [deviceId]: device,
        },
      };
    }),

  loadDevices: async (api: Api) => {
    const { setDevices } = get();

    const devices = await api.listDevices();

    console.log("Fetched devices:", devices);
    setDevices(devices);
  },

  init: async (api: Api) => {
    const { setDevices, updatePlayback, updateContext } = get();

    // const devices = await api.listDevices();
    // console.log("Fetched devices:", devices);
    // setDevices(devices);

    api.eventSubscribe(Event.DeviceConnected, async () => {
      const devices = await api.listDevices();

      console.log("Client connected devices:", devices);
      setDevices(devices);
    });

    api.eventSubscribe(Event.DeviceUpdated, async () => {
      const devices = await api.listDevices();

      console.log("Client connected devices:", devices);
      setDevices(devices);
    });

    api.eventSubscribe(Event.DeviceDisconnected, async () => {
      const devices = await api.listDevices();

      console.log("Client disconnected devices:", devices);
      setDevices(devices);
    });

    api.eventSubscribe(Event.InternalClientAllReplaced, (d: types.DeviceDetailed[]) => {
      setDevices(d);
    });

    api.eventSubscribe(Event.DevicePlaybackState, (pb: types.PlaybackStateDTO) => {
      updatePlayback(pb.device_id, pb);
    });

    api.eventSubscribe(Event.DevicePlayContext, (ctx: types.PlayContextDTO) => {
      updateContext(ctx.device_id, ctx);
    });
  },
}));
