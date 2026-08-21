import EditModal from "./EditModal";
import { Input, Textarea } from "@/primitives/Inputs";
import { Text, useOverlayState } from "@heroui/react";
import { TracksFilterPane } from "./TracksFilterPane";
import { types } from "@/types/core";
import { useState } from "react";
import { useApi } from "@/api/ApiContext";

type State = ReturnType<typeof useOverlayState>;

interface SmartPlaylistModalProps {
    state: State;
    artists: types.Artist[];
    initFilter?: types.TrackFilter;
}


export const SmartPlaylistModal: React.FC<SmartPlaylistModalProps> = ({ state, artists, initFilter = {mood: {}} as types.TrackFilter }) => {
    const [filter, setFilter] = useState<types.TrackFilter>({mood: {}} as types.TrackFilter)

    const api = useApi();

    const onAddSmartPlaylist = (fd: FormData) => {
        const source = {
            name: String(fd.get("name") ?? ""),
            description: fd.get("description") || undefined,
        };

        const plistData = types.PlaylistBase.createFrom(source);

        api.addSmartPlaylist(plistData, filter)
    };

    return (
        <EditModal
            state={state}
            header={`Add Smart Playlist`}
            onSubmit={onAddSmartPlaylist}
        >
            <Input
                name="name"
                label="Name"
                // defaultValue={playlist?.name}
                required
            />
            <Textarea
                name="description"
                label="Description"
                // defaultValue={playlist?.description}
            />
            <Text type="h5" className="pt-4">Filter</Text>
            <TracksFilterPane initFilter={initFilter} onApply={setFilter} artists={artists} dynamicApply/>
        </EditModal>
    )
}
