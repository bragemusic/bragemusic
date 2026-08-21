import { useSession } from "@/session/SessionContext";
import { Button, Table, Tabs, Text, toast, useOverlayState } from "@heroui/react";
import { useEffect, useState } from "react";
import { Check, ImageUp, KeyRound, LucideIcon, Palette, Plus, Smartphone, Trash, Undo2, User } from "lucide-react";
import { Input } from "@/primitives/Inputs";
import { userAvatarLink } from "@/util/images";
import { useApi } from "@/api/ApiContext";
import EditModal from "@/components/EditModal";
import { types } from "@/types/core";
import { useDeviceStore } from "@/store/deviceStore";
import CardLayout from "@/layouts/CardLayout";
import { DeviceRow } from "@/components/DeviceRow";
import { ThemeCard } from "@/components/ThemeCard";
import DefaultLayout from "@/layouts/default";
import { dateHR, parseDate } from "@/util/functions";
import { ActionButton } from "@/primitives/Buttons";
import { responses } from "@/types/server";
import { useMediaQuery } from "react-responsive";
import { mqMobile } from "@/config/config";


interface TokensTabProps {
}

const TokensTab: React.FC<TokensTabProps> = ({
}) => {
    const [tokens, setTokens] = useState<types.TokenLimited[]>([])
    const [apiToken, setApiToken] = useState<responses.CreateApiToken|undefined>(undefined)

    const api = useApi();
    const session = useSession();
    const addAPITokenState = useOverlayState();
    const copyAPITokenState = useOverlayState();

    const loadTokens = () => {
        api.listUserTokens().then(setTokens);
    }

    useEffect(() => {
        loadTokens()
    }, []);

    useEffect(() => {
        if (apiToken === undefined) {
            copyAPITokenState.setOpen(false)
        } else {
            copyAPITokenState.setOpen(true)
        }
    }, [apiToken]);

    const tokenType = (tt: string): string => {
        switch (tt) {
            case "frontend_short":
                return "Short Lived"
            case "frontend_long":
                return "Long Lived"
            default:
                return "Unknown"
        }
    }

    const onAddAPIToken = (fd: FormData) => {
        const name = String(fd.get("name") ?? "")

        if (name == "") {
            toast("Cannot create API token", {
                description: "No name provided",
                variant: "danger",
            });
            return
        }

        api.addApiToken(name).then(setApiToken).finally(loadTokens);
    };

    return (
        <div className="flex flex-col">
            <EditModal
                header="API Token Created"
                icon={Check}
                state={copyAPITokenState}
                canSubmit={false}
                onSubmit={() => {}}
            >
                <Text type="body">Save your token now, this is the only time you can see it.</Text>
                <Text type="code">{apiToken ? apiToken.token : ""}</Text>
            </EditModal>
            <EditModal
                header="Add API Token"
                state={addAPITokenState}
                onSubmit={onAddAPIToken}
            >
                <Input required label="Name" name="name" />
            </EditModal>
            <div className="flex gap-6 items-center pb-4">
                <Text type="h5">API Tokens</Text>
                <ActionButton
                    size="sm"
                    icon={Plus}
                    isDisabled={!session.serverFullyAvailable}
                    tooltip="Add new API Token"
                    onClick={() => {
                        addAPITokenState.setOpen(true);
                    }}
                />
            </div>
            <Table>
                <Table.ScrollContainer>
                    <Table.Content aria-label="Team members" className="min-w-[600px]">
                        <Table.Header>
                            <Table.Column isRowHeader>Name</Table.Column>
                            <Table.Column>Scopes</Table.Column>
                            <Table.Column>Last Used</Table.Column>
                            <Table.Column>Expires At</Table.Column>
                            <Table.Column>Created</Table.Column>
                            <Table.Column></Table.Column>
                        </Table.Header>
                        <Table.Body>
                            {tokens
                                .filter((t) => t.type == "api")
                                .map((t) => (
                                        <Table.Row key={t.id} id={t.id}>
                                            <Table.Cell>{t.name}</Table.Cell>
                                            <Table.Cell>{t.scopes}</Table.Cell>
                                            <Table.Cell>{dateHR(parseDate(t.last_used_at))}</Table.Cell>
                                            <Table.Cell>{dateHR(parseDate(t.expires_at))}</Table.Cell>
                                            <Table.Cell>{dateHR(parseDate(t.created_at))}</Table.Cell>
                                            <Table.Cell>
                                                <ActionButton
                                                    size="sm"
                                                    icon={Trash}
                                                    variant="danger-soft"
                                                    isDisabled={!session.serverFullyAvailable}
                                                    tooltip="Remove Token"
                                                    onClick={() => {
                                                        api.deleteToken(t.id).finally(loadTokens);
                                                    }}
                                                    confirm
                                                />
                                            </Table.Cell>
                                        </Table.Row>
                            ))}
                        </Table.Body>
                    </Table.Content>
                </Table.ScrollContainer>
            </Table>
            <Text type="h5" className="pt-8 pb-4">Session Tokens</Text>
            <Table>
                <Table.ScrollContainer>
                    <Table.Content aria-label="Team members" className="min-w-[600px]">
                        <Table.Header>
                            <Table.Column isRowHeader>Type</Table.Column>
                            <Table.Column>Scopes</Table.Column>
                            <Table.Column>Last Used</Table.Column>
                            <Table.Column>Expires At</Table.Column>
                            <Table.Column>Created</Table.Column>
                            <Table.Column></Table.Column>
                        </Table.Header>
                        <Table.Body>
                            {tokens
                                .filter((t) => t.type !== "api")
                                .map((t) => (
                                        <Table.Row key={t.id} id={t.id}>
                                            <Table.Cell>{tokenType(t.type)}</Table.Cell>
                                            <Table.Cell>{t.scopes}</Table.Cell>
                                            <Table.Cell>{dateHR(parseDate(t.last_used_at))}</Table.Cell>
                                            <Table.Cell>{dateHR(parseDate(t.expires_at))}</Table.Cell>
                                            <Table.Cell>{dateHR(parseDate(t.created_at))}</Table.Cell>
                                            <Table.Cell>
                                                <ActionButton
                                                    size="sm"
                                                    variant="danger-soft"
                                                    icon={Trash}
                                                    isDisabled={!session.serverFullyAvailable}
                                                    tooltip="Remove Token"
                                                    onClick={() => {
                                                        api.deleteToken(t.id).finally(loadTokens);
                                                    }}
                                                    confirm
                                                />
                                            </Table.Cell>
                                        </Table.Row>
                            ))}
                        </Table.Body>
                    </Table.Content>
                </Table.ScrollContainer>
            </Table>
        </div>
    )
}

function ThemesTab() {
  const [themes, setThemes] = useState<types.ThemeDescription[]>([]);

  const api = useApi();

  useEffect(() => {
    api.listThemes().then(setThemes);
  }, []);

  return (
    <CardLayout>
      {themes.map((theme) => (
        <ThemeCard theme_id={theme.id} theme_name={theme.name} />
      ))}
    </CardLayout>
  );
}

function DevicesTab() {
  const devices: types.DeviceDetailed[] = Object.values(
    useDeviceStore((s) => s.devices),
  );

    return (
       <div>
            <Table>
                <Table.ScrollContainer>
                    <Table.Content aria-label="Team members" className="min-w-[600px]">
                        <Table.Header>
                            <Table.Column></Table.Column>
                            <Table.Column isRowHeader>Name</Table.Column>
                            <Table.Column>Last Seen</Table.Column>
                            <Table.Column>Last IP Address</Table.Column>
                            <Table.Column>Platform</Table.Column>
                            <Table.Column>Version</Table.Column>
                            <Table.Column>Playback Support</Table.Column>
                            <Table.Column></Table.Column>
                        </Table.Header>
                        <Table.Body>
                            { devices.map((d) => (
                                <DeviceRow device={d}/>
                            ))}
                        </Table.Body>
                    </Table.Content>
                </Table.ScrollContainer>
            </Table>
        </div>
    )
}

interface ProfileTabProps {
}

const ProfileTab: React.FC<ProfileTabProps> = ({
}) => {

    const [initials, setInitials] = useState("??");
    const [userName, setUserName] = useState("");
    const [userEmail, setUserEmail] = useState("");
    const [profileCanSubmit, setProfileCanSubmit] = useState(false);
    const [profileCanEdit, setProfileCanEdit] = useState(false);

    const session = useSession();
    const api = useApi();
    const modalState = useOverlayState();

    useEffect(() => {
        if (session.user && session.user.username !== "") {
            setInitials(
                session.user.username
                .trim()
                .split(/\s+/)
                .map((word) => word[0].toUpperCase())
                .join(""),
            );
            setUserName(session.user.username)
            setUserEmail(session.user.email)
            setProfileCanEdit(!session.user.isSuperAdmin)
        }
    }, [session.user]);


    const setProfileValue = (setter: (v: string) => void, value?: string) => {
        if (value === undefined) {
            value = ""
        }

        setProfileCanSubmit(true)
        setter(value)
    }

    const revertProfile = () => {
        if (session.user) {
            setProfileCanSubmit(false)
            setUserName(session.user.username)
            setUserEmail(session.user.email)
        }
    }

    const onPasswordSubmit = (fd: FormData) => {
        api.updateProfile(undefined, undefined, fd.get("password")?.toString(), fd.get("new_password")?.toString(), fd.get("new_password_confirm")?.toString())
    };

    if (!session.user) {
        return (
            <div>ERROR</div>
        )
    }

    return (
        <div className="flex flex-col items-center pt-4">
            <EditModal
                state={modalState}
                header="Change Password"
                onSubmit={onPasswordSubmit}
            >
                <Input
                    name="password"
                    label="Current Password"
                    type="password"
                    required
                />
                <Input
                    name="new_password"
                    label="New Password"
                    type="password"
                    required
                />
                <Input
                    name="new_password_confirm"
                    label="Confirm New Password"
                    type="password"
                    required
                />
            </EditModal>
            <div className="flex flex-col gap-8 items-center sm:flex-row sm:items-start">
                <div className="flex overflow-hidden relative justify-center items-center w-60 h-60 font-medium rounded-full group text-8xl/60 aspect-square bg-accent-soft text-accent border-1 border-border">
                    <img className="object-cover w-full h-full text-center align-middle" src={userAvatarLink(session.user.id, 320)} alt={initials}/>
                    <Button
                        variant="secondary"
                        className="absolute bottom-8 opacity-80 group-hover:opacity-100"
                        size="sm"
                        onClick={api.uploadUserImage}
                    >
                        <ImageUp/>Upload Image
                    </Button>
                </div>
                <div className="flex flex-col gap-3">
                    <Input isDisabled={!profileCanEdit} name="username" label="Name" value={userName} onChange={(v) => {setProfileValue(setUserName, v)}} description="Your name, it will be visible to logged in users."/>
                    <Input isDisabled={!profileCanEdit} name="email" label="Email" type="email" value={userEmail} onChange={(v) => {setProfileValue(setUserEmail, v)}} description="Your email, this is what you use to log in."/>
                    <div className="flex flex-col gap-2 pt-4 sm:flex-row">
                        <Button
                            isDisabled={!profileCanEdit}
                            variant="secondary"
                            onClick={() => {modalState.setOpen(true)}}
                        >
                            <KeyRound/>Change Password
                        </Button>
                        <Button variant="outline" isDisabled={!profileCanEdit || !profileCanSubmit} onClick={revertProfile}><Undo2/>Revert</Button>
                        <Button
                            isDisabled={!profileCanEdit || !profileCanSubmit}
                            onClick={() => {api.updateProfile(userEmail, userName)}}
                        >
                            <KeyRound/>Update Profile
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default function UserSettingsPage() {
    const isMobile = useMediaQuery({ query: mqMobile });

    const navbarItem = (key: string, Icon: LucideIcon, desc: string) => {
        return (
        <Tabs.Tab id={key}>
            <Icon size={17} />
            {!isMobile &&
                <p className="pl-1">{desc}</p>
            }
            <Tabs.Indicator />
        </Tabs.Tab>
        );
    };

    const tabsPanelClass = "flex flex-col gap-6 py-10 pt-4 w-full justify-center"

    return (
        <DefaultLayout paddingTop={false}>
        <Tabs className="w-full" variant="secondary">
            <Tabs.ListContainer className="self-center pt-6 sm:w-150">
                <Tabs.List aria-label="Options">
                    {navbarItem("profile", User, "Profile")}
                    {navbarItem("tokens", KeyRound, "Tokens")}
                    {navbarItem("devices", Smartphone, "Devices")}
                    {navbarItem("theme", Palette, "Theme")}
                </Tabs.List>
            </Tabs.ListContainer>
            <Tabs.Panel
                className={tabsPanelClass}
                id="profile"
            >
                <ProfileTab/>
            </Tabs.Panel>
            <Tabs.Panel
                className={tabsPanelClass}
                id="tokens"
            >
                <TokensTab/>
            </Tabs.Panel>
            <Tabs.Panel
                className={tabsPanelClass}
                id="devices"
            >
                <DevicesTab/>
            </Tabs.Panel>
            <Tabs.Panel
                className={tabsPanelClass}
                id="theme"
            >
                <ThemesTab/>
            </Tabs.Panel>
        </Tabs>
        </DefaultLayout>
    )
}
