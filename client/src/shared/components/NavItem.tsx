import { Icon } from "@iconify/react/dist/iconify.js";
import { firstToUpper } from "../utils/utils";
import { useNavigate } from "react-router";
import { useUIStore } from "../hooks/uiStore";

interface NavItemProps {
    name: "chat" | "home" | "connections" | "profile";
}

const Icons = {
    chat: <Icon icon="mdi:chat" className="icon" />,
    home: <Icon icon="mdi:home" className="icon" />,
    connections: <Icon icon="mdi:account-group" className="icon" />,
    profile: <Icon icon="mdi:user" className="icon" />,
};

const NavItem = ({ name }: NavItemProps) => {
    const navigate = useNavigate();
    const { view, setView} = useUIStore()

    const handleClick = () => {
        navigate(`/${name}`);
        setView(name)
    };

    return (
        <div
            className={`nav-item ${view === name ? "active" : ""}`}
            onClick={handleClick}
        >
            <div className="nav-icon">
                {Icons[name]}
            </div>
            <div className="nav-label">{firstToUpper(name)}</div>
        </div>
    );
};

export default NavItem;
