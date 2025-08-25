import React from 'react';
import { Icon } from '@iconify/react/dist/iconify.js';
import { getInitials } from '../../../shared/utils/utils';
import { type ChatListItem } from '../types/chat';
import { type User } from '../../../shared/types/user';
import MessageBubble from './MessageBubble';
import MessageInput from './MessageInput';
import { useConnectionMessages } from '../hooks/useChatMessage';
import { useWebSocketContext } from '../../../shared/hooks/useWebSocketContext';
import { useChatEffects } from '../hooks/useChatEffects';
import ChatHeader from './ChatHeader';

interface ChatMainProps {
  selectedChat: ChatListItem | null;
  isSidebarOpen: boolean;
  setIsSidebarOpen: (isOpen: boolean) => void;
  messagesEndRef: React.RefObject<HTMLDivElement | null>;
  messageInput: string;
  setMessageInput: (input: string) => void;
  handleSendMessage: () => void;
  handleSendMediaMessage?: (file: File, caption?: string) => void;
  handleInputKeyDown: (e: React.KeyboardEvent, selectedFile?: File) => void;
  messageInputRef: React.RefObject<HTMLTextAreaElement | null>;
  isMobile?: boolean;
  currentUser: User | null | undefined;
}

const ChatMain: React.FC<ChatMainProps> = ({
  selectedChat,
  setIsSidebarOpen,
  messagesEndRef,
  messageInput,
  setMessageInput,
  handleSendMessage,
  handleSendMediaMessage,
  handleInputKeyDown,
  messageInputRef,
  isMobile = false,
  currentUser,
}) => {
  const { onlineUsers, userStatuses } = useWebSocketContext();

  // Fetch server messages
  const { data: messagesData, isLoading: isLoadingMessages } = useConnectionMessages(
    selectedChat?.connection_id || '',
    50,
    0
  );
  const serverMessages = messagesData && 'messages' in messagesData ? messagesData.messages : [];

  // Use the chat effects hook
  const { messages, isOtherUserTyping, handleTyping } = useChatEffects({
    selectedChat,
    currentUser,
    messagesEndRef,
    serverMessages,
  });

  if (!selectedChat) {
    return (
      <div className="chat-main">
        <div className="empty-chat">
          <Icon icon="mdi:message-text-outline" className="empty-chat-icon" />
          <div className="empty-chat-title">Select a conversation</div>
          <div className="empty-chat-subtitle">Choose a chat to start messaging</div>
        </div>
      </div>
    );
  }

  const otherUser = selectedChat.other_user;
  if (!otherUser || !otherUser.first_name || !otherUser.last_name) return null;

  if (isLoadingMessages) {
    return (
      <div className="chat-main">
        <ChatHeader
          selectedChat={selectedChat}
          isMobile={isMobile}
          setIsSidebarOpen={setIsSidebarOpen}
          onlineUsers={onlineUsers}
          userStatuses={userStatuses}
        />
        <div className="messages-container">
          <div className="messages-loading">
            <div className="loading-spinner"></div>
            <p>Loading messages...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="chat-main">
      <ChatHeader
        selectedChat={selectedChat}
        isMobile={isMobile}
        setIsSidebarOpen={setIsSidebarOpen}
        onlineUsers={onlineUsers}
        userStatuses={userStatuses}
      />

      <div className="messages-container">
        {messages.length === 0 ? (
          <div className="messages-empty">
            <Icon icon="mdi:message-text-outline" className="empty-messages-icon" />
            <p>No messages yet</p>
            <p className="empty-messages-subtitle">Start a conversation with {otherUser.first_name}</p>
          </div>
        ) : (
          <>
            <div className="message-date">
              <span>Today</span>
            </div>
            <div className="message-group">
              {messages.map((message, index) => {
                const isOwnMessage = message.sender_id === currentUser?.id;
                const prevMessage = index > 0 ? messages[index - 1] : null;
                const showAvatar = !prevMessage || prevMessage.sender_id !== message.sender_id;
                return (
                  <MessageBubble
                    key={message.id}
                    message={message}
                    isOwnMessage={isOwnMessage}
                    showAvatar={showAvatar}
                  />
                );
              })}

              {isOtherUserTyping && (
                <div className="typing-indicator">
                  <div className="typing-avatar">
                    {getInitials(otherUser.first_name, otherUser.last_name)}
                  </div>
                  <div className="typing-bubble">
                    <div className="typing-dots">
                      <div className="typing-dot"></div>
                      <div className="typing-dot"></div>
                      <div className="typing-dot"></div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </>
        )}
        <div ref={messagesEndRef} />
      </div>

      <MessageInput
        messageInput={messageInput}
        setMessageInput={setMessageInput}
        handleSendMessage={handleSendMessage}
        handleSendMediaMessage={handleSendMediaMessage}
        handleInputKeyDown={handleInputKeyDown}
        messageInputRef={messageInputRef}
        onTyping={handleTyping}
      />
    </div>
  );
};

export default ChatMain;
