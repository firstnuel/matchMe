import React, { useRef, useEffect } from 'react';
import { Icon } from '@iconify/react/dist/iconify.js';

interface MessageInputProps {
  messageInput: string;
  setMessageInput: (input: string, skipTypingEvent?: boolean) => void;
  handleSendMessage: () => void;
  handleInputKeyDown: (e: React.KeyboardEvent) => void;
  messageInputRef: React.RefObject<HTMLTextAreaElement | null>;
  onTyping?: (isTyping: boolean) => void;
}

const MessageInput: React.FC<MessageInputProps> = ({
  messageInput,
  setMessageInput,
  handleSendMessage,
  handleInputKeyDown,
  messageInputRef,
  onTyping
}) => {
  const typingTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const isCurrentlyTypingRef = useRef(false);
  const skipNextTypingEventRef = useRef(false);

  const handleTypingEvent = (isTyping: boolean) => {
    if (!onTyping) return;

    // Clear any existing timeout first to prevent race conditions
    if (typingTimeoutRef.current) {
      clearTimeout(typingTimeoutRef.current);
      typingTimeoutRef.current = null;
    }

    if (isTyping && !isCurrentlyTypingRef.current) {
      // User started typing - emit immediately
      isCurrentlyTypingRef.current = true;
      onTyping(true);
    }

    if (isTyping) {
      // Set timeout to stop typing indicator after 2.5 seconds of no activity
      typingTimeoutRef.current = setTimeout(() => {
        if (isCurrentlyTypingRef.current) {
          isCurrentlyTypingRef.current = false;
          onTyping(false);
        }
        typingTimeoutRef.current = null;
      }, 2500);
    } else {
      // User stopped typing - emit immediately
      if (isCurrentlyTypingRef.current) {
        isCurrentlyTypingRef.current = false;
        onTyping(false);
      }
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newValue = e.target.value;
    setMessageInput(newValue);
    
    // Auto-resize textarea
    if (messageInputRef.current) {
      messageInputRef.current.style.height = 'auto';
      messageInputRef.current.style.height = `${messageInputRef.current.scrollHeight}px`;
    }

    // Skip typing events if flagged (after message sending)
    if (skipNextTypingEventRef.current) {
      skipNextTypingEventRef.current = false;
      return;
    }

    // Handle typing events with proper debouncing
    // Always indicate typing on input, even for spaces
    if (newValue.length > 0) {
      // User is actively typing (including spaces)
      handleTypingEvent(true);
    } else {
      // Input is completely empty (user deleted everything)
      handleTypingEvent(false);
    }
  };

  // Watch for external input clearing and set the skip flag
  useEffect(() => {
    if (messageInput === '' && isCurrentlyTypingRef.current) {
      // Input was cleared externally while typing was active
      skipNextTypingEventRef.current = true;
    }
  }, [messageInput]);

  const handleInputFocus = () => {
    // Don't emit typing on focus unless there's text (including spaces)
    if (messageInput.length > 0) {
      handleTypingEvent(true);
    }
  };

  const handleInputBlur = () => {
    // Only stop typing if input is actually empty, not just on focus loss
    if (messageInput.trim().length === 0) {
      handleTypingEvent(false);
    }
    // If there's text, let the natural timeout handle it
  };

  return (
    <div className="message-input-container">
      <div className="message-input-wrapper">
        <div className="input-actions">
          <button className="input-action-btn" title="Attach photo">
            <Icon icon="mdi:camera" />
          </button>
          <button className="input-action-btn" title="Add emoji">
            <Icon icon="mdi:emoticon-happy-outline" />
          </button>
        </div>
        <textarea
          ref={messageInputRef}
          className="message-input"
          placeholder="Type a message..."
          rows={1}
          value={messageInput}
          onInput={handleInputChange}
          onChange={handleInputChange}
          onFocus={handleInputFocus}
          onBlur={handleInputBlur}
          onKeyDown={handleInputKeyDown}
        />
        <button
          className="send-btn"
          disabled={!messageInput.trim()}
          onClick={handleSendMessage}
        >
          <Icon icon="mdi:send" />
        </button>
      </div>
    </div>
  );
};

export default MessageInput;
