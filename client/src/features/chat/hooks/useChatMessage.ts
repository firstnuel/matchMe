/* eslint-disable @typescript-eslint/no-explicit-any */
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  sendTextMessage,
  sendMediaMessage,
  getConnectionMessages,
  markMessagesAsRead,
  getUnreadCount,
  getChatList
} from './../api/chat';
import {
  type SendTextMessageBody,
  type SendMediaMessageBody,
} from "../types/chat";
import { useAuthStore } from "../../auth/hooks/authStore";
import { useUIStore } from "../../../shared/hooks/uiStore";

// Fetch chat list
export const useChatList = () => {
  const { authToken } = useAuthStore();
  return useQuery({
    queryKey: ["chatList"],
    queryFn: getChatList,
    enabled: !!authToken,
    retry: false,
  });
};

// Fetch messages in a connection (paginated)
export const useConnectionMessages = (connectionId: string, limit?: number, offset?: number) => {
  const { authToken } = useAuthStore();
  return useQuery({
    queryKey: ["connectionMessages", connectionId, limit, offset],
    queryFn: () => getConnectionMessages(connectionId, { limit, offset }),
    enabled: !!authToken && !!connectionId,
    retry: false,
  });
};

// Fetch unread count
export const useUnreadCount = () => {
  const { authToken } = useAuthStore();
  return useQuery({
    queryKey: ["unreadCount"],
    queryFn: getUnreadCount,
    enabled: !!authToken,
    retry: false,
  });
};

// Send text message with optimistic update
export const useSendTextMessage = () => {
  const queryClient = useQueryClient();
  const { setInfo, setError } = useUIStore();

  return useMutation({
    mutationFn: (body: SendTextMessageBody) => {
      console.log('📤 Sending message to server:', body);
      return sendTextMessage(body);
    },
    onMutate: async (newMessage) => {
      // Cancel queries for this connection - need to match the exact query key format
      await queryClient.cancelQueries({ 
        queryKey: ["connectionMessages", newMessage.connection_id],
        exact: false // This will cancel all variants with different limit/offset
      });

      // Get the current data from the specific query being used (with limit=50, offset=0)
      const queryKey = ["connectionMessages", newMessage.connection_id, 50, 0];
      const previousMessages = queryClient.getQueryData(queryKey);

      // Create optimistic message with proper structure
      const optimisticMessage = {
        id: `temp-${Date.now()}`,
        connection_id: newMessage.connection_id,
        content: newMessage.content,
        type: 'text' as const,
        sender_id: newMessage.sender_id || '',
        receiver_id: newMessage.receiver_id || '',
        is_read: false,
        created_at: new Date().toISOString(),
        sending: true
      };

      console.log('🚀 Adding optimistic message to query key:', queryKey);
      console.log('🚀 Optimistic message:', optimisticMessage);

      // Update the specific query that ChatMain is reading from
      queryClient.setQueryData(queryKey, (old: any) => {
        const currentData = old || { messages: [] };
        const existingMessages = currentData.messages || [];
        
        console.log('📊 Existing messages count:', existingMessages.length);
        const newData = {
          ...currentData,
          messages: [...existingMessages, optimisticMessage]
        };
        console.log('📊 New messages count:', newData.messages.length);
        return newData;
      });

      return { previousMessages, queryKey };
    },
    onError: (err: any, newMessage, context: any) => {
      console.error('❌ Message sending failed:', err);
      console.error('❌ Error details:', {
        message: err.message,
        response: err.response?.data,
        status: err.response?.status
      });
      
      if (context?.previousMessages && context?.queryKey) {
        queryClient.setQueryData(context.queryKey, context.previousMessages);
      }
      setError(err.message || 'Failed to send message');
    },
    onSuccess: (response, variables) => {
      setInfo("Message sent successfully");
      
      console.log('✅ Message sent successfully, server response:', response);
      
      // Small delay before invalidation to let optimistic update settle
      setTimeout(() => {
        // Invalidate all connectionMessages queries for this connection (with any limit/offset)
        queryClient.invalidateQueries({ 
          queryKey: ["connectionMessages", variables.connection_id],
          exact: false // This will invalidate all variants
        });
        queryClient.invalidateQueries({ queryKey: ["chatList"] });
        queryClient.invalidateQueries({ queryKey: ["unreadCount"] });
        
        console.log('✅ Invalidated queries for connection:', variables.connection_id);
      }, 100);
    },
  });
};

// Send media message with optimistic update
export const useSendMediaMessage = () => {
  const queryClient = useQueryClient();
  const { setInfo, setError } = useUIStore();

  return useMutation({
    mutationFn: (body: SendMediaMessageBody) => sendMediaMessage(body),
    onMutate: async (newMessage) => {
      await queryClient.cancelQueries({ queryKey: ["connectionMessages", newMessage.connection_id] });

      const previousMessages = queryClient.getQueryData(["connectionMessages", newMessage.connection_id]);

      queryClient.setQueryData(["connectionMessages", newMessage.connection_id], (old: any = []) => [
        ...old,
        { ...newMessage, id: `temp-${Date.now()}`, sending: true },
      ]);

      return { previousMessages };
    },
    onError: (err: Error, newMessage, context: any) => {
      if (context?.previousMessages) {
        queryClient.setQueryData(["connectionMessages", newMessage.connection_id], context.previousMessages);
      }
      setError(err.message);
    },
    onSuccess: (_, variables) => {
      setInfo("Media message sent successfully");
      // Invalidate all connectionMessages queries for this connection (with any limit/offset)
      queryClient.invalidateQueries({ 
        queryKey: ["connectionMessages", variables.connection_id],
        exact: false // This will invalidate all variants
      });
      queryClient.invalidateQueries({ queryKey: ["chatList"] });
      queryClient.invalidateQueries({ queryKey: ["unreadCount"] });
      
      console.log('✅ Invalidated queries for connection:', variables.connection_id);
    },
  });
};

// Mark messages as read
export const useMarkMessagesAsRead = () => {
  const queryClient = useQueryClient();
  const { setInfo, setError } = useUIStore();

  return useMutation({
    mutationFn: (connectionId: string) => markMessagesAsRead(connectionId),
    onSuccess: (_, connectionId) => {
      setInfo("Messages marked as read");
      queryClient.invalidateQueries({ queryKey: ["connectionMessages", connectionId] });
      queryClient.invalidateQueries({ queryKey: ["chatList"] });
      queryClient.invalidateQueries({ queryKey: ["unreadCount"] });
    },
    onError: (err: Error) => {
      setError(err.message);
    },
  });
};
