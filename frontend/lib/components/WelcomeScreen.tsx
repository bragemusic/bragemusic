// import { useEffect, useState } from "react";
// import {
//   Checkbox,
//   Description,
//   Form,
//   Input,
//   Label,
//   TextField,
// } from "@heroui/react";
// import { Button } from "@heroui/react";

// import { types } from "@/types/core";

// import Logo from "../icons/logo.svg?react";

// import { UserBadge } from "./UserBadge";

// import { useApi } from "@/api/ApiContext";
// import { useSession } from "@/session/SessionContext";

// interface ViewUserSelectProps {
//   users: types.UserDetails[];
//   serverAvailable: boolean;
//   setView: (v: "user_select" | "login") => void;
// }

// const ViewUserSelect: React.FC<ViewUserSelectProps> = ({
//   users,
//   serverAvailable,
//   setView,
// }) => {
//   const api = useApi();

//   return (
//     <>
//       <div className="flex gap-10">
//         {users.map((u: types.UserDetails) => {
//           return (
//             <div
//               key={u.id}
//               className="flex flex-col items-center cursor-pointer"
//               onClick={() => {
//                 api.selectLocalUser(u.id);
//               }}
//             >
//               <UserBadge user={u} />
//               <p className="pt-2 font-bold">{u.username}</p>
//             </div>
//           );
//         })}
//       </div>
//       <Button
//         className="mt-15"
//         isDisabled={!serverAvailable}
//         size="lg"
//         variant="primary"
//         onPressEnd={() => {
//           setView("login");
//         }}
//       >
//         Log In
//       </Button>
//     </>
//   );
// };

// interface ViewLoginProps {
//   setView: (v: "user_select" | "login") => void;
// }

// const ViewLogin: React.FC<ViewLoginProps> = ({ setView }) => {
//   const [rememberMe, setRememberMe] = useState(true);

//   const api = useApi();

//   return (
//     <Form
//       className="flex flex-col gap-4 w-full"
//       onReset={() => {}}
//       onSubmit={(e) => {
//         e.preventDefault();
//         const formData = new FormData(e.currentTarget);

//         const email = formData.get("email") as string;
//         const password = formData.get("password") as string;

//         api.loginServerUser(email, password, rememberMe);
//       }}
//     >
//       <TextField isRequired name="email" type="email">
//         <Label>Email</Label>
//         <Input />
//         <Description />
//       </TextField>

//       <TextField isRequired name="password" type="password">
//         <Label>Password</Label>
//         <Input />
//         <Description />
//       </TextField>
//       <Checkbox
//         className="mt-2"
//         isSelected={rememberMe}
//         name="long_lived_token"
//         onChange={setRememberMe}
//       >
//         <Checkbox.Control>
//           <Checkbox.Indicator />
//         </Checkbox.Control>
//         <Checkbox.Content>
//           <Label>Keep me logged in</Label>
//         </Checkbox.Content>
//       </Checkbox>
//       <div className="flex gap-2 justify-end mt-2 mb-4 w-full">
//         <Button
//           variant="danger"
//           onPress={() => {
//             setView("user_select");
//           }}
//         >
//           Back
//         </Button>
//         <Button type="submit" variant="primary">
//           Log in
//         </Button>
//       </div>
//     </Form>
//   );
// };

// interface WelcomeScreenProps {}

// export const WelcomeScreen: React.FC<WelcomeScreenProps> = ({}) => {
//   const [users, setUsers] = useState<types.UserDetails[]>([]);
//   const [view, setView] = useState<"user_select" | "login">("user_select");
//   const [serverAvailable, setServerAvailable] = useState<boolean>(false);

//   const api = useApi();
//   const session = useSession();

//   useEffect(() => {
//     api.getCachedUsers().then(setUsers);
//   }, [api]);

//   useEffect(() => {
//     setServerAvailable(session.serverInfo.status != "unavailable");
//   }, [session.serverInfo]);

//   const views = {
//     user_select: (
//       <ViewUserSelect
//         serverAvailable={serverAvailable}
//         setView={setView}
//         users={users}
//       />
//     ),
//     login: <ViewLogin setView={setView} />,
//   };

//   const titles = {
//     user_select: "Welcome to Brage Music",
//     login: "Login to server",
//   };

//   const subtitles = {
//     user_select: "select a logged in user or log in a new",
//     login: `${session.serverInfo.name}, ${session.serverInfo.version}`,
//   };

//   return (
//     <div className="flex flex-col justify-center items-center w-full h-full bg-surface-secondary text-foreground fill-foreground">
//       <div className="flex flex-col justify-between items-center p-10 w-full h-full rounded-xl shadow-xl md:w-1/2 md:h-1/2 bg-background">
//         <Logo className="mb-8 fill-content2" height={140} width={140} />
//         <p className="pb-2 text-2xl font-semibold md:text-3xl lg:text-5xl">
//           {titles[view]}
//         </p>
//         <p className="pb-20 lg:text-xl">{subtitles[view]}</p>
//         {views[view]}
//       </div>
//     </div>
//   );
// };
