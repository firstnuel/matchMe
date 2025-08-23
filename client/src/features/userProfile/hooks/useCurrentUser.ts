import { useQuery } from '@tanstack/react-query';
import { getCurrentUser } from '../api/userProfile';
import { useAuthStore } from '../../auth/hooks/authStore';

export function useCurrentUser() {
    const { authToken } = useAuthStore();
  return useQuery({
    queryKey: ['currentUser'],
    queryFn: getCurrentUser,
    enabled: !!authToken,
    retry: false,
  });
}
