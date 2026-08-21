import {
  Disc,
  Music4,
  ListMusic,
  ChevronRight,
  UsersIcon,
  Heart,
  CloudSync,
  CloudUpload,
  Wrench,
  Home,
} from "lucide-react";
import { cn } from "@heroui/react";
import { ReactNode, useState, useEffect } from "react";
import clsx from "clsx";
import { useNavigate } from "react-router-dom";
import { Chip, Header, ListBox } from "@heroui/react";

import Logo from "@/icons/logo.svg?react";

import { Searchbar } from "./Searchbar";
import { ServerStatus } from "./ServerStatus";

import { useApi } from "@/api/ApiContext";
import { useSession } from "@/session/SessionContext";
import { UserRoles } from "@/types/roles";
import { Event } from "@/types/events";

interface IconWrapperProps {
  children: ReactNode;
  className?: string;
}

const ServerIconWrapper = ({ children, className }: IconWrapperProps) => (
  <div
    className={cn(
      className,
      "flex items-center rounded-small justify-center w-6 h-6 bg-background/90 p-1 text-foreground",
    )}
  >
    {children}
  </div>
);

interface ItemCounterProps {
  number?: number;
}

const ItemCounter = ({ number }: ItemCounterProps) => (
  <div className="hidden gap-1 items-center lg:flex text-default-600 fill-foreground text-foreground">
    {number !== undefined && (
      <Chip color="accent" size="sm" variant="primary">
        {number.toString()}
      </Chip>
    )}
    <ChevronRight size={14} />
  </div>
);

export const Menu = () => {
  const [artistsCount, setArtistsCount] = useState(0);
  const [albumCount, setAlbumCount] = useState(0);
  const [trackCount, setTrackCount] = useState(0);
  const [likedTrackCount, setLikedTrackCount] = useState(0);
  const [playlistCount, setPlaylistCount] = useState(0);
  const [syncInProgress, setSyncInProgress] = useState(false);

  const mediaItems = [
    { key: "Artists", href: "/artists", icon: <UsersIcon />, count: artistsCount },
    { key: "Albums", href: "/albums", icon: <Disc />, count: albumCount },
    { key: "Tracks", href: "/tracks", icon: <Music4 />, count: trackCount },
  ];
  const homeItems = [
    {
      key: "Playlists",
      href: "/playlists",
      icon: <ListMusic />,
      count: playlistCount,
    },
    {
      key: "Liked Tracks",
      href: "/liked-tracks",
      icon: <Heart />,
      count: likedTrackCount,
    },
  ];

  const api = useApi();
  const session = useSession();
  const navigate = useNavigate();

  const loadData = () => {
    api.countArtists().then(setArtistsCount);
    api.countAlbums().then(setAlbumCount);
    api.countTracks().then(setTrackCount);
    api.countLikedTracks().then(setLikedTrackCount);
    api.countPlaylists().then(setPlaylistCount);
  };

  useEffect(() => {
    loadData();
  }, [api]);

  useEffect(() => {
    const syncF = api.eventSubscribe(Event.EntitiesUpdated, () => {
      loadData();
    });

    return () => {
      syncF();
    };
  }, []);

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

  return (
    <div className="flex overflow-visible flex-col gap-0 justify-between p-0 py-4 border-r-1 border-border bg-surface-secondary w-[80px] lg:w-[250px] lg:min-w-[250px]">
      <div>
        <div className="flex items-center pl-5 fill-foreground text-foreground">
          <Logo className="" height={40} width={40} />
          <span className="hidden ml-2 text-xl font-bold lg:block">
            Brage Music
          </span>
        </div>
        <div className="px-4 mt-10">
          <Searchbar />
        </div>
        <ListBox
          aria-label="servermenu"
          className="overflow-visible gap-0 p-0 pr-2 pb-5 pl-4 mt-10 bg-none"
          // itemClasses={{
          //   base: "px-3 rounded-none gap-3 h-8 data-[hover=true]:bg-default-100/80",
          // }}
          disabledKeys={disabledServerKeys()}
          onAction={(key) => {
            if (key === "synchronise" && !syncInProgress) {
              api.syncLibrary();
              return;
            } else if (key === "import_music") {
              navigate("/import");
              return;
            } else if (key === "admin") {
              navigate("/server-admin");
              return;
            } else if (key === "home") {
              navigate("/");
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
          }}
        >
          <ListBox.Item className="justify-center" id="home" key="home">
            <ServerIconWrapper><Home/></ServerIconWrapper>
            <p className="hidden flex-grow lg:block text-foreground">
             Home
            </p>
          </ListBox.Item>
          {homeItems.map((item) => (
            <ListBox.Item className="justify-center" id={item.key} key={item.key}>
              <ServerIconWrapper>{item.icon}</ServerIconWrapper>
              <p className="hidden flex-grow lg:block text-foreground">
                {item.key}
              </p>
              <ItemCounter number={item.count} />
            </ListBox.Item>
          ))}
          <ListBox.Section>
            <Header>Media</Header>
            {mediaItems.map((item) => (
              <ListBox.Item className="justify-center" id={item.key} key={item.key}>
                <ServerIconWrapper>{item.icon}</ServerIconWrapper>
                <p className="hidden flex-grow lg:block text-foreground">
                  {item.key}
                </p>
                <ItemCounter number={item.count} />
              </ListBox.Item>
            ))}
          </ListBox.Section>

          <ListBox.Section className="text-foreground">
            <Header>Server</Header>
            {session.supportsSync ? (
              <ListBox.Item
                className={clsx(
                  syncInProgress ? "text-success" : "",
                  "justify-center",
                )}
                id="synchronise"
                key="synchronise"
              >
                <ServerIconWrapper>
                  <CloudSync />
                </ServerIconWrapper>
                <p className="hidden flex-grow lg:block text-foreground">
                  {syncInProgress ? "Synchronising..." : "Synchronise"}
                </p>
              </ListBox.Item>
            ) : null}
            {session.user?.hasAnyRole(
              UserRoles.Admin,
              UserRoles.ImporterWrite,
            ) ? (
              <ListBox.Item className="justify-center" id="import_music" key="import_music">
                <ServerIconWrapper>
                  <CloudUpload />
                </ServerIconWrapper>
                <p className="hidden flex-grow lg:block">Import Music</p>
              </ListBox.Item>
            ) : null}
            {session.user?.hasAnyRole(UserRoles.Admin) ? (
              <ListBox.Item className="justify-center" id="admin" key="admin">
                <ServerIconWrapper>
                  <Wrench />
                </ServerIconWrapper>
                <p className="hidden flex-grow lg:block">Administration</p>
              </ListBox.Item>
            ) : null}
          </ListBox.Section>
        </ListBox>
      </div>
      <div className="flex flex-col gap-4 lg:px-4 px-4.5">
        <ServerStatus serverInfo={session.serverInfo} />
      </div>
    </div>
  );
};
