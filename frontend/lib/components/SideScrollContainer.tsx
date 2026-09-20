import { ScrollShadow } from "@heroui/react";

interface SideScrollContainerProps {
  children: React.ReactNode;
}

export const SideScrollContainer: React.FC<SideScrollContainerProps> = ({
    children,
}) => {
    return (
            <ScrollShadow
                className="py-4 max-w-min scrollbar-thin"
                orientation="horizontal"
            >
                <div className="flex flex-row gap-6 w-max">
                    {children}
                </div>
            </ScrollShadow>
    );
};
