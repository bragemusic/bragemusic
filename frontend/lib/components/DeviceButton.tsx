import { Button, ButtonProps, Popover, Tooltip } from "@heroui/react";
import {
  type LucideIcon,
  Speaker,
  Laptop,
  PauseCircle,
  PlayCircle,
  Link,
  Monitor,
  Smartphone,
  Tv,
} from "lucide-react";
import { useEffect, useState } from "react";
import { cn } from "@heroui/react";

import { types } from "@/types/core";

import { useDeviceStore } from "@/store/deviceStore.ts";
import { useApi } from "@/api/ApiContext";
import { Event } from "@/types/events";

interface DeviceButtonProps {
  minimized?: boolean;
  size?: ButtonProps["size"];
}

export const DeviceButton: React.FC<DeviceButtonProps> = ({minimized=false, size="sm"}) => {
  const [popoverOpen, setPopoverOpen] = useState(false);
  const [selectedDevice, setSelectedDevice] =
    useState<types.DeviceDetailed | null>(null);
  const devices: types.DeviceDetailed[] = Object.values(
    useDeviceStore((s) => s.devices),
  );

  const api = useApi();

  const isEnabled = (): boolean => {
    return true;
  };


  useEffect(() => {
    api.getConnectedDeviceID().then(setDeviceByID);
  }, [api, devices]);

  const setDeviceByID = (id: string|null) => {
    if (id === null) {
      setSelectedDevice(null)
      return
    }

    devices.forEach((d) => {
      if (d.id === id) {
        setSelectedDevice(d)
      }
    })
  }


  useEffect(() => {
      const unsub1 = api.eventSubscribe(Event.DeviceConnectionID, (id: string|null) => {
        setDeviceByID(id);
      });
      return () => {
        unsub1()
      };
  }, [api]);


  const deviceClicked = (d: types.DeviceDetailed) => {
    if (selectedDevice?.id === d.id) {
      api.disconnectDevice();
      return;
    }

    setSelectedDevice(d);
    api.connectDevice(d.id);
  };

  const popoverContent = (
    <Popover.Content className="overflow-y-scroll w-[240px] max-h-[500px]">
      <Popover.Arrow />
      <Popover.Dialog>
        <Popover.Heading />
        <div className="py-2 px-1 w-full">
          <div className="flex flex-col gap-4">
            {devices.map(
              (d) =>
                d.active &&
                d.supports_playback && (
                  <DeviceInfo
                    key={d.id}
                    currentlyConnected={selectedDevice?.id === d.id}
                    device={d}
                    onClick={deviceClicked}
                  />
                ),
            )}
          </div>
        </div>
      </Popover.Dialog>
    </Popover.Content>
  );

  return (
    <Tooltip closeDelay={100} delay={300}>
      <span className="inline-flex">
        <Popover
          isOpen={popoverOpen}
          onOpenChange={(open) => {
            if (!open) {
              setPopoverOpen(false);

              return;
            }

            if (isEnabled()) {
              setPopoverOpen(true);
            }
          }}
        >
          <Popover.Trigger>
            <Button
              className={cn(selectedDevice !== null ? "" : "")}
              isDisabled={!isEnabled()}
              isIconOnly={selectedDevice === null}
              size={size}
              variant={selectedDevice !== null ? "primary" : "ghost"}
            >
              <Speaker
                size={20}
              />
                {selectedDevice !== null && !minimized &&
              <p className="text-xs">
                  Playing from {selectedDevice.name}
              </p>
                  }
            </Button>
          </Popover.Trigger>
          {popoverContent}
        </Popover>
      </span>
      <Tooltip.Content>
        <p>Connect to a device</p>
      </Tooltip.Content>
    </Tooltip>
  );
};

interface DeviceInfoProps {
  device: types.DeviceDetailed;
  currentlyConnected: boolean;
  onClick: (d: types.DeviceDetailed) => void;
}

const DeviceInfo: React.FC<DeviceInfoProps> = ({
  device,
  currentlyConnected,
  onClick,
}) => {
  const deviceIcon = (iconName: string): LucideIcon => {
    switch (iconName) {
      case "laptop":
        return Laptop;
      case "computer":
        return Monitor;
      case "phone":
        return Smartphone;
      case "speaker":
        return Speaker;
      case "tv":
        return Tv;
      default:
        return Laptop;
    }
  };

  const Icon = deviceIcon(device.icon);

  return (
    <div
      className={cn(
        "flex gap-2 items-center rounded-md cursor-pointer select-none",
        currentlyConnected
          ? "bg-accent/60 hover:bg-accent/70"
          : "hover:bg-background",
      )}
      onClick={() => {
        onClick(device);
      }}
    >
      <div className="p-2 rounded-md bg-background/80 text-foreground">
        <Icon size={30} />
      </div>
      <div className="flex flex-col">
        <div className="flex gap-1 items-center text-foreground">
          <p>{device.name}</p>
          {currentlyConnected && (
            <Link className="text-foreground/90" size={12} />
          )}
        </div>
        {device.player_state ? (
          <div className="flex gap-1 items-center">
            {device.player_state?.playback?.playing ? (
              <>
                <PlayCircle className="text-foreground/60" size={15} />
                <p className="text-sm text-foreground/60">Playing</p>
              </>
            ) : (
              <>
                <PauseCircle className="text-foreground/60" size={15} />
                <p className="text-sm text-foreground/60">Paused</p>
              </>
            )}
          </div>
        ) : (
          <p className="text-sm text-content4">Not playing</p>
        )}
      </div>
    </div>
  );
};
