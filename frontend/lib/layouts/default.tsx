import { mqMobile } from "@/config/config";
import { cn } from "@heroui/styles";
import { useMediaQuery } from "react-responsive";

export default function DefaultLayout({
  scroll=true,
  paddingTop=true,
  children,
  className="",
}: {
  scroll?: boolean;
  paddingTop?: boolean;
  children: React.ReactNode;
  className?: string;
}) {
  const isMobile = useMediaQuery({ query: mqMobile });

  return (
      <main
      className={cn(
        "container flex flex-col flex-grow px-2 py-2 mx-auto sm:px-12 sm:py-4 sm:pb-40 pb-35",
        scroll ? "overflow-y-scroll max-h-max" : isMobile ? "overflow-y-auto max-h-[calc(100dvh)]" : "overflow-y-auto max-h-[calc(100dvh-100px)]",
        paddingTop ? "" : "pt-0 sm:pt-0",
        className,
      )}>
      {children}
    </main>
  );
}
