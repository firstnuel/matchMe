/* eslint-disable react-hooks/exhaustive-deps */
import React, { createContext, useEffect, useState, type ReactNode, useCallback, useMemo } from 'react';
import { WebSocketClient, EventType, type UserStatusEvent, type UserStatusInitialEvent, type MessageEvent } from '../services/websocket';
import { useAuthStore } from '../../features/auth/hooks/authStore';
import { useCurrentUser } from '../../features/userProfile/hooks/useCurrentUser';
import { useQueryClient } from '@tanstack/react-query';

export interface WebSocketContextType {
  statusClient: WebSocketClient | null;
  isStatusConnected: boolean;
  onlineUsers: Set<string>;
  userStatuses: Map<string, 'online' | 'offline' | 'away'>;
  connectToChat: (connectionId: string) => WebSocketClient | null;
  connectToTyping: (connectionId: string) => WebSocketClient | null;
  disconnectFromChat: (connectionId: string) => void;
  disconnectFromTyping: (connectionId: string) => void;
}

// eslint-disable-next-line react-refresh/only-export-components
export const WebSocketContext = createContext<WebSocketContextType | undefined>(undefined);

interface WebSocketProviderProps {
  children: ReactNode;
}
export const WebSocketProvider: React.FC<WebSocketProviderProps> = ({ children }) => {
  const { authToken } = useAuthStore();
  const { data: currentUserData } = useCurrentUser();
  const queryClient = useQueryClient();
  const [statusClient, setStatusClient] = useState<WebSocketClient | null>(null);
  const [isStatusConnected, setIsStatusConnected] = useState(false);
  const [onlineUsers, setOnlineUsers] = useState<Set<string>>(new Set());
  const [userStatuses, setUserStatuses] = useState<Map<string, 'online' | 'offline' | 'away'>>(new Map());
  const [chatClients] = useState<Map<string, WebSocketClient>>(new Map());
  const [typingClients] = useState<Map<string, WebSocketClient>>(new Map());

  // Initialize status WebSocket connection
  useEffect(() => {
    if (!authToken || !currentUserData) {
      return;
    }

    // Disconnect existing client if any to avoid multiple connections
    if (statusClient) {
      statusClient.disconnect();
      setStatusClient(null);
      setIsStatusConnected(false);
    }

    const client = new WebSocketClient('/ws/status', authToken);
    
    // Set up event listeners
    client.addEventListener(EventType.USER_ONLINE, (data: UserStatusEvent) => {
      setOnlineUsers(prev => new Set([...prev, data.user_id]));
      setUserStatuses(prev => new Map([...prev, [data.user_id, 'online']]));
    });

    client.addEventListener(EventType.USER_OFFLINE, (data: UserStatusEvent) => {
      setOnlineUsers(prev => {
        const newSet = new Set([...prev]);
        newSet.delete(data.user_id);
        return newSet;
      });
      setUserStatuses(prev => new Map([...prev, [data.user_id, 'offline']]));
    });

    client.addEventListener(EventType.USER_AWAY, (data: UserStatusEvent) => {
      setUserStatuses(prev => new Map([...prev, [data.user_id, 'away']]));
    });

    // Handle initial status snapshot when client first connects
    client.addEventListener(EventType.USER_STATUS_INITIAL, (data: UserStatusInitialEvent) => {
      console.log('📊 Received initial user statuses:', data);
      
      // Process the array of online users
      const onlineUserIds = new Set<string>();
      const statusMap = new Map<string, 'online' | 'offline' | 'away'>();
      
      data.forEach((userStatus: UserStatusEvent) => {
        if (userStatus.status === 'online') {
          onlineUserIds.add(userStatus.user_id);
        }
        statusMap.set(userStatus.user_id, userStatus.status);
      });
      
      // Update state with all online users at once
      setOnlineUsers(onlineUserIds);
      setUserStatuses(statusMap);
    });

    client.addEventListener(EventType.CONNECTION_REQUEST, () => {
    });

    client.addEventListener(EventType.CONNECTION_ACCEPTED, () => {
    });

    // Handle new message notifications from StatusHub for background chats
    client.addEventListener(EventType.MESSAGE_NEW, (messageData: MessageEvent) => {
      // Invalidate chat list to refresh unread counts and last messages
      // This ensures the chat list updates even when the specific chat isn't open
      queryClient.invalidateQueries({ queryKey: ['chatList'] });
      queryClient.invalidateQueries({ queryKey: ['unreadCount'] });
      
      // Also invalidate the specific chat's messages if it exists in cache
      // This helps keep currently open chats in sync even if they're receiving via StatusHub
      if (messageData.connection_id) {
        queryClient.invalidateQueries({ 
          queryKey: ['connectionMessages', messageData.connection_id],
          exact: false
        });
      }
    });

    client.connect()
      .then(() => {
        setStatusClient(client);
        setIsStatusConnected(true);
      })
      .catch((error) => {
        console.error('Failed to connect to status WebSocket:', error);
        setIsStatusConnected(false);
      });

    return () => {
      client.disconnect();
      setStatusClient(null);
      setIsStatusConnected(false);
    };
  }, [authToken, currentUserData]);

  const connectToChat = useCallback((connectionId: string): WebSocketClient | null => {
    if (!authToken) {
      console.error('🚫 No auth token available for chat WebSocket');
      return null;
    }
    
    if (chatClients.has(connectionId)) {
      return chatClients.get(connectionId)!;
    }

    const client = new WebSocketClient(`/ws/chat/${connectionId}`, authToken);
    chatClients.set(connectionId, client);
    
    client.connect().then(() => {
    }).catch((error) => {
      console.error(`❌ Failed to connect to chat WebSocket for connection ${connectionId}:`, error);
      chatClients.delete(connectionId);
    });

    return client;
  }, [authToken]);

  const connectToTyping = useCallback((connectionId: string): WebSocketClient | null => {
    if (!authToken) {
      console.error('🚫 No auth token available for typing WebSocket');
      return null;
    }
    
    if (typingClients.has(connectionId)) {
      return typingClients.get(connectionId)!;
    }

    const client = new WebSocketClient(`/ws/typing/${connectionId}`, authToken);
    typingClients.set(connectionId, client);
    
    client.connect().then(() => {
    }).catch((error) => {
      console.error(`❌ Failed to connect to typing WebSocket for connection ${connectionId}:`, error);
      typingClients.delete(connectionId);
    });

    return client;
  }, [authToken]);

  const disconnectFromChat = useCallback((connectionId: string) => {
    const client = chatClients.get(connectionId);
    if (client) {
      client.disconnect();
      chatClients.delete(connectionId);
    }
  }, []);

  const disconnectFromTyping = useCallback((connectionId: string) => {
    const client = typingClients.get(connectionId);
    if (client) {
      client.disconnect();
      typingClients.delete(connectionId);
    }
  }, []);

  // Cleanup all connections when component unmounts or auth changes
  useEffect(() => {
    return () => {
      // Disconnect all chat clients
      chatClients.forEach((client) => {
        client.disconnect();
      });
      chatClients.clear();

      // Disconnect all typing clients
      typingClients.forEach((client) => {
        client.disconnect();
      });
      typingClients.clear();

      // Disconnect status client
      if (statusClient) {
        statusClient.disconnect();
      }
    };
  }, [authToken]);

  const contextValue = useMemo(() => ({
    statusClient,
    isStatusConnected,
    onlineUsers,
    userStatuses,
    connectToChat,
    connectToTyping,
    disconnectFromChat,
    disconnectFromTyping
  }), [
      statusClient, 
      isStatusConnected, 
      onlineUsers, 
      userStatuses, 
      connectToChat, 
      connectToTyping, 
      disconnectFromChat, 
      disconnectFromTyping
  ]);

  return (
    <WebSocketContext.Provider value={contextValue}>
      {children}
    </WebSocketContext.Provider>
  );
};