/* eslint-disable @typescript-eslint/no-explicit-any */
import { useWebSocketContext } from '../../../shared/hooks/useWebSocketContext';
import { EventType, type ConnectionRequestEvent, type ConnectionEvent } from '../../../shared/services/websocket';
import {  useQueryClient } from '@tanstack/react-query';
import { useUIStore } from '../../../shared/hooks/uiStore';
import { useEffect } from 'react';

// Hook to handle incoming connection request WebSocket events
export const useConnectionRequestWebSocket = () => {
  const { statusClient } = useWebSocketContext();
  const queryClient = useQueryClient();
  const { setInfo } = useUIStore();

  useEffect(() => {
    if (!statusClient) return;

    const handleConnectionRequest = (data: ConnectionRequestEvent) => {
      if (data.action === 'new') {
        // New connection request received
        setInfo(`New connection request from ${data.request.sender?.first_name || 'someone'}`);
        
        // Update the connection requests query cache
        queryClient.invalidateQueries({ queryKey: ['connectionRequests'] });
        queryClient.refetchQueries({ queryKey: ['connectionRequests'] });
        
        // Optionally update the cache directly for instant UI updates
        queryClient.setQueryData(['connectionRequests'], (oldData: any) => {
          if (!oldData) return oldData;
          return {
            ...oldData,
            requests: [data.request, ...oldData.requests],
            count: oldData.count + 1
          };
        });
      } else if (data.action === 'accepted') {
        // Connection request was accepted
        setInfo('Your connection request was accepted!');
        
        // Remove from connection requests and trigger refetch
        queryClient.invalidateQueries({ queryKey: ['connectionRequests'] });
        queryClient.refetchQueries({ queryKey: ['connectionRequests'] });
      } else if (data.action === 'declined') {
        // Connection request was declined
        setInfo('Your connection request was declined');
        
        // Remove from connection requests
        queryClient.invalidateQueries({ queryKey: ['connectionRequests'] });
        queryClient.refetchQueries({ queryKey: ['connectionRequests'] });
      }
    };

    statusClient.addEventListener(EventType.CONNECTION_REQUEST, handleConnectionRequest);

    return () => {
      statusClient.removeEventListener(EventType.CONNECTION_REQUEST, handleConnectionRequest);
    };
  }, [statusClient, queryClient, setInfo]);
};

// Hook to handle connection accepted WebSocket events  
export const useConnectionAcceptedWebSocket = () => {
  const { statusClient } = useWebSocketContext();
  const queryClient = useQueryClient();
  const { setInfo } = useUIStore();

  useEffect(() => {
    if (!statusClient) return;

    const handleConnectionAccepted = (data: ConnectionEvent) => {
      if (data.action === 'established') {
        // New connection established
        setInfo(`Connection established with ${data.connection.user_a?.first_name || data.connection.user_b?.first_name || 'user'}`);
        
        // Update both connections and connection requests caches
        queryClient.invalidateQueries({ queryKey: ['connections'] });
        queryClient.invalidateQueries({ queryKey: ['connectionRequests'] });
        queryClient.refetchQueries({ queryKey: ['connections'] });
        queryClient.refetchQueries({ queryKey: ['connectionRequests'] });

        // Optionally update the connections cache directly for instant UI updates
        queryClient.setQueryData(['connections'], (oldData: any) => {
          if (!oldData) return oldData;
          return {
            ...oldData,
            connections: [data.connection, ...oldData.connections],
            count: oldData.count + 1
          };
        });
      } else if (data.action === 'dropped') {
        // Connection was dropped
        setInfo('A connection was dropped');
        
        // Remove from connections
        queryClient.invalidateQueries({ queryKey: ['connections'] });
        queryClient.refetchQueries({ queryKey: ['connections'] });
      }
    };

    statusClient.addEventListener(EventType.CONNECTION_ACCEPTED, handleConnectionAccepted);

    return () => {
      statusClient.removeEventListener(EventType.CONNECTION_ACCEPTED, handleConnectionAccepted);
    };
  }, [statusClient, queryClient, setInfo]);
};

// Combined hook to handle all connection-related WebSocket events
export const useConnectionWebSocketEvents = () => {
  useConnectionRequestWebSocket();
  useConnectionAcceptedWebSocket();
};