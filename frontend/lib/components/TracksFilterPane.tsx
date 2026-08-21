import { ArtistSelect } from "@/primitives/Selects";
import { types } from "@/types/core";
import { Button, Text, Label, Slider, Switch } from "@heroui/react";
import { useEffect, useState } from "react";

interface PercentSliderProps {
    label: string;
    onChange: (value: number|undefined) => void;
    defaultValue: number | undefined;
}


const PercentSlider: React.FC<PercentSliderProps> = ({ label, onChange, defaultValue }) => {
    const [value, setValue] = useState<number>(defaultValue ? defaultValue * 100 : 0);
    const [isEnabled, setIsEnabled] = useState(defaultValue ? true : false);

    useEffect(() => {
        if (isEnabled) {
            onChange(value)
        } else {
            onChange(undefined)
        }
    }, [isEnabled]);

    return (
        <div className="flex gap-2 items-end">
            <Switch aria-label={label} isSelected={isEnabled} onChange={setIsEnabled}>
                <Switch.Control>
                    <Switch.Thumb />
                </Switch.Control>
            </Switch>
                <Slider
                    className="w-full max-w-sm"
                    isDisabled={!isEnabled}
                    value={value}
                    onChange={(v) => {!Array.isArray(v) && setValue(v)}}
                    onChangeEnd={() => {onChange(value)}}>
                <Label>{label}</Label>
                <Slider.Output />
                <Slider.Track>
                    <Slider.Fill />
                    <Slider.Thumb />
                </Slider.Track>
            </Slider>
        </div>
    )
}

interface BPMSliderProps {
    onChange: (lowerValue: number|undefined, upperValue: number|undefined) => void;
    defaultLowerValue: number | undefined;
    defaultUpperValue: number | undefined;
}


const BPMSlider: React.FC<BPMSliderProps> = ({ onChange, defaultUpperValue, defaultLowerValue }) => {
    const [upperValue, setUpperValue] = useState<number>(defaultUpperValue ? defaultUpperValue : 125);
    const [lowerValue, setLowerValue] = useState<number>(defaultLowerValue ? defaultLowerValue : 115);
    const [isEnabled, setIsEnabled] = useState(defaultUpperValue ? true : false);

    useEffect(() => {
        if (isEnabled) {
            onChange(lowerValue, upperValue)
        } else {
            onChange(undefined, undefined)
        }
    }, [isEnabled]);

    const setValue = (v: number[]) => {
        setLowerValue(v[0])
        setUpperValue(v[1])
    }

    return (
        <div className="flex gap-2 items-end">
            <Switch aria-label="bpm" isSelected={isEnabled} onChange={setIsEnabled}>
                <Switch.Control>
                    <Switch.Thumb />
                </Switch.Control>
            </Switch>

            <Slider
                className="w-full max-w-sm"
                maxValue={250}
                minValue={50}
                step={5}
                isDisabled={!isEnabled}
                value={[lowerValue, upperValue]}
                onChange={(v) => {Array.isArray(v) && setValue(v)}}
                onChangeEnd={() => {onChange(lowerValue, upperValue)}}
            >
            <Label>BPM Range</Label>
            <Slider.Output />
            <Slider.Track>
                {({state}) => (
                <>
                    <Slider.Fill />
                    {state.values.map((_, i) => (
                    <Slider.Thumb key={i} index={i} />
                    ))}
                </>
                )}
            </Slider.Track>
            </Slider>
        </div>
    )
}


interface TracksFilterPaneProps {
    onApply: (filter: types.TrackFilter) => void;
    initFilter: types.TrackFilter;
    artists: types.Artist[];
    dynamicApply?: boolean;
}


export const TracksFilterPane: React.FC<TracksFilterPaneProps> = ({ onApply, initFilter, artists, dynamicApply=false }) => {
    const [bpm, setBpm] = useState<number[]|undefined>(initFilter.bpm ? [initFilter.bpm.lower, initFilter.bpm.upper] : undefined)
    const [moodAggressive, setMoodAggressive] = useState<number | undefined>(initFilter.mood.aggressive)
    const [moodCalm, setMoodCalm] = useState<number | undefined>(initFilter.mood.calm)
    const [moodHappy, setMoodHappy] = useState<number | undefined>(initFilter.mood.happy)
    const [moodSad, setMoodSad] = useState<number | undefined>(initFilter.mood.sad)

    const [artistIds, setArtistIds] = useState<string[]>(initFilter.artists ? initFilter.artists : [])

    const onApplyFilter = () => {
        const filter = new types.TrackFilter
        if (bpm !== undefined) {
            filter.bpm = {lower: bpm[0], upper: bpm[1]}
        }
        filter.mood = {
            aggressive: moodAggressive ? moodAggressive / 100 : undefined,
            calm: moodCalm ? moodCalm / 100 : undefined,
            happy: moodHappy ? moodHappy / 100 : undefined,
            sad: moodSad ? moodSad / 100 : undefined,
        }

        if (artistIds.length === 0) {
            filter.artists = undefined
        } else {
            filter.artists = artistIds
        }

        onApply(filter)
    }

    if (dynamicApply) {
        useEffect(() => {
            onApplyFilter()
        }, [bpm, moodAggressive, moodCalm, moodHappy, moodSad, artistIds]);
    }

    return (
        <>
            <ArtistSelect
                artists={artists}
                name="artists"
                selectedArtistIds={artistIds}
                onChange={(k) => {setArtistIds(k.map((i) => (i.toString())))}}
            />
            <Text type="h6"className="pt-4 pb-2">BPM</Text>
            <BPMSlider
                defaultUpperValue={initFilter.bpm?.upper}
                defaultLowerValue={initFilter.bpm?.lower}
                onChange={(vl, vu) => {
                    if (vu === undefined || vl === undefined) {
                        setBpm(undefined)
                        return
                    }
                    setBpm([vl, vu])
                }}
            />
            <Text type="h6"className="pt-4 pb-2">Mood</Text>
            <div className="flex flex-col gap-4">
                <PercentSlider
                    label="Aggressive"
                    defaultValue={moodAggressive}
                    onChange={setMoodAggressive}
                />
                <PercentSlider
                    label="Calm"
                    defaultValue={moodCalm}
                    onChange={setMoodCalm}
                />
                <PercentSlider
                    label="Happy"
                    defaultValue={moodHappy}
                    onChange={setMoodHappy}
                />
                <PercentSlider
                    label="Sad"
                    defaultValue={moodSad}
                    onChange={setMoodSad}
                />
            </div>
            { !dynamicApply &&
                <div className="flex gap-2 justify-end pt-8">
                    <Button slot="close" variant="secondary">
                        Cancel
                    </Button>
                    <Button slot="close" onClick={onApplyFilter}>Apply Filter</Button>
                </div>
            }
        </>
    )
}
