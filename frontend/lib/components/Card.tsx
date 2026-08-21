import { Card as BaseCard } from "@heroui/react";
import { LucideIcon } from "lucide-react";
import { useMediaQuery } from 'react-responsive'
import { mqMobile }from "@/config/config";

import { Image } from "./Image";
import { ActionButton, ActionButtonProps } from "@/primitives/Buttons";

interface CardProps {
  title: string;
  imgSrc?: string;
  fallbackIcon: LucideIcon;
  description?: string;
  onClick?: () => void;
  radius?:
    | "xs"
    | "sm"
    | "md"
    | "lg"
    | "xl"
    | "2xl"
    | "3xl"
    | "4xl"
    | "full"
    | "none";
  overlayIcon?: LucideIcon;
  actionButtons?: ActionButtonProps[];
}

export const Card: React.FC<CardProps> = ({
  title,
  description,
  imgSrc,
  fallbackIcon,
  onClick,
  radius = "3xl",
  overlayIcon: OverlayIcon,
  actionButtons,
}) => {
  const isMobile = useMediaQuery({ query: mqMobile })

  if (isMobile) {
      return (
        <BaseCard
          className={`relative gap-0 p-0 ${onClick ? "cursor-pointer" : ""} w-full border-1 border-border rounded-sm`}
          onClick={onClick}
        >
          <div className="flex relative gap-4 w-full">
            {OverlayIcon &&
              <div className="absolute top-0 left-0 p-1 mt-1 ml-1 rounded-full shadow-md bg-accent text-accent-foreground"><OverlayIcon size={10}/></div>
            }
            {actionButtons &&
              <div className="flex gap-2 justify-center pb-2 w-full">
                {actionButtons.map((d) => (
                  <ActionButton
                    confirm={d.confirm}
                    icon={d.icon}
                    size={d.size}
                    tooltip={d.tooltip}
                    variant={d.variant}
                    onClick={d.onClick}
                  />
                ))}
              </div>
            }
            <Image
              fallbackIcon={fallbackIcon}
              height={80}
              width={80}
              radius="none"
              src={imgSrc}
            />
            <div className="flex overflow-hidden flex-col pt-3 h-full">
              <div className="w-full text-base font-medium sm:text-lg text-foreground truncate">
                {title}
              </div>
              {description && (
                <div className="pt-1 w-full text-xs font-medium sm:text-sm text-foreground/50">
                  {description}
                </div>
              )}
            </div>
          </div>
        </BaseCard>
    );
  } else {
    return (
    <div
      onClick={onClick}
      className={`flex flex-col items-center pb-4 ${onClick ? "cursor-pointer" : ""}`}
      >
      <div className="relative">
        <Image
          height={180}
          width={180}
          fallbackIcon={fallbackIcon}
          src={imgSrc}
          radius={radius}
          className="shadow-md border-1 border-border"
        />
        {OverlayIcon &&
          <div className="absolute top-0 left-0 p-2 mt-2 ml-2 rounded-full shadow-md bg-accent text-accent-foreground"><OverlayIcon /></div>
        }
        {actionButtons &&
          <div className="flex absolute bottom-0 left-0 gap-2 justify-center pb-2 w-full">
            {actionButtons.map((d) => (
              <ActionButton
                confirm={d.confirm}
                icon={d.icon}
                size={d.size}
                tooltip={d.tooltip}
                variant={d.variant}
                onClick={d.onClick}
              />
            ))}
          </div>
        }
      </div>
      <p className="pt-2 w-full font-medium text-center truncate text-foreground">{title}</p>
      <p className="w-full text-center truncate text-foreground/80">
        {description}
      </p>
    </div>
    )
  }
};
