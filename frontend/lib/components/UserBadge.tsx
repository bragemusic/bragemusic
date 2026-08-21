import { Avatar, AvatarProps, cn } from "@heroui/react";
import { useEffect, useState } from "react";

import { types } from "@/types/core";

import { userAvatarLink } from "@/util/images";

interface UserBadgeProps {
  user?: types.UserDetails;
  size?: AvatarProps["size"];
  radius?: "sm" | "md" | "lg" | "xl";
  showAsLink?: boolean;
}

export const UserBadge: React.FC<UserBadgeProps> = ({
  user,
  size = "lg",
  radius,
  showAsLink = false,
}) => {
  const [initials, setInitials] = useState("??");

  useEffect(() => {
    if (user && user.username !== "") {
      setInitials(
        user.username
          .trim()
          .split(/\s+/)
          .map((word) => word[0].toUpperCase())
          .join(""),
      );
    }
  }, [user]);

  return (
    <Avatar
      className={cn(
        showAsLink ? `cursor-pointer` : "",
        radius !== undefined ? `rounded-${radius}` : "",
      )}
      color="accent"
      size={size}
      variant="soft"
    >
      {user && <Avatar.Image className="object-cover" src={userAvatarLink(user.id, 320)} />}
      <Avatar.Fallback>{initials}</Avatar.Fallback>
    </Avatar>
  );
};
