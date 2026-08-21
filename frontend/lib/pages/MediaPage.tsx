import PageTitleLayout from "@/layouts/page_title";
import { Key, ListBox } from "@heroui/react";
import { ChevronRight, Disc, Heart, ListMusic, LucideIcon, Music4, UsersIcon } from "lucide-react";
import { useNavigate } from "react-router-dom";

interface IconWrapperProps {
  icon: LucideIcon
}

const IconWrapper = ({ icon: Icon }: IconWrapperProps) => {
  return (
    <div
      className="flex justify-center items-center p-1 w-12 h-12 rounded-sm border-1 border-border bg-surface-secondary text-foreground"
    >
      <Icon/>
    </div>
  )
}

export default function ImportPage() {
  const navigate = useNavigate()

  const items = [
    { key: "Artists", href: "/artists", icon: UsersIcon },
    { key: "Albums", href: "/albums", icon: Disc },
    { key: "Tracks", href: "/tracks", icon: Music4 },
    { key: "Playlists", href: "/playlists", icon: ListMusic },
    { key: "Liked Tracks", href: "/liked-tracks", icon: Heart },
  ];

  const onAction = (key: Key) => {
    items.forEach((item) => {
      if (item.key == key) {
        navigate(item.href);
        return;
      }
    });
  }
  return (
    <PageTitleLayout title="Media">
      <ListBox
        aria-label="mediamenu"
        className="overflow-visible gap-2 p-0 pr-2 pb-5 pl-4 mt-5 w-full bg-none"
        onAction={onAction}
      >
          {items.map((item) => (
            <ListBox.Item className="justify-between w-full" id={item.key} key={item.key}>
              <div className="flex gap-3 items-center">
                <IconWrapper icon={item.icon} />
                <p className="text-xl font-semibold text-foreground">{item.key}</p>
              </div>
              <ChevronRight className="stroke-foreground"/>
            </ListBox.Item>
          ))}
      </ListBox>
    </PageTitleLayout>
  )
}
