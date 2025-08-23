import { useRef, useState, useEffect } from "react";
import { Icon } from "@iconify/react/dist/iconify.js";
import { useUploadPhotos, useDeletePhoto } from "../hooks/useCurrentUser";
import { type UserPhoto } from "../types/user";

interface PhotoUploadProps {
  maxPhotos?: number;
  existingPhotos?: UserPhoto[];
}

const PhotoUploadSection = ({ maxPhotos = 6, existingPhotos = [] }: PhotoUploadProps) => {
  const [photos, setPhotos] = useState<string[]>([]);
  const [photoFiles, setPhotoFiles] = useState<File[]>([]);
  const [serverPhotos, setServerPhotos] = useState<UserPhoto[]>([]);
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const [dragIndex, setDragIndex] = useState<number | null>(null);
  const uploadPhotosMutation = useUploadPhotos();
  const deletePhotoMutation = useDeletePhoto();

  // Initialize with existing photos from server
  useEffect(() => {
    if (existingPhotos.length > 0) {
      // Sort by order and extract URLs
      const sortedPhotos = [...existingPhotos].sort((a, b) => a.order - b.order);
      const photoUrls = sortedPhotos.map(photo => photo.photo_url);
      setPhotos(photoUrls);
      setServerPhotos(sortedPhotos);
      setPhotoFiles([]); // Clear any local files when server photos are loaded
    }
  }, [existingPhotos]);

  const handleRemove = (index: number) => {
    const serverPhoto = serverPhotos[index];
    
    // If this is a server photo, delete it via API
    if (serverPhoto?.id) {
      deletePhotoMutation.mutate(serverPhoto.id);
    }
    
    // Remove from local state regardless (for immediate UI feedback)
    setPhotos((prev) => prev.filter((_, i) => i !== index));
    setPhotoFiles((prev) => prev.filter((_, i) => i !== index));
    setServerPhotos((prev) => prev.filter((_, i) => i !== index));
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files) return;

    const fileArray: string[] = [];
    const newFiles: File[] = [];
    
    for (let i = 0; i < files.length; i++) {
      const file = files[i];
      const url = URL.createObjectURL(file);
      fileArray.push(url);
      newFiles.push(file);
    }

    setPhotos((prev) => [...prev, ...fileArray].slice(0, maxPhotos));
    setPhotoFiles((prev) => [...prev, ...newFiles].slice(0, maxPhotos));
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
    const newPhotoFiles = [...photoFiles];
    const newServerPhotos = [...serverPhotos];
    
    const [draggedPhoto] = newPhotos.splice(dragIndex, 1);
    const [draggedFile] = newPhotoFiles.splice(dragIndex, 1);
    const [draggedServerPhoto] = newServerPhotos.splice(dragIndex, 1);
    
    newPhotos.splice(index, 0, draggedPhoto);
    newPhotoFiles.splice(index, 0, draggedFile);
    newServerPhotos.splice(index, 0, draggedServerPhoto);
    
    setPhotos(newPhotos);
    setPhotoFiles(newPhotoFiles);
    setServerPhotos(newServerPhotos);
    setDragIndex(null);
  };

  const handleUpload = () => {
    const newFilesToUpload = photoFiles.filter(file => file instanceof File);
    if (newFilesToUpload.length === 0) {
      return;
    }
    uploadPhotosMutation.mutate(newFilesToUpload);
  };

  const hasNewFiles = photoFiles.some(file => file instanceof File);

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
        <button 
          className="photo-remove" 
          onClick={() => handleRemove(i)}
          disabled={!!(deletePhotoMutation.isPending && serverPhotos[i]?.id)}
        >
          {deletePhotoMutation.isPending && serverPhotos[i]?.id ? '...' : '×'}
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
        <button 
          className="upload-btn" 
          onClick={handleUpload}
          disabled={!hasNewFiles || uploadPhotosMutation.isPending}
        >
          {uploadPhotosMutation.isPending ? "Uploading..." : "Update Photo(s)"}
        </button>
    </div>

  );
};

export default PhotoUploadSection;