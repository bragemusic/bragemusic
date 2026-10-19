import {
  Disc,
  Music4,
  ListMusic,
  UsersIcon,
  Heart,
  CloudSync,
  CloudUpload,
  Wrench,
  Home,
  LucideIcon,
  Search,
  Server,
} from "lucide-react";
import { Badge, cn, Key, Separator, useOverlayState } from "@heroui/react";
import { useState, useEffect } from "react";
import clsx from "clsx";
import { useNavigate } from "react-router-dom";
import { ListBox } from "@heroui/react";

import Logo from "@/icons/logo.svg?react";

import { Searchbar } from "./Searchbar";

import { useApi } from "@/api/ApiContext";
import { useSession } from "@/session/SessionContext";
import { UserRoles } from "@/types/roles";
import { Event } from "@/types/events";
import { UserMenu } from "./UserMenu";

interface MenuButtonProps {
  id: string;
  title: string;
  icon: LucideIcon;
  badgeColor?: "accent"| "default"| "success"| "warning"| "danger" ;
  className?:string;
}

const MenuButton = ({ id, title, icon: Icon, className, badgeColor }: MenuButtonProps) => {
  return (
    <ListBox.Item
      id={id}
      textValue={title}
      className={cn("py-2 px-4 text-xs rounded-lg hover:stroke-accent stroke-foreground", className)}
    >
      <Badge.Anchor className="flex flex-col gap-1 justify-center items-center w-full">
        <Icon size={18} className="stroke-inherit"/>
        <p className="text-center">{title}</p>
        {badgeColor && <Badge color={badgeColor} size="sm" variant="primary" className="mr-1 -mt-1" />}
      </Badge.Anchor>
    </ListBox.Item>
  );
};

export const Menu = () => {
  const [syncInProgress, setSyncInProgress] = useState(false);

  const searchState = useOverlayState();

  const mediaItems = [
    { key: "Artists", href: "/artists", icon: UsersIcon },
    { key: "Albums", href: "/albums", icon: Disc },
    { key: "Tracks", href: "/tracks", icon: Music4 },
  ];
  const homeItems = [
    {
      key: "Search",
      href: "",
      icon: Search,
    },
    {
      key: "Home",
      href: "/",
      icon: Home,
    },
    {
      key: "Playlists",
      href: "/playlists",
      icon: ListMusic,
    },
    {
      key: "Liked Tracks",
      href: "/liked-tracks",
      icon: Heart,
    },
  ];

  const api = useApi();
  const session = useSession();
  const navigate = useNavigate();


  useEffect(() => {
    const syncP = api.eventSubscribe(Event.SyncInProgress, (p: boolean) => {
      setSyncInProgress(p);
    });

    return () => {
      syncP();
    };
  }, []);

  const disabledServerKeys = (): string[] => {
    if (!session.serverFullyAvailable) {
      return ["synchronise", "import_music", "admin"];
    }

    if (syncInProgress) {
      return ["synchronise"];
    }

    return [];
  };

  const listActions = (key: Key) => {
    console.log(session.serverInfo.status)
    if (key === "synchronise" && !syncInProgress) {
      api.syncLibrary();
      return;
    } else if (key === "import_music") {
      navigate("/import");
      return;
    } else if (key === "admin") {
      navigate("/server-admin");
      return;
    } else if (key === "Search") {
      searchState.setOpen(true)
      return;
    }

    mediaItems.forEach((item) => {
      if (item.key == key) {
        navigate(item.href);

        return;
      }
    });
    homeItems.forEach((item) => {
      if (item.key == key) {
        navigate(item.href);

        return;
      }
    });
  }

  return (
    <div className="flex overflow-visible flex-col gap-0 justify-between items-center p-0 py-6 bg-surface-secondaryd w-[100px]">
      <div className="flex flex-col flex-grow items-center">
        <Searchbar state={searchState} />
        <Logo className="fill-accent" height={40} width={40} />
        <ListBox
          aria-label="servermenu"
          className="flex overflow-visible flex-col flex-1 flex-grow gap-0 p-0 px-3 pb-5 mt-10 min-h-0 bg-none w-[100px]"
          disabledKeys={disabledServerKeys()}
          onAction={listActions}
        >
          <ListBox.Section className="gap-2 mb-4">
            {homeItems.map((item) => (
              <MenuButton key={item.key} id={item.key} title={item.key} icon={item.icon}/>
            ))}
          </ListBox.Section>

          <Separator />

          <ListBox.Section className="gap-2 mt-4 mb-2">
            {mediaItems.map((item) => (
              <MenuButton key={item.key} id={item.key} title={item.key} icon={item.icon}/>
            ))}
          </ListBox.Section>

          <ListBox.Section className="gap-2 mt-auto mb-2">
            {session.supportsSync ? (
              <MenuButton
                key="synchronise"
                id="synchronise"
                title={syncInProgress ? "Syncing..." : "Sync"}
                icon={CloudSync}
                className={clsx(
                    syncInProgress ? "text-success" : "",
                  )}
              />
            ) : null}
            {session.user?.hasAnyRole(
              UserRoles.Admin,
              UserRoles.ImporterWrite,
            ) ? (
              <MenuButton key="import_music" id="import_music" title="Import" icon={CloudUpload}/>
            ) : null}
            {session.user?.hasAnyRole(UserRoles.Admin) ? (
              <MenuButton key="admin" id="admin" title="Server Admin" icon={Wrench}/>
            ) : null}
            <MenuButton
              key="server_status"
              id="server_status"
              title="Server"
              icon={Server}
              badgeColor={
                session.serverInfo.status == "running" ? "success" :
                  session.serverInfo.status == "no_auth" ? "warning" :
                  session.serverInfo.status == "unavailable" ? "danger" : undefined
              }
            />
          </ListBox.Section>

        </ListBox>
      </div>

      <UserMenu />
    </div>
  );
};
