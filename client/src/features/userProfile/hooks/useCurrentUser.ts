/* eslint-disable @typescript-eslint/no-explicit-any */
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentUser, updateUser, uploadUserPhotos, deleteUserPhoto, getLocationCity } from '../api/userProfile';
import { useAuthStore } from '../../auth/hooks/authStore';
import { type UpdateUserRequest } from '../types/user';
import { useUIStore } from '../../../shared/hooks/uiStore';

// Fetch current user
export const useCurrentUser = () => {
  const { authToken } = useAuthStore();
  return useQuery({
    queryKey: ['currentUser'],
    queryFn: getCurrentUser,
    enabled: !!authToken,
    retry: false,
  });
};

// Update user info with optimistic update
export const useUpdateUser = () => {
  const queryClient = useQueryClient();
  const { setInfo, setError } = useUIStore();

  return useMutation({
    mutationFn: (userData: Partial<UpdateUserRequest>) => updateUser(userData),
    onMutate: async (newData) => {
      await queryClient.cancelQueries({ queryKey: ['currentUser'] });
      const previousUser = queryClient.getQueryData(['currentUser']);
      queryClient.setQueryData(['currentUser'], (old: any) => ({ ...old, ...newData }));
      return { previousUser };
    },
    onError: (err: Error, _newData, context: any) => {
      if (context?.previousUser) {
        queryClient.setQueryData(['currentUser'], context.previousUser);
      }
      setError(err.message);
    },
    onSuccess: () => {
      setInfo('User updated successfully');
      queryClient.invalidateQueries({ queryKey: ['currentUser'] });
    },
  });
};

// Upload photos with optimistic update
export const useUploadPhotos = () => {
  const queryClient = useQueryClient();
  const { setInfo, setError } = useUIStore();

  return useMutation({
    mutationFn: (photos: File[]) => uploadUserPhotos(photos),
    onMutate: async (newPhotos) => {
      await queryClient.cancelQueries({ queryKey: ['currentUser'] });
      const previousUser = queryClient.getQueryData(['currentUser']);
      queryClient.setQueryData(['currentUser'], (old: any) => ({
        ...old,
        photos: [...(old?.photos || []), ...newPhotos.map((file, i) => ({ id: `temp-${Date.now()}-${i}`, url: URL.createObjectURL(file), uploading: true }))],
      }));
      return { previousUser };
    },
    onError: (err: Error, _newPhotos, context: any) => {
      if (context?.previousUser) {
        queryClient.setQueryData(['currentUser'], context.previousUser);
      }
      setError(err.message);
    },
    onSuccess: () => {
      setInfo('User photo updated successfully');
      queryClient.invalidateQueries({ queryKey: ['currentUser'] });
    },
  });
};

// Delete photo with optimistic update
export const useDeletePhoto = () => {
  const queryClient = useQueryClient();
  const { setInfo, setError } = useUIStore();

  return useMutation({
    mutationFn: (photoId: string) => deleteUserPhoto(photoId),
    onMutate: async (photoId) => {
      await queryClient.cancelQueries({ queryKey: ['currentUser'] });
      const previousUser = queryClient.getQueryData(['currentUser']);
      queryClient.setQueryData(['currentUser'], (old: any) => ({
        ...old,
        photos: (old?.photos || []).filter((p: any) => p.id !== photoId),
      }));
      return { previousUser };
    },
    onError: (err: Error, _photoId, context: any) => {
      if (context?.previousUser) {
        queryClient.setQueryData(['currentUser'], context.previousUser);
      }
      setError(err.message);
    },
    onSuccess: () => {
      setInfo('User photo deleted successfully');
      queryClient.invalidateQueries({ queryKey: ['currentUser'] });
    },
  });
};

// Get location city (mutation, not cached)
export const useGetLocationCity = () => {
  return useMutation({
    mutationFn: ({ latitude, longitude }: { latitude: number; longitude: number }) =>
      getLocationCity({ latitude, longitude }),
  });
};
