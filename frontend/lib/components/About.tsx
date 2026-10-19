import { Modal } from "@heroui/react";
import { Info } from "lucide-react";
import { Fragment, useState } from "react";

interface AboutProps {
    text?: string;
    tableData?: Record<string, string>;
}

export const About: React.FC<AboutProps> = ({
    text,
    tableData
}) => {
    const [isOpen, setIsOpen] = useState(false)
    return (
        <div className="">
            <h3 className="mb-2 text-2xl font-bold font-display">About</h3>
            {text &&
                <>
                <p className="line-clamp-5 text-foreground/70">{text}</p>
                <p
                    className="mt-2 mb-6 text-sm font-bold cursor-pointer hover:underline text-accent-deep"
                    onClick={() => {
                        setIsOpen(true)
                    }}
                >
                    Read more
                </p>
                </>
            }
            {tableData && (
                <dl className="grid gap-y-2.5 gap-x-5 text-sm grid-cols-[auto_1fr]">
                {Object.entries(tableData).map(([label, value]) => (
                    <Fragment key={label}>
                        <dt className="text-muted">{label}</dt>
                        <dd>{value}</dd>
                    </Fragment>
                ))}
                </dl>
            )}
            <Modal.Backdrop isOpen={isOpen} onOpenChange={setIsOpen}>
            <Modal.Container size="lg">
                <Modal.Dialog className="">
                <Modal.CloseTrigger />
                <Modal.Header>
                    <Modal.Icon className="bg-accent-soft text-accent-soft-foreground">
                    <Info className="size-5" />
                    </Modal.Icon>
                    <Modal.Heading>About</Modal.Heading>
                </Modal.Header>
                <Modal.Body>
                    <p>
                        {text}
                    </p>
                </Modal.Body>
                </Modal.Dialog>
            </Modal.Container>
            </Modal.Backdrop>
        </div>
    )
}
