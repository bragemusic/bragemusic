import {
  Button as BaseButton,
  ButtonProps,
  PressEvent,
  Popover,
  Tooltip,
} from "@heroui/react";
import { type LucideIcon } from "lucide-react";
import { useState } from "react";

import FileUploadButton from "@/components/FileUploadButton";

export type BannerButtonProps = {
  icon: LucideIcon;
  tooltip?: string;
  onClick?: (e: PressEvent) => void;
  onUpload?: (files: File[]) => void;
  variant?: ButtonProps["variant"];
  type?: "normal" | "fileupload";
  isDisabled?: boolean;
  confirm?: boolean;
};

export function BannerButton({
  icon: Icon,
  tooltip,
  onClick,
  onUpload,
  variant = "primary",
  type = "normal",
  isDisabled = false,
  confirm = false,
}: BannerButtonProps) {
  if (type == "normal") {
    return (
      <ActionButton
        confirm={confirm}
        icon={Icon}
        isDisabled={isDisabled}
        size="sm"
        tooltip={tooltip}
        variant={variant}
        onClick={onClick}
      />
    );
  }

  return (
    <FileUploadButton
      isIconOnly
      accept="image/*"
      isDisabled={isDisabled}
      multiple={false}
      size="sm"
      variant={variant}
      onUpload={onUpload}
    >
      <Icon size={20} />
    </FileUploadButton>
  );
}

export type ActionButtonProps = {
  icon: LucideIcon;
  tooltip?: string;
  onClick?: (e: PressEvent) => void;
  size?: ButtonProps["size"];
  variant?: ButtonProps["variant"];
  confirm?: boolean;
  isDisabled?: boolean;
  buttonClassName?: string;
};

export function ActionButton({
  icon: Icon,
  tooltip,
  onClick,
  size = "md",
  variant = "primary",
  confirm = false,
  isDisabled = false,
  buttonClassName,
}: ActionButtonProps) {
  const [popoverOpen, setPopoverOpen] = useState(false);

  const iconSize = (size: ButtonProps["size"]): number => {
    if (size === "sm") {
      return 20;
    } else {
      return 30;
    }
  };

  const onPressAction = (e: PressEvent) => {
    onClick?.(e);
    setPopoverOpen(false);
  };

  const button = (
    <BaseButton
      isIconOnly
      className={buttonClassName}
      isDisabled={isDisabled}
      size={size}
      variant={variant}
      onPressEnd={(e) => {
        if (!confirm) {
          onClick?.(e);
        } else {
          setPopoverOpen(true);
        }
      }}
    >
      <Icon size={iconSize(size)} />
    </BaseButton>
  );

  const popoverContent = (
    <Popover.Content className="w-[240px]">
      <Popover.Arrow />
      <Popover.Dialog>
        <Popover.Heading />
        <div className="py-2 px-1 w-full">
          <p className="font-bold text-small text-foreground">Are you sure?</p>
          <div className="flex gap-2 mt-2 w-full">
            <BaseButton
              className="w-full"
              variant="ghost"
              onPressEnd={onPressAction}
            >
              <Icon />
              Yes
            </BaseButton>
          </div>
        </div>
      </Popover.Dialog>
    </Popover.Content>
  );

  return (
    <Tooltip closeDelay={100} delay={300} isDisabled={!tooltip}>
      <span className="inline-flex">
        {confirm ? (
          <Popover isOpen={popoverOpen} onOpenChange={setPopoverOpen}>
            <Popover.Trigger>{button}</Popover.Trigger>
            {popoverContent}
          </Popover>
        ) : (
          button
        )}
      </span>
      <Tooltip.Content>
        <p>{tooltip}</p>
      </Tooltip.Content>
    </Tooltip>
  );
}
