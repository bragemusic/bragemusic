import { UserMenu } from "./UserMenu";
import { Searchbar } from "./Searchbar";

interface StatusbarProps {}

export const Statusbar: React.FC<StatusbarProps> = ({}) => {
  return (
    <div className="flex justify-between items-center h-20 border-b-1 bg-content1 border-background border-default">
      <div className="flex w-1/3" />
      <div className="px-10 w-2/3">
        <Searchbar />
      </div>
      <div className="flex gap-4 justify-end items-center pr-4 w-1/3">
        <UserMenu />
      </div>
    </div>
  );
};
