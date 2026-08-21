import clsx from "clsx";

import { types } from "@/types/core";

import { UserMenu } from "./UserMenu";

interface ServerStatusProps {
  serverInfo: types.ServerApiInfo;
  forceBig?: boolean;
}

export const ServerStatus = ({ serverInfo, forceBig=false }: ServerStatusProps) => {
  const serverAvailable = (): boolean => {
    return serverInfo.status === "running";
  };

  const status: Record<string, string> = {
    running: "Running",
    unavailable: "Unavailable",
    no_auth: "No Authentication",
  };

  const getStatus = (key: string) => {
    return status[key];
  };

  return (
    <div className={clsx("flex overflow-hidden flex-col p-1 rounded-sm text-foreground border-1 border-border bg-background", !forceBig ? "lg:flex-row":"flex-row")}>
      <UserMenu />
      <div className={clsx("px-0 py-1 justify-between flex flex-col group-data-[hover=true]:bg-default-200", !forceBig ? "lg:px-3":"px-3")}>
        <span
          className={clsx(
            serverAvailable() ? "text-success" : "text-danger",
            "capitalize text-xl w-full text-center",
            !forceBig ? "lg:hidden" : "hidden",
          )}
        >
          {serverInfo.status ? serverInfo.status[0] : "U"}
        </span>
        <span className={clsx(!forceBig ? "hidden lg:block":"block", "text-sm text-default-800")}>
          {serverAvailable() ? serverInfo.name : "Unknown server"}
        </span>
        <div className={clsx("gap-2 text-sm", !forceBig ? "hidden lg:flex":"flex")}>
          {serverInfo.version && (
            <span className="text-default-600">{serverInfo.version}</span>
          )}
          <span
            className={clsx(
              serverAvailable() ? "text-success" : "text-danger",
              "capitalize",
            )}
          >
            {getStatus(serverInfo.status)}
          </span>
        </div>
      </div>
    </div>
  );
};
