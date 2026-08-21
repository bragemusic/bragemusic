import {
  KeyRound,
  Laptop,
  LucideIcon,
  Monitor,
  ShieldX,
  Smartphone,
  Speaker,
  Trash2,
  Tv,
} from "lucide-react";

import { types } from "@/types/core";

import { Card } from "./Card";
import { useApi } from "@/api/ApiContext";

interface DeviceCardProps {
  device: types.DeviceDetailed;
}

export const DeviceCard: React.FC<DeviceCardProps> = ({ device }) => {
  const api = useApi();

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
    <div className="flex flex-col">
      <Card
        description={device.active
            ? "Active now"
            : `Last seen: ${new Date(device.last_seen).toDateString()}`
                    }
        fallbackIcon={Icon}
        radius="lg"
        title={device.name}
        actionButtons={[
          {
            confirm: true,
            icon: Trash2,
            size: "sm",
            tooltip: "Remove Device",
            variant: "danger",
            onClick: () => {
              api.removeDevice(device.id);
            },
          },
          {
            confirm: true,
            icon: KeyRound,
            size: "sm",
            tooltip: "Remove Token Used by Device",
            variant: "danger",
            onClick: () => {
              api.removeDeviceToken(device.id);
            },
          },
          {
            confirm: true,
            icon: ShieldX,
            size: "sm",
            tooltip: "Remove Device and Token",
            variant: "danger",
            onClick: () => {
              api.removeDeviceAndToken(device.id);
            },
          },
        ]}
      />
    </div>
  )
  // return (
  //   <Card className="flex flex-col gap-0 p-0 cursor-default select-none border-1 border-border w-[250px] h-[250px]">
  //     <div className="flex flex-grow justify-center items-center p-2 bg-default text-foreground">
  //       <Icon size={70} />
  //     </div>
  //     <Card.Footer className="flex overflow-hidden flex-col justify-between items-center py-1 w-full border-border border-t-1 text-foreground">
  //       <p className="text-lg font-bold">{device.name}</p>
  //       <p className="text-sm">
  //         {device.active
  //           ? "Active now"
  //           : `Last seen: ${new Date(device.last_seen).toDateString()}`}
  //       </p>
  //       <p className="text-sm">{device.last_ip}</p>
  //       <div className="flex gap-2 justify-center py-2 w-full">
  //         <ActionButton
  //           confirm={true}
  //           icon={Trash2}
  //           size="sm"
  //           tooltip="Remove Device"
  //           variant="danger"
  //           onClick={() => {
  //             api.removeDevice(device.id);
  //           }}
  //         />
  //         <ActionButton
  //           confirm={true}
  //           icon={KeyRound}
  //           size="sm"
  //           tooltip="Remove Token Used by Device"
  //           variant="danger"
  //           onClick={() => {
  //             api.removeDeviceToken(device.id);
  //           }}
  //         />
  //         <ActionButton
  //           confirm={true}
  //           icon={ShieldX}
  //           size="sm"
  //           tooltip="Remove Device and Token"
  //           variant="danger"
  //           onClick={() => {
  //             api.removeDeviceAndToken(device.id);
  //           }}
  //         />
  //       </div>
  //     </Card.Footer>
  //   </Card>
  // );
};
