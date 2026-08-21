import { createContext, useContext, useEffect, useState } from "react";

import { useApi } from "../api/ApiContext";

import { types } from "@/types/core";
import { UserDetails } from "@/models/UserDetails";
import { Event } from "@/types/events";

export type SessionState = {
  user: UserDetails | null | undefined;
  serverInfo: types.ServerApiInfo;
  serverFullyAvailable: boolean;
  supportsSync: boolean;
};

const SessionContext = createContext<SessionState | null>(null);

export function SessionProvider({ children }: { children: React.ReactNode }) {
  const api = useApi();

  const [user, setUser] = useState<UserDetails | null | undefined>(undefined);
  const [serverStatus, setServerStatus] = useState<types.ServerApiInfo>({
    application: "",
    name: "",
    version: "",
    status: "unknown",
    id: "",
  });
  const [serverFullyAvailable, setServerFullyAvailable] =
    useState<boolean>(false);
  const [supportsSync, setSupportsSync] = useState<boolean>(false);

  useEffect(() => {
    const userEvent = api.eventSubscribe(
      Event.UserUpdated,
      (p: types.UserDetails | null) => {
        setUser(p ? new UserDetails(p) : null);
        if (p === null || p === undefined) {
          api.publicServerStatus().then(setServerStatus);
        } else {
          api.serverStatus().then(setServerStatus);
        }
      },
    );

    const serverS = api.eventSubscribe(Event.ServerStatus, (p: types.ServerApiInfo) => {
      setServerStatus(p);
    });

    return () => {
      userEvent();
      serverS();
    };
  }, []);

  useEffect(() => {
    setServerFullyAvailable(serverStatus.status === "running");
  }, [serverStatus]);

  const loadData = () => {
    api.publicServerStatus().then(setServerStatus);
    api.supportsSync().then(setSupportsSync);
  };

  useEffect(() => {
    loadData();
  }, [api]);

  return (
    <SessionContext.Provider
      value={{
        user: user,
        serverInfo: serverStatus,
        serverFullyAvailable: serverFullyAvailable,
        supportsSync: supportsSync,
      }}
    >
      {children}
    </SessionContext.Provider>
  );
}

export function useSession(): SessionState {
  const ctx = useContext(SessionContext);

  if (!ctx) {
    throw new Error("SessionProvider missing");
  }

  return ctx;
}
