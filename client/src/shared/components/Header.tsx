import { useState } from "react";
import { useCurrentUser } from "../../features/userProfile/hooks/useCurrentUser";
import { useAuthStore } from "../../features/auth/hooks/authStore";
import { getInitials } from "../utils/utils";
import { useUIStore } from "../hooks/uiStore";


const Header = () => {
  const [show, setShow] = useState(false);
  const { clearAuth, } = useAuthStore();
  const { data: currentUser } = useCurrentUser();
  const { profileView } = useUIStore()
  const user = currentUser && 'user' in currentUser ? currentUser.user : undefined;

  const handleLogout = () => {
    clearAuth()
    setShow(false);
  };

  return (
    <header className="header">
      <div className="app-name logo">MatchMe</div>

      <div
        className="profile-icon"
        onClick={() => setShow(!show)}
        role="button"
        aria-expanded={show}
      >
        {getInitials(user?.first_name ?? '',user?.last_name ?? '')}
      </div>

      {show && (
        <div className="profile-menu">
          <div className="profile-name">
            {`${user?.first_name ?? ''} ${user?.last_name ?? ''}`}
          </div>
          <button className="profile-btn logout" onClick={handleLogout}>
            Logout
          </button>
          <button className="profile-btn cancel" onClick={() => setShow(false)}>
            Cancel
          </button>
        </div>
      )}
    </header>
  );
};

export default Header;
