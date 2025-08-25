import { getRecommendations, getUserProfile, getDistanceFromUser } from './../api/matches';
import { useAuthStore } from '../../auth/hooks/authStore';
import { useQuery } from '@tanstack/react-query';
import { useCurrentUser } from '../../userProfile/hooks/useCurrentUser';

// Fetch user recommendations (only if profile completion > 90)
export const useUserRecommendations = () => {
  const { data } = useCurrentUser();
  const isProfileComplete =
    data &&
    'user' in data &&
    data.user?.profile_completion &&
    data.user.profile_completion > 90;

  return useQuery({
    queryKey: ['userRecommendations'],
    queryFn: getRecommendations,
    enabled: !!isProfileComplete,
    retry: false,
  });
};

// Fetch another user's profile by id
export const useUserProfile = (id: string) => {
  const { authToken } = useAuthStore();
  return useQuery({
    queryKey: ['userProfile', id],
    queryFn: () => getUserProfile(id),
    enabled: !!(authToken && id),
    retry: false,
  });
};

// Fetch distance between current user and another user
export const useUserDistance = (id: string) => {
  const { authToken } = useAuthStore();
  return useQuery({
    queryKey: ['userDist', id],
    queryFn: () => getDistanceFromUser(id),
    enabled: !!(authToken && id),
    retry: false,
  });
};
