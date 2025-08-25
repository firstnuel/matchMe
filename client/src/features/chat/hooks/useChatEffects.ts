/* eslint-disable react-hooks/exhaustive-deps */
import { useState, useEffect, useRef } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { useWebSocketContext } from '../../../shared/hooks/useWebSocketContext';
import { EventType, type MessageEvent, type TypingEvent } from '../../../shared/services/websocket';
import { type ChatListItem, type Message } from '../types/chat';
import { type User } from '../../../shared/types/user';
import { useMarkMessagesAsRead } from '../hooks/useChatMessage';

interface UseChatEffectsProps {
  selectedChat: ChatListItem | null;
  currentUser: User | null | undefined;
  messagesEndRef: React.RefObject<HTMLDivElement | null>;
  serverMessages: Message[];
}

interface UseChatEffectsReturn {
  messages: Message[];
  isOtherUserTyping: boolean;
  handleTyping: (isTyping: boolean) => void;
}

export const useChatEffects = ({
  selectedChat,
  currentUser,
  messagesEndRef,
  serverMessages,
}: UseChatEffectsProps): UseChatEffectsReturn => {
  const [localMessages, setLocalMessages] = useState<Message[]>([]);
  const [isOtherUserTyping, setIsOtherUserTyping] = useState(false);
  const markedAsReadRef = useRef<Set<string>>(new Set());
  const typingTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const sendTypingTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const queryClient = useQueryClient();
  const markAsReadMutation = useMarkMessagesAsRead();
  const { connectToChat, disconnectFromChat, connectToTyping, disconnectFromTyping } = useWebSocketContext();

  // Combine and sort messages
  const allMessages = [
    ...serverMessages,
    ...localMessages.filter((localMsg) => !serverMessages.some((serverMsg) => serverMsg.id === localMsg.id)),
  ];
  const messages = allMessages.sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime());

  // Mark messages as read when chat is opened
  useEffect(() => {
    const connectionId = selectedChat?.connection_id;
    if (
      connectionId &&
      selectedChat.unread_count > 0 &&
      !markedAsReadRef.current.has(connectionId) &&
      !markAsReadMutation.isPending
    ) {
      markedAsReadRef.current.add(connectionId);
      markAsReadMutation.mutate(connectionId);
    }

    return () => {
      if (connectionId) {
        markedAsReadRef.current.delete(connectionId);
      }
    };
  }, [selectedChat?.connection_id, selectedChat?.unread_count]);

  // WebSocket connection for real-time messages
  useEffect(() => {
    const connectionId = selectedChat?.connection_id;
    if (!connectionId) {
      setLocalMessages([]);
      return;
    }

    console.log('🔌 Connecting to chat WebSocket for:', connectionId);
    const chatClient = connectToChat(connectionId);
    if (!chatClient) {
      console.error('❌ Failed to get chat client for:', connectionId);
      return;
    }
    console.log('✅ Chat client connected:', chatClient.isConnected);

    const handleNewMessage = (messageData: MessageEvent) => {
      console.log('📨 Received new message:', messageData);
      if (messageData.connection_id === connectionId) {
        console.log('✅ Adding new message to local state');
        setLocalMessages((prev) => {
          if (prev.some((msg) => msg.id === messageData.message.id)) {
            console.log('⚠️ Duplicate message, ignoring');
            return prev;
          }
          return [...prev, messageData.message];
        });
        queryClient.invalidateQueries({ queryKey: ['connectionMessages', connectionId] });
        queryClient.invalidateQueries({ queryKey: ['chatList'] });
      } else {
        console.log('⚠️ Message not for this connection:', messageData.connection_id, 'vs', connectionId);
      }
    };

    // Add a generic event listener to catch ALL events
    console.log('🎯 Setting up MESSAGE_NEW event listener for connection:', connectionId);
    chatClient.addEventListener(EventType.MESSAGE_NEW, handleNewMessage);
    setLocalMessages([]);

    return () => {
      disconnectFromChat(connectionId);
    };
  }, [selectedChat?.connection_id, connectToChat, disconnectFromChat, queryClient]);

  // WebSocket connection for typing indicators
  useEffect(() => {
    const connectionId = selectedChat?.connection_id;
    if (!connectionId) {
      setIsOtherUserTyping(false);
      return;
    }

    console.log('⌨️ Connecting to typing WebSocket for:', connectionId);
    const typingClient = connectToTyping(connectionId);
    if (!typingClient) {
      console.error('❌ Failed to get typing client for:', connectionId);
      return;
    }
    console.log('✅ Typing client connected:', typingClient.isConnected);

    const handleTypingEvent = (typingData: TypingEvent) => {
      console.log('⌨️ Received typing event:', typingData);
      if (typingData.connection_id === connectionId && typingData.user_id !== currentUser?.id) {
        console.log('✅ Setting typing indicator:', typingData.is_typing);
        setIsOtherUserTyping(typingData.is_typing);
        if (typingData.is_typing) {
          if (typingTimeoutRef.current) {
            clearTimeout(typingTimeoutRef.current);
          }
          typingTimeoutRef.current = setTimeout(() => {
            setIsOtherUserTyping(false);
          }, 3000);
        } else {
          if (typingTimeoutRef.current) {
            clearTimeout(typingTimeoutRef.current);
          }
        }
      }
    };

    typingClient.addEventListener(EventType.MESSAGE_TYPING, handleTypingEvent);
    setIsOtherUserTyping(false);

    return () => {
      disconnectFromTyping(connectionId);
      if (typingTimeoutRef.current) {
        clearTimeout(typingTimeoutRef.current);
      }
      if (sendTypingTimeoutRef.current) {
        clearTimeout(sendTypingTimeoutRef.current);
      }
    };
  }, [selectedChat?.connection_id, connectToTyping, disconnectFromTyping]);

  // Auto-scroll to bottom
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [serverMessages.length, localMessages.length, isOtherUserTyping, messagesEndRef]);

  // Handle typing events with debouncing
  const handleTyping = (isTyping: boolean) => {
    const connectionId = selectedChat?.connection_id;
    if (!connectionId || !currentUser?.id) return;

    console.log('⌨️ Sending typing event:', { connectionId, isTyping, userId: currentUser.id });
    const typingClient = connectToTyping(connectionId);
    if (!typingClient) {
      console.error('❌ No typing client available');
      return;
    }

    // Clear any existing server-side timeout first
    if (sendTypingTimeoutRef.current) {
      clearTimeout(sendTypingTimeoutRef.current);
      sendTypingTimeoutRef.current = null;
    }

    if (isTyping) {
      // Send typing=true immediately
      typingClient.sendMessage(EventType.MESSAGE_TYPING, {
        connection_id: connectionId,
        user_id: currentUser.id,
        is_typing: true,
        updated_at: new Date().toISOString(),
      });
      
      // Set server-side backup timeout (slightly longer than client-side)
      sendTypingTimeoutRef.current = setTimeout(() => {
        typingClient.sendMessage(EventType.MESSAGE_TYPING, {
          connection_id: connectionId,
          user_id: currentUser.id,
          is_typing: false,
          updated_at: new Date().toISOString(),
        });
        sendTypingTimeoutRef.current = null;
      }, 3000); // 3 seconds - longer than client-side 2.5s
    } else {
      // Send typing=false immediately
      typingClient.sendMessage(EventType.MESSAGE_TYPING, {
        connection_id: connectionId,
        user_id: currentUser.id,
        is_typing: false,
        updated_at: new Date().toISOString(),
      });
    }
  };

  return { messages, isOtherUserTyping, handleTyping };
};