import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import { cn } from "@heroui/react";
import { BlurredBackdrop } from "@/components/BlurredBackdrop";

type Backdrop = {
  src: string;
  saturate?: number;
  height?: number;
};

type BackgroundContextValue = {
  backdrop: Backdrop | null;
  setBackdrop: (backdrop: Backdrop | null) => void;
  clearBackdrop: () => void;
};

const BackgroundContext = createContext<BackgroundContextValue | null>(null);

type BackgroundProviderProps = {
  children: ReactNode;
  className?: string;
};

export function BackgroundProvider({ children, className }: BackgroundProviderProps) {
  const [backdrop, setBackdropState] = useState<Backdrop | null>(null);

  const setBackdrop = useCallback((next: Backdrop | null) => {
    setBackdropState(next);
  }, []);

  const clearBackdrop = useCallback(() => {
    setBackdropState(null);
  }, []);

  const value = useMemo(
    () => ({ backdrop, setBackdrop, clearBackdrop }),
    [backdrop, setBackdrop, clearBackdrop]
  );

  return (
    <BackgroundContext.Provider value={value}>
      {/* `relative isolate` keeps the backdrop's -z-10 above this div's bg-background */}
      <div className={cn("relative isolate min-h-full w-full bg-background", className)}>
        {backdrop && (
          <BlurredBackdrop
            src={backdrop.src}
            saturate={backdrop.saturate}
            height={backdrop.height}
          />
        )}
        {children}
      </div>
    </BackgroundContext.Provider>
  );
}

/** Imperative access: set or clear the backdrop whenever you like. */
export function useBackground() {
  const ctx = useContext(BackgroundContext);
  if (!ctx) {
    throw new Error("useBackground must be used inside <BackgroundProvider>");
  }
  return ctx;
}

/**
 * Declarative shortcut: "while this component is mounted, show this backdrop".
 * Passing no src shows the plain bg-background. Clears on unmount.
 */
export function useBackdrop(
  src: string | undefined,
  options: { saturate?: number; height?: number } = {}
) {
  const { setBackdrop, clearBackdrop } = useBackground();
  const { saturate, height } = options;

  useEffect(() => {
    if (!src) {
      clearBackdrop();
      return;
    }
    setBackdrop({ src, saturate, height });
    return () => clearBackdrop();
  }, [src, saturate, height, setBackdrop, clearBackdrop]);
}
