import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentUser, updateUser, uploadUserPhotos, deleteUserPhoto, getLocationCity } from '../api/userProfile';
import { useAuthStore } from '../../auth/hooks/authStore';
import { type UpdateUserRequest } from '../types/user';
import { useUIStore } from '../../../shared/hooks/uiStore';


export const useCurrentUser = () => {
  const { authToken } = useAuthStore();
  return useQuery({
    queryKey: ['currentUser'],
    queryFn: getCurrentUser,
    enabled: !!authToken,
    retry: false,
  });
};

export const useUpdateUser = () => {
  const queryClient = useQueryClient();
   const { setInfo, setError } = useUIStore();
  
  return useMutation({
    mutationFn: (userData: Partial<UpdateUserRequest>) => updateUser(userData),
    onSuccess: () => {
      setInfo("User updated successfully")
      queryClient.invalidateQueries({ queryKey: ['currentUser'] });
      queryClient.refetchQueries({ queryKey: ['currentUser'] });
    },
    onError: (err: Error) => {
      setError(err.message)
    }
  });
};

export const useUploadPhotos = () => {
  const queryClient = useQueryClient();
  const { setInfo, setError } = useUIStore();
  
  return useMutation({
    mutationFn: (photos: File[]) => uploadUserPhotos(photos),
    onSuccess: () => {
      setInfo("User photo updated successfully")
      queryClient.invalidateQueries({ queryKey: ['currentUser'] });
      queryClient.refetchQueries({ queryKey: ['currentUser'] });
    },
    onError: (err: Error) => {
      setError(err.message)
    }
  });
};

export const useDeletePhoto = () => {
  const queryClient = useQueryClient();
  const { setInfo, setError } = useUIStore();
  
  return useMutation({
    mutationFn: (photoId: string) => deleteUserPhoto(photoId),
    onSuccess: () => {
      setInfo("User photo deleted successfully")
      queryClient.invalidateQueries({ queryKey: ['currentUser'] });
      queryClient.refetchQueries({ queryKey: ['currentUser'] });
    },
    onError: (err: Error) => {
      setError(err.message)
    }
  });
};

export const useGetLocationCity = () => {
  return useMutation({
    mutationFn: ({ latitude, longitude }: { latitude: number; longitude: number }) => 
      getLocationCity({ latitude, longitude }),
  });
};
