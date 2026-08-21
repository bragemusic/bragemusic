import {
  CalendarSync,
  CircleAlert,
  CircleCheck,
  CircleQuestionMark,
  CircleX,
  Logs,
  LucideIcon,
  Settings,
  SquarePen,
  Trash,
  User,
  UserPlus,
} from "lucide-react";
import { useEffect, useState } from "react";
import { Table, Tabs } from "@heroui/react";
import {
  Checkbox,
  CheckboxGroup,
  Label,
  Chip,
  useOverlayState,
  toast,
} from "@heroui/react";

import { types } from "@/types/core";

import { useApi } from "@/api/ApiContext";
import TightLayout from "@/layouts/tight";
import { UserDetails } from "@/models/UserDetails";
import { ActionButton } from "@/primitives/Buttons";
import EditModal from "@/components/EditModal";
import { Input } from "@/primitives/Inputs";
import { UserBadge } from "@/components/UserBadge";

type View = "settings" | "users" | "jobs" | "entity_log";

interface ViewSettingsProps {}

const ViewSettings: React.FC<ViewSettingsProps> = ({}) => {
  return (
    <div className="flex flex-col gap-6 pt-10 w-full xl:px-40 lg:px-30">
      Settings is yet to be implemented
    </div>
  );
};

interface ViewUsersProps {}

const ViewUsers: React.FC<ViewUsersProps> = ({}) => {
  const [users, setUsers] = useState<UserDetails[]>([]);
  const [userRoles, setUserRoles] = useState<string[]>([]);
  const [selectedEditUser, setSelectedEditUser] = useState<UserDetails | null>(
    null,
  );

  const addUserState = useOverlayState();
  const editUserState = useOverlayState();

  const api = useApi();

  const loadData = () => {
    api.listUsers(false).then(setUsers);
    api.listUserRoles().then(setUserRoles);
  };

  useEffect(() => {
    loadData();
  }, [api]);

  const columns = [
    {
      key: "avatar",
      label: "",
      isHeader: false,
    },
    {
      key: "username",
      label: "Username",
      isHeader: true,
    },
    {
      key: "email",
      label: "Email",
      isHeader: false,
    },
    {
      key: "roles",
      label: "Number of Roles",
      isHeader: false,
    },
    {
      key: "provider",
      label: "Provider",
      isHeader: false,
    },
    {
      key: "created_at",
      label: "Created At",
      isHeader: false,
    },
    {
      key: "actions",
      label: "",
      isHeader: false,
    },
  ];

  const onAddUser = (fd: FormData) => {
    const email = fd.get("email");
    const username = fd.get("username");
    const password1 = fd.get("password");
    const password2 = fd.get("password2");
    const roles = fd.getAll("roles");

    if (password1 != password2) {
      toast("Cannot create user", {
        description: "Passwords does not match",
        variant: "danger",
      });

      return;
    }

    if (email == null || username == null || password1 == null) {
      toast("Cannot create user", {
        description: "All fields are not filled",
        variant: "danger",
      });

      return;
    }

    api.createUser(
      email.toString(),
      username.toString(),
      password1.toString(),
      roles.map((item) => {
        return item.toString();
      }),
    );
    loadData();
  };

  const onEditUser = (fd: FormData) => {
    const email = fd.get("email");
    const username = fd.get("username");
    const password1 = fd.get("password");
    const password2 = fd.get("password2");
    const roles = fd.getAll("roles");

    if (password1 != password2) {
      toast("Cannot update user", {
        description: "Passwords does not match",
        variant: "danger",
      });

      return;
    }

    if (email == null || username == null || selectedEditUser == null) {
      toast("Cannot update user", {
        description: "All fields are not filled",
        variant: "danger",
      });

      return;
    }

    let password: string | null | undefined = password1?.toString();

    if (password1 == undefined || password1 == "") {
      password = null;
    }

    api.updateUser(
      selectedEditUser.id,
      email.toString(),
      username.toString(),
      password,
      roles.map((item) => {
        return item.toString();
      }),
    );
    setSelectedEditUser(null);
    loadData();
  };

  return (
    <div className="flex flex-col">
      <EditModal
        header={`Edit User`}
        state={editUserState}
        onSubmit={onEditUser}
      >
        {selectedEditUser != null && (
          <>
            <Input
              required
              defaultValue={selectedEditUser.email}
              label="Email"
              name="email"
              type="email"
            />
            <Input
              required
              defaultValue={selectedEditUser.username}
              label="Username"
              name="username"
            />
            <Input label="Password" name="password" type="password" />
            <Input label="Re-enter password" name="password2" type="password" />
            <p className="pt-4">Roles</p>
            <CheckboxGroup
              className="overflow-y-scroll w-full"
              defaultValue={selectedEditUser.role}
              name="roles"
            >
              {userRoles.map((role) => {
                return (
                  <Checkbox value={role} variant="secondary">
                    <Checkbox.Control>
                      <Checkbox.Indicator />
                    </Checkbox.Control>
                    <Checkbox.Content>
                      <Label htmlFor={role}>{role}</Label>
                    </Checkbox.Content>
                  </Checkbox>
                );
              })}
            </CheckboxGroup>
          </>
        )}
      </EditModal>
      <EditModal header={`Add User`} state={addUserState} onSubmit={onAddUser}>
        <Input required label="Email" name="email" type="email" />
        <Input required label="Username" name="username" />
        <Input required label="Password" name="password" type="password" />
        <Input
          required
          label="Re-enter password"
          name="password2"
          type="password"
        />
        <p className="pt-4">Roles</p>
        <CheckboxGroup className="overflow-y-scroll w-full" name="roles">
          {userRoles.map((role) => {
            return (
              <Checkbox value={role} variant="secondary">
                <Checkbox.Control>
                  <Checkbox.Indicator />
                </Checkbox.Control>
                <Checkbox.Content>
                  <Label htmlFor={role}>{role}</Label>
                </Checkbox.Content>
              </Checkbox>
            );
          })}
        </CheckboxGroup>
      </EditModal>
      <div className="flex justify-end mb-5">
        <ActionButton
          icon={UserPlus}
          tooltip="Add new user"
          onClick={() => {
            addUserState.setOpen(true);
          }}
        />
      </div>
      <Table>
        <Table.ScrollContainer>
          <Table.Content aria-label="Users" className="">
            <Table.Header>
              {columns.map((column) => (
                <Table.Column isRowHeader={column.isHeader} key={column.key}>{column.label}</Table.Column>
              ))}
            </Table.Header>
            <Table.Body>
              {users.map((item) => (
                <Table.Row key={item.id}>
                  <Table.Cell className="flex gap-2 items-center">
                    <UserBadge size="md" user={item} />
                    {item.isAdmin &&
                      item.id != "11111111-1111-1111-1111-111111111111" && (
                        <Chip color="success" size="sm" variant="secondary">
                          admin
                        </Chip>
                      )}
                    {item.id == "11111111-1111-1111-1111-111111111111" && (
                      <Chip color="warning" size="sm" variant="primary">
                        super admin
                      </Chip>
                    )}
                  </Table.Cell>
                  <Table.Cell className="items-center">
                    {item.username}
                  </Table.Cell>
                  <Table.Cell className="items-center">{item.email}</Table.Cell>
                  <Table.Cell className="items-center">
                    {item.role ? item.role.length : 0}
                  </Table.Cell>
                  <Table.Cell className="items-center">
                    {item.provider}
                  </Table.Cell>
                  <Table.Cell>
                    {new Date(item.created_at).toISOString()}
                  </Table.Cell>
                  <Table.Cell className="items-center">
                    {item.id != "11111111-1111-1111-1111-111111111111" && (
                      <>
                        <ActionButton
                          buttonClassName="mr-1"
                          icon={SquarePen}
                          size="sm"
                          tooltip="Edit user"
                          onClick={() => {
                            setSelectedEditUser(item);
                            editUserState.setOpen(true);
                          }}
                        />
                        <ActionButton
                          confirm
                          icon={Trash}
                          size="sm"
                          tooltip="Delete user"
                          onClick={() => {
                            api.deleteUser(item.id);
                            loadData();
                          }}
                        />
                      </>
                    )}
                  </Table.Cell>
                </Table.Row>
              ))}
            </Table.Body>
          </Table.Content>
        </Table.ScrollContainer>
      </Table>
    </div>
  );
};

interface ViewJobsProps {}

const ViewJobs: React.FC<ViewJobsProps> = ({}) => {
  return (
    <div className="flex flex-col gap-6 pt-10 w-full xl:px-40 lg:px-30">
      Jobs is yet to be implemented
    </div>
  );
};

const columns = [
  {
    key: "event_type",
    label: "Action",
    isHeader: false,
  },
  {
    key: "entity_type",
    label: "Entity Type",
    isHeader: false,
  },
  {
    key: "item_id",
    label: "Entity ID",
    isHeader: true,
  },
  {
    key: "user_id",
    label: "User",
    isHeader: false,
  },
  {
    key: "event_time",
    label: "Timestamp",
    isHeader: false,
  },
];

interface ViewEntityLogProps {}

const ViewEntityLog: React.FC<ViewEntityLogProps> = ({}) => {
  const [events, setEvents] = useState<types.EntityEvent[]>([]);
  const [users, setUsers] = useState<UserDetails[]>([]);

  const api = useApi();

  const loadData = () => {
    api.listUsers(true).then(setUsers);
    api.listEntityEvents().then(setEvents);
  };

  const chip = (action: string) => {
    let color: "default" | "success" | "danger" | "warning" = "default";
    let icon = <CircleQuestionMark size={18} />;

    switch (action) {
      case "create":
        color = "success";
        icon = <CircleCheck size={18} />;
        break;
      case "update":
        color = "warning";
        icon = <CircleAlert size={18} />;
        break;
      case "delete":
        color = "danger";
        icon = <CircleX size={18} />;
        break;
    }

    return (
      <Chip color={color} size="sm" variant="tertiary">
        {icon}
        {action}
      </Chip>
    );
  };

  const user = (id: string): UserDetails | undefined => {
    return users.find((u) => u.id === id);
  };

  useEffect(() => {
    loadData();
  }, [api]);

  return (
    <Table>
      <Table.ScrollContainer>
        <Table.Content aria-label="Entity Log" className="">
          <Table.Header>
            {columns.map((column) => (
              <Table.Column isRowHeader={column.isHeader} key={column.key}>{column.label}</Table.Column>
            ))}
          </Table.Header>
          <Table.Body>
            {events.map((item) => (
              <Table.Row key={item.id}>
                <Table.Cell className="items-center">
                  {chip(item.event_type)}
                </Table.Cell>
                <Table.Cell>{item.entity_type}</Table.Cell>
                <Table.Cell>{item.item_id}</Table.Cell>
                <Table.Cell>
                  <UserBadge size="sm" user={user(item.user_id)} />
                </Table.Cell>
                <Table.Cell>
                  {new Date(item.event_time).toISOString()}
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table.Content>
      </Table.ScrollContainer>
    </Table>
  );
};

export default function ServerAdminPage() {
  const views = {
    settings: <ViewSettings />,
    users: <ViewUsers />,
    jobs: <ViewJobs />,
    entity_log: <ViewEntityLog />,
  };

  const navbarItem = (key: View, Icon: LucideIcon, desc: string) => {
    return (
      <Tabs.Tab id={key}>
        <Icon size={17} />
        <p className="pl-1">{desc}</p>
        <Tabs.Indicator />
      </Tabs.Tab>
    );
  };

  return (
    <TightLayout>
      <Tabs className="w-full" variant="secondary">
        <Tabs.ListContainer className="self-center pt-6 sm:w-150">
          <Tabs.List aria-label="Options">
            {navbarItem("settings", Settings, "Settings")}
            {navbarItem("users", User, "Users")}
            {navbarItem("jobs", CalendarSync, "Jobs")}
            {navbarItem("entity_log", Logs, "Entity Log")}
          </Tabs.List>
        </Tabs.ListContainer>
        <Tabs.Panel
          className="flex flex-col gap-6 py-10 pt-4 w-full xl:px-40 lg:px-30"
          id="settings"
        >
          {views["settings"]}
        </Tabs.Panel>
        <Tabs.Panel
          className="flex flex-col gap-6 py-10 pt-4 w-full xl:px-40 lg:px-30"
          id="users"
        >
          {views["users"]}
        </Tabs.Panel>
        <Tabs.Panel
          className="flex flex-col gap-6 py-10 pt-4 w-full xl:px-40 lg:px-30"
          id="jobs"
        >
          {views["jobs"]}
        </Tabs.Panel>
        <Tabs.Panel
          className="flex flex-col gap-6 py-10 pt-4 w-full xl:px-40 lg:px-30"
          id="entity_log"
        >
          {views["entity_log"]}
        </Tabs.Panel>
      </Tabs>
    </TightLayout>
  );
}
