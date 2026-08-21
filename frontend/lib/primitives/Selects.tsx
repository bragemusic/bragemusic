import {
  Avatar,
  Chip,
  Select as BaseSelect,
  ListBox,
  Label,
  Description,
  Key,
} from "@heroui/react";

import { types } from "@/types/core";

import { artistImageLink } from "@/util/images";

type ArtistSelectProps = {
  name: string;
  artists: types.Artist[];
  selectedArtistIds: string[] | undefined;
  description?: string;
  onChange?: (key: Key[]) => void;
};

export function ArtistSelect({
  name,
  artists,
  selectedArtistIds,
  description,
  onChange,
}: ArtistSelectProps) {
  return (
    <BaseSelect
      defaultValue={selectedArtistIds}
      name={name}
      placeholder="Select artist(s)"
      selectionMode="multiple"
      variant="secondary"
      onChange={onChange}
    >
      <Label>Artists</Label>
      <BaseSelect.Trigger>
        <BaseSelect.Value className="overflow-hidden">
          {({ defaultChildren, isPlaceholder, state }) => {
            if (isPlaceholder || state.selectedItems.length === 0) {
              return defaultChildren;
            }
            const selectedItems = state.selectedItems;

            return (
              <div className="flex gap-2 items-center">
                {selectedItems.map((item) => {
                  return (
                    <Chip key={item.key} color="accent" variant="primary">
                      {item.textValue}
                    </Chip>
                  );
                })}
              </div>
            );
          }}
        </BaseSelect.Value>
        <BaseSelect.Indicator />
      </BaseSelect.Trigger>
      <BaseSelect.Popover>
        <ListBox>
          {artists.map((artist) => (
            <ListBox.Item
              key={artist.id.toString()}
              id={artist.id.toString()}
              textValue={artist.name}
            >
              <Avatar size="sm">
                <Avatar.Image src={artistImageLink(artist, 320)} />
                <Avatar.Fallback>{artist.name}</Avatar.Fallback>
              </Avatar>
              <div className="flex flex-col">
                <Label>{artist.name}</Label>
              </div>
              <ListBox.ItemIndicator className="text-accent" />
            </ListBox.Item>
          ))}
        </ListBox>
      </BaseSelect.Popover>
      {description && <Description>{description}</Description>}
    </BaseSelect>
    // <BaseSelect
    //   name={name}
    //   classNames={{
    //     base: "",
    //     trigger: "min-h-12 py-2",
    //    popoverContent: "bg-background"
    //   }}
    //   isMultiline={true}
    //   items={artists}
    //   label="Artists"
    //   labelPlacement="inside"
    //   placeholder="Select one or more artists"
    //   defaultSelectedKeys={selectedArtistIds}
    //   renderValue={(items) => {
    //     return (
    //       <div className="flex flex-wrap gap-2">
    //         {items.map((item) => (
    //           <Chip key={item.key} color="accent">{item.data?.name}</Chip>
    //         ))}
    //       </div>
    //     );
    //   }}
    //   selectionMode="multiple"
    //   variant="flat"
    //   color="default"
    //   required
    // >
    //   {(artist) => (
    //     <SelectItem key={artist.id.toString()} textValue={artist.name}>
    //       <div className="flex gap-2 items-center">
    //         <Avatar
    //           size="sm"
    //           className="shrink-0"
    //           color="accent"
    //         >
    //           <Avatar.Image
    //             src={artistImageLink(artist, 320)}
    //           />
    //           <Avatar.Fallback>{artist.name}</Avatar.Fallback>
    //         </Avatar>

    //         <div className="flex flex-col">
    //           <span className="text-small">{artist.name}</span>
    //         </div>
    //       </div>
    //     </SelectItem>
    //   )}
    // </BaseSelect>
  );
}
