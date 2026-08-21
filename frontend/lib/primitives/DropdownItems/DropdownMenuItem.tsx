import {
  cn,
  Description,
  Dropdown,
  DropdownItemProps,
  Key,
  Label,
} from "@heroui/react";
import { LucideIcon } from "lucide-react";

interface DropdownMenuItemProps {
  id: Key;
  icon: LucideIcon;
  label: string;
  description: string;
  variant?: DropdownItemProps["variant"];
}

export const DropdonwMenuItem: React.FC<DropdownMenuItemProps> = ({
  id,
  icon: Icon,
  label,
  description,
  variant = "default",
}) => {
  const iconClasses = cn(
    "text-xl text-foreground pointer-events-none shrink-0",
    variant == "danger" ? "text-danger" : "",
  );

  return (
    <Dropdown.Item id={id} textValue={label} variant={variant}>
      <div className="flex justify-center items-start pt-px h-8">
        <Icon className={iconClasses} />
      </div>
      <div className="flex flex-col">
        <Label>{label}</Label>
        <Description className={variant == "danger" ? "text-danger" : ""}>
          {description}
        </Description>
      </div>
    </Dropdown.Item>
  );
};
