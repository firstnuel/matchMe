import { getRecommendations, getUserProfile, getDistanceFromUser } from './../api/matches';
import { useAuthStore } from '../../auth/hooks/authStore';
import { useQuery } from '@tanstack/react-query';
import { useCurrentUser } from '../../userProfile/hooks/useCurrentUser';


export const useUserRecommendations = () => {
    const { data } = useCurrentUser()
    const q = useQuery({
        queryKey: ['userRecommendations'],
        queryFn: getRecommendations,
        enabled: !!(data && 'user' in data && data.user?.profile_completion && data.user.profile_completion > 90),
        retry: false,
    });

    return q;
}

export const useUserProfile = (id: string) => {
    const { authToken } = useAuthStore();
    const q = useQuery({
        queryKey: ['userProfile', id],
        queryFn: ({ queryKey }) => getUserProfile(queryKey[1] as string),
        enabled: !!(authToken && id),
        retry: false,
    });

    return q;
}

export const useUserDistance = (id: string) => {
    const { authToken } = useAuthStore();
    const q = useQuery({
        queryKey: ['userDist', id],
        queryFn: ({ queryKey }) => getDistanceFromUser(queryKey[1] as string),
        enabled: !!(authToken && id),
        retry: false,
    });

    return q;
}