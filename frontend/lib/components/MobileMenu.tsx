import { cn } from "@heroui/styles";
import { CloudCog, House, LucideIcon, Music, Search } from "lucide-react";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Searchbar } from "./Searchbar";

interface MenuItemProps {
    id: string;
    name: string;
    icon: LucideIcon;
    selected: boolean;
    onClick: (id: string) => void;
}

const MenuItem = ({ id, name, icon: Icon, selected, onClick }: MenuItemProps) => {
    return(
        <div
            onClick={() => {onClick(id)}}
            className={cn(
            "flex flex-col items-center py-1 w-20 rounded-full gap-0.5 border transition-all duration-200 ease-out",
            selected
                    ? "bg-default/80 border-border scale-100 shadow-md"
                    : "bg-transparent border-transparent scale-95 hover:scale-100"
        )}
        >
            <Icon size={24}/>
            <p className="text-[10px]">{name}</p>
        </div>
    );
};

export const MobileMenu = () => {
    const [selectedItem, setSelectedItem] = useState("home");

    const navigate = useNavigate();

    const onMenuClick = (id: string) => {
        setSelectedItem(id)
        switch (id) {
            case "home":
                navigate("/");
                return
            case "media":
                navigate("/media");
                return
            case "server":
                navigate("/server");
                return
        }
    }

    return (
        <div className="absolute bottom-0 z-20 px-4 pb-4 w-full">
            <div className="flex justify-between p-1 w-full rounded-full shadow-md backdrop-blur-sm bg-surface/60 border-1 border-border">
                <MenuItem onClick={onMenuClick} selected={selectedItem == "home"} id={"home"} key={"home"} name="Home" icon={House}/>
                <Searchbar
                    trigger={({ onClick }) => (
                        <MenuItem onClick={onClick} selected={selectedItem == "search"} id={"search"} key={"search"} name="Search" icon={Search}/>
                    )}
                />
                <MenuItem onClick={onMenuClick} selected={selectedItem == "media"} id={"media"} key={"media"} name="Media" icon={Music}/>
                <MenuItem onClick={onMenuClick} selected={selectedItem == "server"} id={"server"} key={"server"} name="Server" icon={CloudCog}/>
            </div>
        </div>
    )
}
