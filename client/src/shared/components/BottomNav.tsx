import NavItem from "./NavItem"



const BottomNav = () => {
    const navOpts: ["home", "chat", "connections", "profile"] = 
    ["home", "chat", "connections", "profile"]

    return (
        <nav className="bottom-nav">
            { navOpts.map((nvItem, idx) => (<NavItem key={idx} name={nvItem}/>)) }
        </nav>
    )
}


export default BottomNav;