import { ScrollShadow } from "@heroui/react";

interface SideScrollContainerProps {
  children: React.ReactNode;
}

export const SideScrollContainer: React.FC<SideScrollContainerProps> = ({
    children,
}) => {
    return (
            <ScrollShadow
                className="py-4 max-w-min scrollbar-none"
                orientation="horizontal"
            >
                <div className="flex flex-row gap-4 w-max">
                    {children}
                </div>
            </ScrollShadow>
    );
};
