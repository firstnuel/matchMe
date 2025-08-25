import React, { useState, useRef, useEffect } from 'react';
import { type ChatListItem } from '../types/chat';
import ChatSidebar from './ChatSidebar';
import ChatMain from './ChatMain';
import { useIsMobile } from '../hooks/useIsMobile';
import { useUIStore } from '../../../shared/hooks/uiStore';
import { useChatList, useSendTextMessage, useSendMediaMessage } from '../hooks/useChatMessage';
import { useCurrentUser } from '../../userProfile/hooks/useCurrentUser';
import '../styles.css';
import type { User } from '../../../shared/types/user';
import type { UserResponse } from '../../userProfile/types/user';

const ChatInterface = () => {
  const isMobile = useIsMobile();
  const { setIsChatMessageViewActive } = useUIStore();
  const [selectedChat, setSelectedChat] = useState<ChatListItem | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [messageInput, setMessageInput] = useState('');
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [mobileView, setMobileView] = useState<'sidebar' | 'chat'>('sidebar');
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const messageInputRef = useRef<HTMLTextAreaElement>(null);
  
  const { data: currentUser } = useCurrentUser();
  const { data: chatListData, isLoading: isLoadingChats, error: chatListError } = useChatList();
  const sendTextMutation = useSendTextMessage();
  const sendMediaMutation = useSendMediaMessage();

  useEffect(() => {
    if (!isMobile && chatListData && 'chats' in chatListData && chatListData.chats.length > 0 && !selectedChat) {
      setSelectedChat(chatListData.chats[0]);
    }
  }, [chatListData, isMobile, selectedChat]);

  useEffect(() => {
    if (isMobile) {
      setIsChatMessageViewActive(mobileView === 'chat');
    } else {
      setIsChatMessageViewActive(false);
    }
    return () => {
      setIsChatMessageViewActive(false);
    };
  }, [isMobile, mobileView, setIsChatMessageViewActive]);

  const handleSendMessage = () => {
    if (!messageInput.trim() || !selectedChat || !currentUser) return;
    sendTextMutation.mutate({
      connection_id: selectedChat.connection_id,
      content: messageInput.trim(),
      sender_id: (currentUser as UserResponse)?.user?.id || '',
      receiver_id: selectedChat.other_user?.id || ''
    });
    setMessageInput('');
  };

  const handleSendMediaMessage = (file: File, caption?: string) => {
    if (!selectedChat || !currentUser) return;
    sendMediaMutation.mutate({
      connection_id: selectedChat.connection_id,
      media: file,
      text: caption
    });
  };

  const handleInputKeyDown = (e: React.KeyboardEvent, selectedFile?: File) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      // Use the same logic as the send button
      if (selectedFile) {
        handleSendMediaMessage(selectedFile, messageInput);
        setMessageInput('');
      } else {
        handleSendMessage();
      }
    }
  };

  const filteredChats = (chatListData && 'chats' in chatListData ? chatListData.chats : []).filter((chat: ChatListItem) => {
    if (!chat.other_user) return false;
    const name = `${chat.other_user.first_name} ${chat.other_user.last_name}`.toLowerCase();
    const lastMessage = chat.last_message?.content?.toLowerCase() || '';
    return name.includes(searchTerm.toLowerCase()) || lastMessage.includes(searchTerm.toLowerCase());
  });

  const handleChatSelect = (chat: ChatListItem | null) => {
    setSelectedChat(chat);
    if (isMobile) {
      setMobileView('chat');
      setIsSidebarOpen(false);
    }
  };

  const handleBackToSidebar = () => {
    if (isMobile) {
      setMobileView('sidebar');
      setSelectedChat(null);
    } else {
      setIsSidebarOpen(true);
    }
  };

  // ✅ NEW: Handle error state gracefully without unmounting
  if (chatListError) {
    return (
      <div className="chat-container">
        <div className="chat-error">
          <p>Failed to load chats. Please try again.</p>
          <button onClick={() => window.location.reload()}>Retry</button>
        </div>
      </div>
    );
  }

  // ✅ MODIFIED: Always render the main structure and pass loading state as a prop
  return (
    <div className="chat-container">
      {(!isMobile || mobileView === 'sidebar') && (
        <ChatSidebar
          searchTerm={searchTerm}
          setSearchTerm={setSearchTerm}
          filteredChats={filteredChats}
          setSelectedChat={handleChatSelect}
          selectedChat={selectedChat}
          isLoading={isLoadingChats} // Pass loading state as a prop
        />
      )}
      {isSidebarOpen && !isMobile && (
        <div 
          className="chat-overlay show" 
          onClick={() => setIsSidebarOpen(false)}
        />
      )}
      {(!isMobile || mobileView === 'chat') && (
        <ChatMain
          selectedChat={selectedChat}
          isSidebarOpen={isSidebarOpen}
          setIsSidebarOpen={handleBackToSidebar}
          messagesEndRef={messagesEndRef}
          messageInput={messageInput}
          setMessageInput={setMessageInput}
          handleSendMessage={handleSendMessage}
          handleSendMediaMessage={handleSendMediaMessage}
          handleInputKeyDown={handleInputKeyDown}
          messageInputRef={messageInputRef}
          isMobile={isMobile}
          currentUser={currentUser && 'user' in currentUser ? currentUser.user as User: null}
        />
      )}
    </div>
  );
};

export default ChatInterface;