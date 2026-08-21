// import { Button, PressEvent } from "@heroui/react";
// import { Card, CardBody } from "@heroui/card";
// import { LucideIcon } from "lucide-react";

// interface ActionCardProps {
//     title: string
//     desc: string
//     icon: LucideIcon;
//     onClick?: (e: PressEvent) => void;
// }

// export const ActionCard: React.FC<ActionCardProps> = ({ title, desc, icon: Icon, onClick}) => {
//   return (
//     <Card
//       className="border-1 border-foreground/30 hover:bg-default/30"
//       radius="lg"
//       as={Button}
//       onPressEnd={onClick}
//     >
//       <CardBody>
//         <div className="flex gap-5 items-center">
//           <div className="p-2 rounded-lg bg-background">
//             <Icon size={50} />
//           </div>
//           <div className="flex flex-col">
//             <p className="text-xl font">{title}</p>
//             <p className="text-content4">{desc}</p>
//           </div>
//         </div>
//       </CardBody>
//     </Card>
//   );
// };
