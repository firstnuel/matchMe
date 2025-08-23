import { useState } from "react";
import { Icon } from "@iconify/react/dist/iconify.js";
import { type User } from "../../../shared/types/user";
import { getInitials } from "../../../shared/utils/utils";

interface ProfileCardProps {
  user?: (Partial<User> & { distance?: number }) | null;
  onLike?: () => void;
  onReject?: () => void;
}

type SwipeDirection = 'left' | 'right' | null;

const ProfileCard = ({ user, onLike, onReject }: ProfileCardProps) => {
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const [swipeDirection, setSwipeDirection] = useState<SwipeDirection>(null);
  const [isAnimating, setIsAnimating] = useState(false);

  // Get all photos from user
  const photos = user?.photos || [];
  const hasMultiplePhotos = photos.length > 1;

  // Get current photo or fallback
  const currentPhoto = photos[currentImageIndex]?.photo_url;
  const fallbackPhoto = user?.profile_photo;
  const displayPhoto = currentPhoto || fallbackPhoto;

  // Get user initials for fallback
  const initials = getInitials(user?.first_name ?? '', user?.last_name ?? '');

  // Navigation handlers
  const goToPrevious = (e: React.MouseEvent) => {
    e.stopPropagation();
    setCurrentImageIndex((prev) => 
      prev === 0 ? photos.length - 1 : prev - 1
    );
  };

  const goToNext = (e: React.MouseEvent) => {
    e.stopPropagation();
    setCurrentImageIndex((prev) => 
      prev === photos.length - 1 ? 0 : prev + 1
    );
  };

  // Dot indicators
  const goToSlide = (index: number) => {
    setCurrentImageIndex(index);
  };

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
        {displayPhoto ? (
          <div className="image-carousel">
            <img 
              src={displayPhoto} 
              alt={`${user?.first_name || 'User'} profile`}
              className="carousel-image"
            />
            
            {/* Navigation arrows for multiple photos */}
            {hasMultiplePhotos && (
              <>
                <button 
                  className="carousel-nav carousel-prev"
                  onClick={goToPrevious}
                  aria-label="Previous photo"
                >
                  <Icon icon="mdi:chevron-left" />
                </button>
                <button 
                  className="carousel-nav carousel-next"
                  onClick={goToNext}
                  aria-label="Next photo"
                >
                  <Icon icon="mdi:chevron-right" />
                </button>
              </>
            )}

            {/* Dot indicators for multiple photos */}
            {hasMultiplePhotos && (
              <div className="carousel-dots">
                {photos.map((_, index) => (
                  <button
                    key={index}
                    className={`carousel-dot ${index === currentImageIndex ? 'active' : ''}`}
                    onClick={() => goToSlide(index)}
                    aria-label={`Go to photo ${index + 1}`}
                  />
                ))}
              </div>
            )}
          </div>
        ) : (
          <div className="profile-placeholder">
            <div className="placeholder-avatar">
              {initials}
            </div>
          </div>
        )}
      </div>

      <div className="card-info">
        <div className="card-name">
          {user?.first_name ? `${user.first_name} ${user.last_name?.[0] || ''}, ${user.age?? ''}` : 'Unknown User'}
        </div>
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