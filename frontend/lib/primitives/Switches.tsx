import { Switch as BaseSwitch, Description, Label } from "@heroui/react";

type SwitchProps = {
  name: string;
  label: string;
  defaultSelected?: boolean;
  required?: boolean;
  className?: string;
  description?: string;
};

export function Switch({
  name,
  label,
  defaultSelected,
  required = false,
  className,
  description,
}: SwitchProps) {
  return (
    <BaseSwitch
      className={className}
      defaultSelected={defaultSelected}
      name={name}
      size="md"
    >
      <BaseSwitch.Control aria-required={required}>
        <BaseSwitch.Thumb />
      </BaseSwitch.Control>
      <BaseSwitch.Content>
        <Label className="text-sm">{label}</Label>
        {description && <Description>{description}</Description>}
      </BaseSwitch.Content>
    </BaseSwitch>
  );
}
