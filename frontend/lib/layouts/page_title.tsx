import { cn } from "@heroui/styles";
import DefaultLayout from "./default";

type PageTitleLayoutProps = {
  title?: string;
  headerContent?: React.ReactNode;
  scroll?: boolean;
  children: React.ReactNode;
};

export default function PageTitleLayout({
  title,
  headerContent,
  scroll=true,
  children,
}: PageTitleLayoutProps) {
  return (
    <DefaultLayout scroll={scroll}>
        <div className="flex justify-between items-center">
          {title &&
            <h2 className="text-2xl font-semibold sm:text-4xl text-foreground">{title}</h2>
          }
          {headerContent && (
            <div className="flex gap-4 items-center">{headerContent}</div>
          )}
        </div>
      <div className={cn("flex overflow-hidden flex-1 justify-between w-full min-h-0select-none", title != undefined ? "mt-12" : "")}>{children}</div>
    </DefaultLayout>
  );
}
