import { Button, Card } from "@heroui/react";
import { BadgeCheck, Music2, Palette } from "lucide-react";
import { useTheme } from "next-themes";

interface ThemeCardProps {
  theme_id: string;
  theme_name: string;
}

export const ThemeCard: React.FC<ThemeCardProps> = ({
  theme_id,
  theme_name,
}) => {
  const { theme, setTheme } = useTheme();

  return (
    <Card
      className={`relative gap-0 p-0 w-[250px] h-[250px] border-1 border-border`}
    >
      <div
        className="flex flex-col justify-between p-4 w-full h-full"
        data-theme={theme_id}
      >
        <div className="flex flex-col flex-grow rounded-md bg-background border-1 border-border">
          <div className="flex justify-between items-center py-1 px-2">
            <div className="flex gap-1">
              <div className="w-4 h-4 rounded-full bg-danger" />
              <div className="w-4 h-4 rounded-full bg-warning" />
              <div className="w-4 h-4 rounded-full bg-success" />
            </div>
            <p className="text-sm font-semibold text-foreground">
              {theme_name}
            </p>
          </div>
          <div className="flex flex-grow w-full border-t-1 border-border">
            <div className="flex flex-col gap-1 p-2 w-1/6 bg-surface-secondary border-r-1 border-border">
              <div className="w-full h-1 bg-foreground" />
              <div className="w-full h-1 bg-foreground" />
              <div className="w-full h-1 bg-foreground" />
            </div>
            <div className="flex flex-col flex-grow bg-background">
              <div className="flex flex-grow gap-4 justify-center items-center">
                <div className="flex justify-center items-center rounded-md w-15 h-15 bg-default text-foreground">
                  <Music2 className="stroke-foreground" />
                </div>
                <div className="flex justify-center items-center rounded-md w-15 h-15 bg-overlay text-foreground">
                  <Palette className="stroke-accent" />
                </div>
              </div>
              <div className="flex gap-1 justify-center items-center w-full h-1/6 bg-surface border-t-1 border-border">
                <div className="w-1 h-1 rounded-full bg-default-foreground" />
                <div className="w-2 h-2 rounded-full bg-accent" />
                <div className="w-1 h-1 rounded-full bg-default-foreground" />
              </div>
            </div>
          </div>
        </div>
        <div className="flex justify-between items-center pt-4">
          <div className="flex gap-1 items-center">
            {theme == theme_id && (
              <>
                <BadgeCheck className="fill-success stroke-surface" />
                <p className="text-lg font-semibold text-accent">Current</p>
              </>
            )}
          </div>
          <Button onClick={() => setTheme(theme_id)}>Select</Button>
        </div>
      </div>
    </Card>
  );
};
