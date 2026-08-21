import { Button, Modal } from "@heroui/react";
import { Input } from "@/primitives/Inputs";
import { useApi } from "@/api/ApiContext";
import { useState } from "react";
import { responses } from "@/types/server";
import { musicbrainz } from "@/types/core";

interface MusicbrainzSearchProps {
    setMbID: React.Dispatch<React.SetStateAction<string|undefined>>;
    setShowMbSearch: React.Dispatch<React.SetStateAction<boolean>>;
}

export const MusicbrainzSearch: React.FC<MusicbrainzSearchProps> = ({ setMbID, setShowMbSearch }) => {
    const api = useApi()

    const [mbSearchResults, setMbSearchResults] = useState<responses.ListPayload<musicbrainz.SearchResults>>()
    const [artist, setArtist] = useState("")
    const [album, setAlbum] = useState("")

    return (
        <>
        <Modal.Body className="px-1 pt-4">
            <div className="flex flex-col gap-2">
                <Input
                    className="flex-grow"
                    color="secondary"
                    label="Artist"
                    name="artist"
                    description="Search for artist name"
                    onChange={(e) => {
                        setArtist(e ? e : "");
                    }}
                    value={artist}
                />
                <Input
                    className="flex-grow"
                    color="secondary"
                    label="Album"
                    name="album"
                    description="Search for album name"
                    onChange={(e) => {
                        setAlbum(e ? e : "");
                    }}
                    value={album}
                />

                    <Button variant="secondary" onClick={() => {api.searchMusicBrainz(artist, album).then(setMbSearchResults)}}>
                    Search
                    </Button>
                <div className="flex flex-col gap-2">
                {mbSearchResults && mbSearchResults.items?.map((item) => (
                    <div
                        className="flex flex-row gap-3 justify-between items-center p-2 rounded-lg cursor-pointer bg-background border-border border-1 hover:bg-accent hover:text-accent-foreground"
                        onClick={() => {setMbID(item.id); setShowMbSearch(false)}}
                    >
                        <div className="flex flex-col">
                            <p>{item.album}</p>
                            <p className="capitalize text-small text-foreground/60">{item.artist}</p>
                        </div>
                        <p>{item.score}% match</p>
                    </div>
                ))}
                </div>
            </div>
        </Modal.Body>
        <Modal.Footer>
            <Button variant="secondary" onClick={() => {setShowMbSearch(false)}}>
                Back
            </Button>
        </Modal.Footer>
        </>
    )
}
