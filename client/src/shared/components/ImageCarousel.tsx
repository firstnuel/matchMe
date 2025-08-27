import { useState } from "react";
import { Icon } from "@iconify/react/dist/iconify.js";
import { type UserPhoto } from "../types/user";

interface ImageCarouselProps {
  photos: UserPhoto[];
  fallbackPhoto?: string;
  altText?: string;
  className?: string;
}

const ImageCarousel = ({ 
  photos, 
  fallbackPhoto, 
  altText = "Profile photo",
  className = ""
}: ImageCarouselProps) => {
  const [currentImageIndex, setCurrentImageIndex] = useState(0);

  // Combine photos with fallback
  const allPhotos = photos.length > 0 ? photos : (fallbackPhoto ? [{ photo_url: fallbackPhoto }] : []);
  const hasMultiplePhotos = allPhotos.length > 1;
  const currentPhoto = allPhotos[currentImageIndex]?.photo_url;

  // Navigation handlers
  const goToPrevious = (e: React.MouseEvent) => {
    e.stopPropagation();
    setCurrentImageIndex((prev) => 
      prev === 0 ? allPhotos.length - 1 : prev - 1
    );
  };

  const goToNext = (e: React.MouseEvent) => {
    e.stopPropagation();
    setCurrentImageIndex((prev) => 
      prev === allPhotos.length - 1 ? 0 : prev + 1
    );
  };

  // Dot indicators
  const goToSlide = (index: number) => {
    setCurrentImageIndex(index);
  };

  if (!currentPhoto) {
    return null;
  }

  return (
    <div className={`image-carousel ${className}`}>
      <img 
        src={currentPhoto} 
        alt={altText}
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
          {allPhotos.map((_, index) => (
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
  );
};

export default ImageCarousel;