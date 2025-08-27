import { useState, useRef, useEffect } from "react";
import { useCurrentUser } from "../../features/userProfile/hooks/useCurrentUser";
import { useAuthStore } from "../../features/auth/hooks/authStore";
import { getInitials } from "../utils/utils";
import { useNavigate } from "react-router";

const Header = () => {
  const [show, setShow] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);
  const { clearAuth, } = useAuthStore();
  const { data: currentUser } = useCurrentUser();
  const user = currentUser && 'user' in currentUser ? currentUser.user : undefined;
  const navigate = useNavigate()
  
  // Get the first photo (main profile photo) or profile_photo if available
  const profilePhoto = user?.profile_photo || (user?.photos && user.photos.length > 0 ? user.photos[0].photo_url : null);

  const handleLogout = () => {
    clearAuth()
    setShow(false);
  };

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setShow(false);
      }
    };

    if (show) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [show]);

  return (
    <header className="header">
      <div className="app-name logo" onClick={() => navigate("/")}>MatchMe</div>

      <div
        className="profile-icon"
        onClick={() => setShow(!show)}
        role="button"
        aria-expanded={show}
        style={profilePhoto ? {
          backgroundImage: `url(${profilePhoto})`,
          backgroundSize: 'cover',
          backgroundPosition: 'center',
          color: 'transparent'
        } : {}}
      >
        {profilePhoto ? '' : getInitials(user?.first_name ?? '',user?.last_name ?? '')}
      </div>

      {show && (
        <div className="profile-menu" ref={menuRef}>
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
