import { useEffect, useRef, useCallback } from 'react';
import { useMarkMessagesAsRead } from './useChatMessage'; 
import { type Message } from '../types/chat'; 

// Interface for the hook's props
interface UseAutoMarkAsReadProps {
  connectionId: string | undefined;
  messages: Message[];
  currentUserId: string | undefined;
  isActive: boolean;
  delay?: number;
}


export const useAutoMarkAsRead = ({
  connectionId,
  messages,
  currentUserId,
  isActive,
  delay = 500, 
}: UseAutoMarkAsReadProps) => {
  const markAsReadMutation = useMarkMessagesAsRead();
  const lastProcessedMessageRef = useRef<string | null>(null);
  const markAsReadTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const isPageVisibleRef = useRef<boolean>(!document.hidden);

  const hasUnreadMessages = useCallback(() => {
    // Efficiently checks if there are any unread messages from other users
    return messages.some(
      (message) => !message.is_read && message.sender_id !== currentUserId
    );
  }, [messages, currentUserId]);

  // Define cleanup logic once to keep the code DRY
  const clearMarkAsReadTimeout = useCallback(() => {
    if (markAsReadTimeoutRef.current) {
      clearTimeout(markAsReadTimeoutRef.current);
    }
  }, []);

  const scheduleMarkAsRead = useCallback(() => {
    if (!connectionId || !isPageVisibleRef.current || !isActive) return;

    // Clear any existing timeout to reset the debounce timer
    clearMarkAsReadTimeout();

    markAsReadTimeoutRef.current = setTimeout(() => {
      // Only mutate if there are unread messages and no mutation is already pending
      if (hasUnreadMessages() && !markAsReadMutation.isPending) {
        markAsReadMutation.mutate(connectionId);
      }
    }, delay); // Use the configurable delay
  }, [connectionId, isActive, delay, hasUnreadMessages, markAsReadMutation, clearMarkAsReadTimeout]);

  // Effect to track page visibility
  useEffect(() => {
    const handleVisibilityChange = () => {
      isPageVisibleRef.current = !document.hidden;
      // If page becomes visible and chat is active, check for unread messages
      if (isPageVisibleRef.current && isActive && hasUnreadMessages()) {
        scheduleMarkAsRead();
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);
    // Cleanup function to prevent memory leaks
    return () => {
      document.removeEventListener('visibilitychange', handleVisibilityChange);
      clearMarkAsReadTimeout();
    };
  }, [isActive, hasUnreadMessages, scheduleMarkAsRead, clearMarkAsReadTimeout]);

  // Effect to auto-mark messages as read when new ones arrive
  useEffect(() => {
    if (!isActive || !connectionId || !currentUserId || messages.length === 0) {
      return;
    }

    // Performance Improvement: Find the last message from another user without sorting.
    let latestMessage: Message | undefined;
    for (let i = messages.length - 1; i >= 0; i--) {
      const msg = messages[i];
      if (msg.sender_id !== currentUserId) {
        latestMessage = msg;
        break; // Found the latest one, exit the loop
      }
    }

    if (!latestMessage) return;

    // Trigger mark-as-read if the latest message is new, unread, and the page is visible.
    if (
      latestMessage.id !== lastProcessedMessageRef.current &&
      !latestMessage.is_read &&
      isPageVisibleRef.current
    ) {
      lastProcessedMessageRef.current = latestMessage.id;
      scheduleMarkAsRead();
    }

    // Correctness Improvement: Add cleanup for this effect to clear timeout on unmount
    return clearMarkAsReadTimeout;
    
  }, [messages, connectionId, currentUserId, isActive, scheduleMarkAsRead, clearMarkAsReadTimeout]);

  return {
    hasUnreadMessages: hasUnreadMessages(),
    isMarkingAsRead: markAsReadMutation.isPending,
  };
};