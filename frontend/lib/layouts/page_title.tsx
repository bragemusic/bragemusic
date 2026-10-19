import { cn } from "@heroui/styles";
import DefaultLayout from "./default";

type PageTitleLayoutProps = {
  title?: string;
  subtitle?: string;
  headerContent?: React.ReactNode;
  scroll?: boolean;
  children: React.ReactNode;
};

export default function PageTitleLayout({
  title,
  subtitle,
  headerContent,
  scroll=true,
  children,
}: PageTitleLayoutProps) {
  return (
    <div className="overflow-scroll h-full sm:pt-6 bg-surface-secondary border-l-1 border-border">
    <DefaultLayout scroll={scroll} className="">
        <div className="flex flex-col justify-between px-1 mt-7 w-full sm:px-0 sm:mt-0 xl:flex-row xl:items-end">
          <div>
            {title &&
              <h2 className="text-4xl font-black sm:text-6xl font-display text-foreground">{title}</h2>
            }
            {subtitle &&
              <p className="mt-1 text-sm sm:text-base text-foreground/70">{subtitle}</p>
            }
          </div>
          {headerContent && (
            <div className="flex gap-4 mt-4 xl:mt-0">{headerContent}</div>
          )}
        </div>
      <div className={cn("flex overflow-hidden flex-1 justify-between w-full min-h-0 select-none", title != undefined ? "mt-12" : "")}>{children}</div>
    </DefaultLayout>
    </div>
  );
}
