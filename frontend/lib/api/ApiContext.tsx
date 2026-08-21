import { createContext, useContext } from "react";

import { Api } from "../api/Api";

import { PlayerApi } from "./PlayerApi";

const ApiContext = createContext<Api | null>(null);
const PlayerApiContext = createContext<PlayerApi | null>(null);

export function ApiProvider({
  api,
  playerApi,
  children,
}: {
  api: Api;
  playerApi: PlayerApi;
  children: React.ReactNode;
}) {
  return (
    <ApiContext.Provider value={api}>
      <PlayerApiContext.Provider value={playerApi}>
        {children}
      </PlayerApiContext.Provider>
    </ApiContext.Provider>
  );
}

export function useApi(): Api {
  const api = useContext(ApiContext);

  if (!api) {
    throw new Error("ApiProvider missing");
  }

  return api;
}

export function usePlayerApi(): PlayerApi {
  const api = useContext(PlayerApiContext);

  if (!api) {
    throw new Error("ApiProvider missing");
  }

  return api;
}
