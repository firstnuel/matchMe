/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect, useRef, useState, useCallback } from 'react';
import { WebSocketClient, EventType, } from '../services/websocket';
import { useAuthStore } from '../../features/auth/hooks/authStore';

export const useWebSocket = (endpoint: string) => {
  const { authToken } = useAuthStore();
  const [isConnected, setIsConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const clientRef = useRef<WebSocketClient | null>(null);

  useEffect(() => {
    if (!authToken) {
      setError('No authentication token available');
      return;
    }

    const client = new WebSocketClient(endpoint, authToken);
    clientRef.current = client;

    client.connect()
      .then(() => {
        setIsConnected(true);
        setError(null);
      })
      .catch((err) => {
        setError(err.message);
        setIsConnected(false);
      });

    return () => {
      client.disconnect();
      clientRef.current = null;
      setIsConnected(false);
    };
  }, [endpoint, authToken]);

  const addEventListener = useCallback((eventType: EventType, callback: (data: any) => void) => {
    if (clientRef.current) {
      clientRef.current.addEventListener(eventType, callback);
    }
  }, []);

  const removeEventListener = useCallback((eventType: EventType, callback: (data: any) => void) => {
    if (clientRef.current) {
      clientRef.current.removeEventListener(eventType, callback);
    }
  }, []);

  const sendMessage = useCallback((type: EventType, data?: any) => {
    if (clientRef.current) {
      clientRef.current.sendMessage(type, data);
    }
  }, []);

  return {
    isConnected,
    error,
    addEventListener,
    removeEventListener,
    sendMessage
  };
};

// Hook for general user status WebSocket
export const useStatusWebSocket = () => {
  return useWebSocket('/ws/status');
};

// Hook for chat WebSocket with specific connection
export const useChatWebSocket = (connectionId: string) => {
  return useWebSocket(`/ws/chat/${connectionId}`);
};

// Hook for typing indicators WebSocket with specific connection
export const useTypingWebSocket = (connectionId: string) => {
  return useWebSocket(`/ws/typing/${connectionId}`);
};