import { useEffect } from "react";
import { Icon } from "@iconify/react/dist/iconify.js";
import { type User } from "../../../shared/types/user";
import { getInitials } from "../../../shared/utils/utils";
import { useGetLocationCity } from "../hooks/useCurrentUser";

interface ProfileImageCardProps {
  user?: Partial<User> | null | undefined;
  onEditClick?: () => void;
  locationDisplay: string | null;
  setLocationDisplay: (value: string | null) => void;
  fetchingCity: boolean;
  setFetchingCity: (value: boolean) => void;
  retryCount: number;
  setRetryCount: (value: number) => void;
}

const ProfileImageCard = ({
  user,
  onEditClick,
  locationDisplay,
  setLocationDisplay,
  fetchingCity,
  setFetchingCity,
  retryCount,
  setRetryCount,
}: ProfileImageCardProps) => {
  const profilePhoto =
    user?.profile_photo ||
    (user?.photos && user.photos.length > 0 ? user.photos[0].photo_url : null);

  const initials = getInitials(user?.first_name ?? "", user?.last_name ?? "");

  const cityMutation = useGetLocationCity();

  useEffect(() => {
    if (!user?.coordinates) {
      setLocationDisplay(null);
      setRetryCount(0);
      return;
    }

    if (user.coordinates.city && user.coordinates.city.trim()) {
      setLocationDisplay(user.coordinates.city);
      setRetryCount(0);
      return;
    }

    if (
      user.coordinates.latitude &&
      user.coordinates.longitude &&
      user.coordinates.latitude !== 0 &&
      user.coordinates.longitude !== 0
    ) {
      setFetchingCity(true);

      cityMutation.mutate(
        {
          latitude: user.coordinates.latitude,
          longitude: user.coordinates.longitude,
        },
        {
          onSuccess: (data) => {
            setLocationDisplay(data.city);
            setFetchingCity(false);
            setRetryCount(0);
          },
          onError: () => {
            if (retryCount < 1) {
              // retry one more time
              setRetryCount(retryCount + 1);
            } else {
              // fallback
              if (user?.coordinates?.latitude && user?.coordinates?.longitude) {
                setLocationDisplay(
                  `(${user.coordinates.latitude.toFixed(
                    4
                  )}, ${user.coordinates.longitude.toFixed(4)})`
                );
              }
              setFetchingCity(false);
            }
          },
        }
      );
    } else {
      setLocationDisplay(null);
      setRetryCount(0);
    }
    // include retryCount in deps so we re-run on retry
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user?.coordinates, retryCount, setLocationDisplay, setFetchingCity, setRetryCount]);

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