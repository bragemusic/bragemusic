import {
  Input as BaseInput,
  InputProps as BaseInputProps,
  DateField as BaseDateInput,
  Label,
  TextField,
  Description,
  TextArea as BaseTextarea,
  NumberField,
} from "@heroui/react";
import { parseDate } from "@internationalized/date";

type InputProps = {
  name: string;
  label: string;
  defaultValue?: string;
  value?: string;
  onChange?: (value: string | undefined) => void;
  required?: boolean;
  color?: BaseInputProps["color"];
  type?: string;
  description?: string;
  className?: string;
  isDisabled?: boolean;
};

type NumberInputProps = {
  name: string;
  label: string;
  description?: string;
  defaultValue?: number;
  required?: boolean;
};

export function Input({
  name,
  label,
  defaultValue,
  onChange,
  value,
  description,
  required = false,
  color = "default",
  type = "text",
  className = "",
  isDisabled = false,
}: InputProps) {
  return (
    <TextField isDisabled={isDisabled} className={className} isRequired={required} defaultValue={defaultValue} name={name} type={type} >
      <Label color={color}>{label}</Label>
      <BaseInput
        color={color}
          required={required}
          variant="secondary"
          onChange={(e) => {
            const val = e.target.value;
            onChange?.(val === "" ? undefined : val);
          }}
          value={value}
      />
      {description && <Description color={color}>{description}</Description>}
    </TextField>
  );
}

export function NumberInput({
  name,
  label,
  description,
  defaultValue,
  required = false,
}: NumberInputProps) {
  return (
    <NumberField
      defaultValue={defaultValue}
      name={name}
      formatOptions={{ useGrouping: false }}
      isRequired={required}
      variant="secondary"
    >
      <Label>{label}</Label>
      <NumberField.Group>
        <NumberField.DecrementButton />
        <NumberField.Input className="w-[120px]" />
        <NumberField.IncrementButton />
      </NumberField.Group>
      {description && <Description>{description}</Description>}
    </NumberField>
  );
}

export function DateInput({ name, label, defaultValue }: InputProps) {
  return (
    <BaseDateInput
      defaultValue={
       defaultValue
         ? parseDate(defaultValue.split("T")[0])
         : null
      }
      name={name}
    >
      <Label>{label}</Label>
      <BaseDateInput.Group variant="secondary">
        <BaseDateInput.Input>
          {(segment) => <BaseDateInput.Segment segment={segment} />}
        </BaseDateInput.Input>
      </BaseDateInput.Group>
    </BaseDateInput>
  );
}

export function Textarea({
  name,
  label,
  defaultValue,
  required = false,
}: InputProps) {
  return (
    <div className="flex flex-col gap-2">
      <Label htmlFor={name}>{label}</Label>
      <BaseTextarea
        aria-label={label}
        color="default"
        defaultValue={defaultValue}
        name={name}
        required={required}
        variant="secondary"
      />
    </div>
  );
}
