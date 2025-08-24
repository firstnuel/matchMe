import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  getConnections, 
  deleteConnections, 
  getConnectionRequests, 
  sendConnectionRequest, 
  acceptConnectionRequest, 
  rejectConnectionRequest 
} from '../api/connections';
import { useAuthStore } from '../../auth/hooks/authStore';
import { type SendConnectionRequestBody } from '../types/connections';
import { useUIStore } from '../../../shared/hooks/uiStore';

export const useConnections = () => {
  const { authToken } = useAuthStore();
  return useQuery({
    queryKey: ['connections'],
    queryFn: getConnections,
    enabled: !!authToken,
    retry: false,
  });
};

export const useConnectionRequests = () => {
  const { authToken } = useAuthStore();
  return useQuery({
    queryKey: ['connectionRequests'],
    queryFn: getConnectionRequests,
    enabled: !!authToken,
    retry: false,
  });
};

export const useDeleteConnection = () => {
  const queryClient = useQueryClient();
  const { setInfo, setError } = useUIStore();
  
  return useMutation({
    mutationFn: (connectionId: string) => deleteConnections(connectionId),
    onSuccess: () => {
      setInfo("Connection deleted successfully");
      queryClient.invalidateQueries({ queryKey: ['connections'] });
      queryClient.refetchQueries({ queryKey: ['connections'] });
    },
    onError: (err: Error) => {
      setError(err.message);
    }
  });
};

export const useSendConnectionRequest = () => {
  const queryClient = useQueryClient();
  const { setInfo, setError } = useUIStore();
  
  return useMutation({
    mutationFn: (requestData: SendConnectionRequestBody) => sendConnectionRequest(requestData),
    onSuccess: () => {
      setInfo("Connection request sent successfully");
      queryClient.invalidateQueries({ queryKey: ['connectionRequests'] });
      queryClient.refetchQueries({ queryKey: ['connectionRequests'] });
    },
    onError: (err: Error) => {
      setError(err.message);
    }
  });
};

export const useAcceptConnectionRequest = () => {
  const queryClient = useQueryClient();
  const { setInfo, setError } = useUIStore();
  
  return useMutation({
    mutationFn: (requestId: string) => acceptConnectionRequest(requestId),
    onSuccess: () => {
      setInfo("Connection request accepted successfully");
      queryClient.invalidateQueries({ queryKey: ['connectionRequests'] });
      queryClient.invalidateQueries({ queryKey: ['connections'] });
      queryClient.refetchQueries({ queryKey: ['connectionRequests'] });
      queryClient.refetchQueries({ queryKey: ['connections'] });
    },
    onError: (err: Error) => {
      setError(err.message);
    }
  });
};

export const useRejectConnectionRequest = () => {
  const queryClient = useQueryClient();
  const { setInfo, setError } = useUIStore();
  
  return useMutation({
    mutationFn: (requestId: string) => rejectConnectionRequest(requestId),
    onSuccess: () => {
      setInfo("Connection request declined successfully");
      queryClient.invalidateQueries({ queryKey: ['connectionRequests'] });
      queryClient.refetchQueries({ queryKey: ['connectionRequests'] });
    },
    onError: (err: Error) => {
      setError(err.message);
    }
  });
};