import { types as wtypes } from "@/types/core";
import { UserRoles } from "@/types/roles";

export class UserDetails extends wtypes.UserDetails {
  get isAdmin(): boolean {
    return this.role?.includes(UserRoles.Admin) ?? false;
  }

  get isSuperAdmin(): boolean {
    return this.id == "11111111-1111-1111-1111-111111111111";
  }

  hasAnyRole(...roles: UserRoles[]): boolean {
    return this.role?.some((r) => roles.includes(r as UserRoles)) ?? false;
  }

  hasAllRoles(...roles: UserRoles[]): boolean {
    return roles.every((role) => this.role?.includes(role) ?? false);
  }
}
