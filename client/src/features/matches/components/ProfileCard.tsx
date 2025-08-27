import { useState } from "react";
import { Icon } from "@iconify/react/dist/iconify.js";
import { type User } from "../../../shared/types/user";
import { getInitials } from "../../../shared/utils/utils";
import ImageCarousel from "../../../shared/components/ImageCarousel";
import { useNavigate } from "react-router";
import { firstToUpper } from "../../../shared/utils/utils";

interface ProfileCardProps {
  user?: (Partial<User> & { distance?: number }) | null;
  onLike?: () => void;
  onReject?: () => void;
}

type SwipeDirection = 'left' | 'right' | null;

const ProfileCard = ({ user, onLike, onReject }: ProfileCardProps) => {
  const [swipeDirection, setSwipeDirection] = useState<SwipeDirection>(null);
  const [isAnimating, setIsAnimating] = useState(false);
  const navigate = useNavigate()

  // Get all photos from user
  const photos = user?.photos || [];
  const fallbackPhoto = user?.profile_photo;

  // Get user initials for fallback
  const initials = getInitials(user?.first_name ?? '', user?.last_name ?? '');

  // Swipe animation handlers
  const handleLikeWithSwipe = () => {
    if (isAnimating) return;
    
    setIsAnimating(true);
    setSwipeDirection('right');
    
    // Call the actual like handler after a short delay
    setTimeout(() => {
      onLike?.();
    }, 300);
    
    // Reset animation state after animation completes
    setTimeout(() => {
      setSwipeDirection(null);
      setIsAnimating(false);
    }, 600);
  };

  const handleRejectWithSwipe = () => {
    if (isAnimating) return;
    
    setIsAnimating(true);
    setSwipeDirection('left');
    
    // Call the actual reject handler after a short delay
    setTimeout(() => {
      onReject?.();
    }, 300);
    
    // Reset animation state after animation completes
    setTimeout(() => {
      setSwipeDirection(null);
      setIsAnimating(false);
    }, 600);
  };

  return (
    <div className={`profile-card ${swipeDirection ? `swipe-${swipeDirection}` : ''} ${isAnimating ? 'animating' : ''}`}>
      <div className="card-image">
        {photos.length > 0 || fallbackPhoto ? (
          <ImageCarousel 
            photos={photos}
            fallbackPhoto={fallbackPhoto ?? undefined}
            altText={`${user?.first_name || 'User'} profile`}
          />
        ) : (
          <div className="profile-placeholder">
            <div className="placeholder-avatar">
              {initials}
            </div>
          </div>
        )}
      </div>

      <div className="card-info">
        <div className="card-name" onClick={() => navigate(`/users/${user?.id}`)}>
          {user?.first_name ? `${user.first_name} ${user.last_name?.[0] || ''}, ${user.age ?? ''}` : 'Unknown User'}
        </div>
        <div className="gender">{firstToUpper(user?.gender ?? '')}</div>
        <div className="card-age">
          {user?.distance !== undefined 
            ? user.distance > 0 && user.distance < 1 
              ? 'Less than 1km away'
              : user.distance === 0
                ? 'Same location'
                : `${Math.round(user.distance)}km away`
            : 'Distance unknown'}
        </div>
        {user?.about_me && (
          <div className="card-bio">
            {user.about_me}
          </div>
        )}
      </div>

      <div className="action-buttons">
        <button 
          className="action-btn reject-btn" 
          onClick={handleRejectWithSwipe}
          disabled={isAnimating}
        >
          <Icon icon="mdi:close" />
        </button>
        <button 
          className="action-btn like-btn" 
          onClick={handleLikeWithSwipe}
          disabled={isAnimating}
        >
          <Icon icon="mdi:heart" />
        </button>
      </div>
    </div>
  );
};

export default ProfileCard;