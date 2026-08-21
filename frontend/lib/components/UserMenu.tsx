import { Dropdown, Label } from "@heroui/react";
import { LogOut, UserCog } from "lucide-react";
import { useNavigate } from "react-router-dom";

import { UserBadge } from "./UserBadge";

import { useApi } from "@/api/ApiContext";
import { useSession } from "@/session/SessionContext";
import { DropdonwMenuItem } from "@/primitives/DropdownItems/DropdownMenuItem";

interface UserMenuProps {}

export const UserMenu: React.FC<UserMenuProps> = ({}) => {
  const api = useApi();
  const session = useSession();
  const navigate = useNavigate();

  const dropDownCallback = (callbackType: React.Key) => {
    switch (callbackType.toString()) {
      case "settings":
        navigate("/settings");
        return;
      case "logout":
        api.logoutLocalUser();
        return;
    }
  };

  return (
    <div className="flex items-center select-none">
      {session.user ? (
        <Dropdown className="mt-4">
          <Dropdown.Trigger>
            <UserBadge radius="sm" user={session.user} />
          </Dropdown.Trigger>
          <Dropdown.Popover>
            <Dropdown.Menu
              aria-label="Profile Actions"
              disabledKeys={
                session.serverInfo.status === "unavailable"
                  ? ["login_server"]
                  : []
              }
              onAction={dropDownCallback}
            >
              <Dropdown.Item
                className="gap-2 h-14"
                id="profile"
                textValue="Profile"
              >
                <Label>
                  <p className="font-semibold">{session.user.username}</p>
                  <p className="font-semibold">{session.user.email}</p>
                </Label>
              </Dropdown.Item>

              <DropdonwMenuItem
                description="Change your personal settings"
                icon={UserCog}
                id="settings"
                label="Settings"
              />
              <DropdonwMenuItem
                description="Logout the current user"
                icon={LogOut}
                id="logout"
                label="Logout"
              />
            </Dropdown.Menu>
          </Dropdown.Popover>
        </Dropdown>
      ) : (
        <></>
      )}
    </div>
  );
};
