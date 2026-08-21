import PageTitleLayout from "@/layouts/page_title";
import { Key, ListBox } from "@heroui/react";
import { useSession } from "@/session/SessionContext";
import { UserRoles } from "@/types/roles";
import {
  ChevronRight,
  CloudUpload,
  LucideIcon,
  Wrench,
} from "lucide-react";
import { useNavigate } from "react-router-dom";
import { ServerStatus } from "@/components/ServerStatus";

interface IconWrapperProps {
  icon: LucideIcon;
}

const IconWrapper = ({ icon: Icon }: IconWrapperProps) => {
  return (
    <div className="flex justify-center items-center p-1 w-12 h-12 rounded-sm border-1 border-border bg-surface-secondary text-foreground">
      <Icon />
    </div>
  );
};

export default function ServerPage() {
  const navigate = useNavigate();
  const session = useSession();

  const onAction = (key: Key) => {
    if (key === "import_music") {
      navigate("/import");
      return;
    } else if (key === "admin") {
      navigate("/server-admin");
      return;
    }
  };

  return (
    <PageTitleLayout title="Server">
      <div className="flex flex-col gap-4 w-full h-full">
        <ServerStatus serverInfo={session.serverInfo} forceBig={true} />
      <ListBox
        aria-label="mediamenu"
        className="overflow-visible gap-2 p-0 pr-2 pb-5 pl-4 mt-5 w-full bg-none"
        onAction={onAction}
      >
        {session.user?.hasAnyRole(
          UserRoles.Admin,
          UserRoles.ImporterWrite,
        ) ? (
          <ListBox.Item className="justify-between w-full" id="import_music" key="import_music">
            <div className="flex gap-3 items-center">
              <IconWrapper icon={CloudUpload} />
              <p className="text-xl font-semibold text-foreground">Import Music</p>
            </div>
            <ChevronRight className="stroke-foreground"/>
          </ListBox.Item>
        ) : null}
        {session.user?.hasAnyRole(UserRoles.Admin) ? (
          <ListBox.Item className="justify-between w-full" id="admin" key="admin">
            <div className="flex gap-3 items-center">
              <IconWrapper icon={Wrench} />
              <p className="text-xl font-semibold text-foreground">Administration</p>
            </div>
            <ChevronRight className="stroke-foreground"/>
          </ListBox.Item>
        ) : null}
      </ListBox>
      </div>
    </PageTitleLayout>
  );
}
