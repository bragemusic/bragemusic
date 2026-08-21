import { Button, ButtonProps } from "@heroui/react";
import { useCallback, useRef, useState } from "react";
import { twMerge } from "tailwind-merge";

export default function FileUploadButton({
  accept,
  onUpload,
  acceptProps = { variant: "primary" },
  rejectProps = { variant: "danger" },
  multiple = false,
  classNames,
  className,
  ...props
}: ButtonProps & {
  classNames?: { wrapper?: string; button?: string };
  onUpload?: (files: File[]) => void;
  acceptProps?: ButtonProps;
  rejectProps?: ButtonProps;
  accept?: string;
  multiple?: boolean;
}) {
  const inputRef = useRef<HTMLInputElement | null>(null);
  const [acceptance] = useState<null | "ACCEPT" | "REJECT">(null);
  const onFileChosen = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      onUpload?.(Array.from(e.target.files!));
    },
    [onUpload],
  );
  const onButtonPress = useCallback(() => {
    inputRef.current?.click();
  }, [inputRef]);

  return (
    <form className={classNames?.wrapper}>
      <label htmlFor="_upload">
        <Button
          {...props}
          {...(acceptance === "ACCEPT"
            ? acceptProps
            : acceptance === "REJECT"
              ? rejectProps
              : {})}
          className={twMerge(classNames?.button)}
          onPress={onButtonPress}
        />
      </label>
      <input
        ref={inputRef}
        accept={accept}
        className="hidden"
        multiple={multiple}
        name="_upload"
        role="presentation"
        type="file"
        onChange={onFileChosen}
      />
    </form>
  );
}
