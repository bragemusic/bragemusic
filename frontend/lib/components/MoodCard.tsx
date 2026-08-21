import { Card, cn } from "@heroui/react";
import {
    Angry,
  Frown,
  Laugh,
  LucideIcon,
  Smile,
} from "lucide-react";


interface MoodCardProps {
  mood: "angry"|"calm"|"happy"|"sad";
  frac: number;
}

export const MoodCard: React.FC<MoodCardProps> = ({ mood, frac }) => {

  const moodIcon = (): LucideIcon => {
    switch (mood) {
      case "angry":
        return Angry;
      case "calm":
        return Smile;
      case "happy":
        return Laugh;
      case "sad":
        return Frown;
      default:
        return Smile;
    }
  };

  const Icon = moodIcon();

  return (
    <Card className="flex flex-col gap-0 p-0 cursor-default select-none border-1 border-border w-[150px] h-[150px]">
      <div className={cn(
        "flex flex-col flex-grow gap-2 justify-center items-center p-2",
        mood == "angry" ? "bg-danger text-danger-foreground" : "",
        mood == "calm" ? "bg-accent text-accent-foreground" : "",
        mood == "happy" ? "bg-success text-success-foreground" : "",
        mood == "sad" ? "bg-warning text-warning-foreground" : "",
        )}
      >
        <Icon size={70} />
        <p className="text-xl font-semi-bold">{Math.round(frac * 100)}%</p>
      </div>
    </Card>
  );
};
