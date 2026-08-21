import { mqMobile } from "@/config/config";
import { cn } from "@heroui/styles";
import { useMediaQuery } from "react-responsive";

export default function DefaultLayout({
  scroll=true,
  paddingTop=true,
  children,
}: {
  scroll?: boolean;
  paddingTop?: boolean;
  children: React.ReactNode;
}) {
  const isMobile = useMediaQuery({ query: mqMobile });

  return (
      <main
      className={cn(
        "container flex flex-col flex-grow p-2 mx-auto sm:p-6 sm:pb-6 pb-35",
        scroll ? "overflow-y-scroll max-h-max" : isMobile ? "overflow-y-auto max-h-[calc(100dvh)]" : "overflow-y-auto max-h-[calc(100dvh-100px)]",
        paddingTop ? "" : "pt-0 sm:pt-0",
      )}>
      {children}
    </main>
  );
}
