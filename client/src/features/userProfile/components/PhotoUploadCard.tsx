import { useRef, useState } from "react";
import { Icon } from "@iconify/react/dist/iconify.js";

interface PhotoUploadProps {
  maxPhotos?: number;
}

const PhotoUploadSection = ({ maxPhotos = 6 }: PhotoUploadProps) => {
  const [photos, setPhotos] = useState<string[]>([]);
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const [dragIndex, setDragIndex] = useState<number | null>(null);

  const handleRemove = (index: number) => {
    setPhotos((prev) => prev.filter((_, i) => i !== index));
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files) return;

    const fileArray: string[] = [];
    for (let i = 0; i < files.length; i++) {
      const file = files[i];
      const url = URL.createObjectURL(file);
      fileArray.push(url);
    }

    setPhotos((prev) => [...prev, ...fileArray].slice(0, maxPhotos));
    e.target.value = "";
  };

  const handleAdd = () => {
    fileInputRef.current?.click();
  };

  const handleDragStart = (index: number) => {
    setDragIndex(index);
  };

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
  };

  const handleDrop = (index: number) => {
    if (dragIndex === null || dragIndex === index) return;

    const newPhotos = [...photos];
    const [draggedPhoto] = newPhotos.splice(dragIndex, 1);
    newPhotos.splice(index, 0, draggedPhoto);
    setPhotos(newPhotos);
    setDragIndex(null);
  };

  const slots = Array.from({ length: maxPhotos }, (_, i) => {
    const photo = photos[i];
    return photo ? (
      <div
        key={i}
        className="photo-slot has-photo"
        style={{ backgroundImage: `url(${photo})` }}
        draggable={true}
        onDragStart={() => handleDragStart(i)}
        onDragOver={handleDragOver}
        onDrop={() => handleDrop(i)}
      >
        <button className="photo-remove" onClick={() => handleRemove(i)}>
          ×
        </button>
      </div>
    ) : (
      <div key={i} className="photo-slot" onClick={handleAdd}>
        <div className="upload-icon">
          <Icon icon="mdi:camera" className="icon" />
        </div>
      </div>
    );
  });

  return (
      <div>
        <div className="photo-upload-section">
          <div className="photo-grid">{slots}</div>
          <input
              type="file"
              ref={fileInputRef}
              accept="image/*"
              multiple
              style={{ display: "none" }}
              onChange={handleFileChange} />
        </div>
        <button className="upload-btn" onClick={() => { } }>
              Update Photo(s)
          </button>
    </div>

  );
};

export default PhotoUploadSection;