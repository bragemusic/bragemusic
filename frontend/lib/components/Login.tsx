import { useApi } from "@/api/ApiContext";
import { useSession } from "@/session/SessionContext";
import Logo from "../icons/logo.svg?react";
import { Button, Checkbox, Form, Label, Surface, Text } from "@heroui/react"
import { useEffect, useState } from "react";
import { Input } from "@/primitives/Inputs";
import { CircleAlert } from "lucide-react";
import { Event, Message } from "@/types/events";

interface LoginProps {}

export const Login: React.FC<LoginProps> = ({}) => {
    const [rememberMe, setRememberMe] = useState(true);
    const [errorMessage, setErrorMessage] = useState<Message | undefined>(undefined);

    const api = useApi();
    const session = useSession();

    useEffect(() => {
        const errHandler = (msg: Message) => {
            setErrorMessage(msg)
        };

        const unsub1 = api.eventSubscribe(Event.MsgErr, errHandler);

        return () => {
        unsub1()
        };
    }, [api]);

    return (
    <div className="flex flex-col justify-center items-center w-full h-full bg-background text-foreground fill-foreground">
        <Surface className="flex flex-col gap-3 justify-between items-center p-6 w-full h-full rounded-3xl shadow-md sm:w-2/3 sm:h-1/2 xl:w-1/3 min-w-[320px] min-h-[640px]" variant="default">
            <div className="flex flex-col items-center w-full">
                <Logo className="mb-4 fill-accent" height={140} width={140} />
                <Text className="mb-2" type="h2">Brage Music</Text>
                <Text type="body">{session.serverInfo.name}</Text>
                {errorMessage &&
                    <div className="flex flex-col gap-1 py-2 px-4 mt-4 w-full rounded-md bg-danger-soft text-danger border-danger border-1">
                        <div className="flex gap-2 font-semibold">
                            <CircleAlert/>
                            <p>{errorMessage.title}</p>
                        </div>
                    <p>{errorMessage.message}</p>
                    </div>
            }
            </div>

            <Form
                className="flex flex-col gap-4 w-full"
                onReset={() => {}}
                onSubmit={(e) => {
                    e.preventDefault();
                    const formData = new FormData(e.currentTarget);

                    const email = formData.get("email") as string;
                    const password = formData.get("password") as string;

                    api.loginServerUser(email, password, rememberMe);
                }}
            >
                <Input
                    required
                    description="Enter the email tied to your account"
                    label="Email"
                    name="email"
                    type="email"
                />
                <Input
                    required
                    description="Enter your password"
                    label="Password"
                    name="password"
                    type="password"
                />
                <Checkbox
                    className="mt-2"
                    isSelected={rememberMe}
                    name="long_lived_token"
                    onChange={setRememberMe}
                >
                    <Checkbox.Control>
                        <Checkbox.Indicator />
                    </Checkbox.Control>
                    <Checkbox.Content>
                        <Label>Keep me logged in</Label>
                    </Checkbox.Content>
                </Checkbox>
                <div className="flex gap-2 justify-end mt-2 mb-4 w-full">
                    <Button type="submit" variant="primary">
                    Log in
                    </Button>
                </div>
            </Form>
        </Surface>
    </div>
    )
}
