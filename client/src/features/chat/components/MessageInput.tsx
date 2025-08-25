import React, { useRef, useEffect, useState } from 'react';
import { Icon } from '@iconify/react/dist/iconify.js';

interface MessageInputProps {
  messageInput: string;
  setMessageInput: (input: string, skipTypingEvent?: boolean) => void;
  handleSendMessage: () => void;
  handleSendMediaMessage?: (file: File) => void;
  handleInputKeyDown: (e: React.KeyboardEvent) => void;
  messageInputRef: React.RefObject<HTMLTextAreaElement | null>;
  onTyping?: (isTyping: boolean) => void;
}

const MessageInput: React.FC<MessageInputProps> = ({
  messageInput,
  setMessageInput,
  handleSendMessage,
  handleSendMediaMessage,
  handleInputKeyDown,
  messageInputRef,
  onTyping
}) => {
  const typingTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const isCurrentlyTypingRef = useRef(false);
  const skipNextTypingEventRef = useRef(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [filePreview, setFilePreview] = useState<string | null>(null);

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

  // Handle file selection
  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setSelectedFile(file);
      
      // Create preview for image (all files are images now)
      const reader = new FileReader();
      reader.onload = () => {
        setFilePreview(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  // Handle media attachment click
  const handleAttachClick = () => {
    fileInputRef.current?.click();
  };

  // Handle sending media
  const handleSendMedia = () => {
    if (selectedFile && handleSendMediaMessage) {
      handleSendMediaMessage(selectedFile);
      clearSelectedFile();
    }
  };

  // Clear selected file
  const clearSelectedFile = () => {
    setSelectedFile(null);
    setFilePreview(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  // Handle send button click (text or media)
  const handleSendClick = () => {
    if (selectedFile) {
      handleSendMedia();
    } else {
      handleSendMessage();
    }
  };

  // Check if we can send (either text or media)
  const canSend = messageInput.trim() || selectedFile;

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
      {/* Image Preview */}
      {selectedFile && filePreview && (
        <div className="file-preview">
          <div className="file-preview-content">
            <img src={filePreview} alt="Preview" className="file-preview-image" />
            <span className="file-name">{selectedFile.name}</span>
          </div>
          <button className="file-preview-remove" onClick={clearSelectedFile} title="Remove image">
            <Icon icon="mdi:close" />
          </button>
        </div>
      )}

      <div className="message-input-wrapper">
        <div className="input-actions">
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            onChange={handleFileSelect}
            style={{ display: 'none' }}
          />
          <button 
            className="input-action-btn" 
            title="Attach image" 
            onClick={handleAttachClick}
          >
            <Icon icon="mdi:camera" />
          </button>
          <button className="input-action-btn" title="Add emoji">
            <Icon icon="mdi:emoticon-happy-outline" />
          </button>
        </div>
        <textarea
          ref={messageInputRef}
          className="message-input"
          placeholder={selectedFile ? "Add a caption..." : "Type a message..."}
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
          disabled={!canSend}
          onClick={handleSendClick}
        >
          <Icon icon={selectedFile ? "mdi:send" : "mdi:send"} />
        </button>
      </div>
    </div>
  );
};

export default MessageInput;
