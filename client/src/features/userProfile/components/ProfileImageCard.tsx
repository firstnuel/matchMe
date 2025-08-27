import { Icon } from "@iconify/react/dist/iconify.js";
import { type User } from "../../../shared/types/user";
import { getInitials } from "../../../shared/utils/utils";
import { useCityFromCoordinates } from "../hooks/useCityFromCoordinates"; // path to custom hook

interface ProfileImageCardProps {
  user?: Partial<User> | null | undefined;
  onEditClick?: () => void;
}

const ProfileImageCard = ({ user, onEditClick }: ProfileImageCardProps) => {
  const profilePhoto =
    user?.profile_photo ||
    (user?.photos && user.photos.length > 0 ? user.photos[0].photo_url : null);
  const initials = getInitials(user?.first_name ?? "", user?.last_name ?? "");
  const { locationDisplay, fetchingCity } = useCityFromCoordinates(user?.coordinates);

  return (
    <div className="profile-view">
      <div className="profile-header">
        <div
          className="profile-avatar"
          style={
            profilePhoto
              ? {
                  backgroundImage: `url(${profilePhoto})`,
                  backgroundSize: "cover",
                  backgroundPosition: "center",
                  color: "transparent",
                }
              : {}
          }
        >
          {profilePhoto ? "" : initials}
        </div>
        <div className="profile-name">
          {user?.first_name && user?.last_name
            ? `${user.first_name} ${user.last_name}`
            : "Unknown User"}
        </div>
        <div className="profile-details">
          {user?.age && `${user.age} years old`}
          {user?.age && locationDisplay && (
            <Icon icon="mdi:circle" className="sep-icon" />
          )}
          {fetchingCity ? "Loading location..." : locationDisplay}
        </div>
        <button className="edit-profile-btn" onClick={onEditClick}>
          Edit Profile
        </button>
      </div>
    </div>
  );
};

export default ProfileImageCard;
