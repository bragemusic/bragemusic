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

import { useApi } from "@/api/ApiContext";
import { Table } from "@heroui/react";
import { ActionButton } from "@/primitives/Buttons";

interface DeviceRowProps{
  device: types.DeviceDetailed;
}

export const DeviceRow: React.FC<DeviceRowProps> = ({ device }) => {
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
    <Table.Row key={device.id} id={device.id}>
        <Table.Cell><Icon /></Table.Cell>
        <Table.Cell>{device.name}</Table.Cell>
        <Table.Cell>{
          device.active
            ? "Active now"
            : new Date(device.last_seen).toDateString()
          }
        </Table.Cell>
        <Table.Cell>{device.last_ip}</Table.Cell>
        <Table.Cell>{device.platform}</Table.Cell>
        <Table.Cell>{device.version}</Table.Cell>
        <Table.Cell>{device.supports_playback ? "yes" : "no"}</Table.Cell>
        <Table.Cell className="flex gap-1">
            <ActionButton
              confirm={true}
              icon={Trash2}
              size="sm"
              tooltip="Remove Device"
              variant="danger-soft"
              onClick={() => {
                api.removeDevice(device.id);
              }}
            />
            <ActionButton
              confirm={true}
              icon={KeyRound}
              size="sm"
              tooltip="Remove Token Used by Device"
              variant="danger-soft"
              onClick={() => {
                api.removeDeviceToken(device.id);
              }}
            />
            <ActionButton
              confirm={true}
              icon={ShieldX}
              size="sm"
              tooltip="Remove Device and Token"
              variant="danger-soft"
              onClick={() => {
                api.removeDeviceAndToken(device.id);
              }}
            />
        </Table.Cell>
    </Table.Row>
  );
};
