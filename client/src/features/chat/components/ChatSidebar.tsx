import React from 'react';
import { Icon } from '@iconify/react/dist/iconify.js';
import { type ChatListItem } from '../types/chat';
import ChatList from './ChatList';

interface ChatSidebarProps {
  searchTerm: string;
  setSearchTerm: (term: string) => void;
  filteredChats: ChatListItem[];
  setSelectedChat: (chat: ChatListItem | null) => void;
  selectedChat: ChatListItem | null;
  isLoading: boolean;
}

const ChatSidebar: React.FC<ChatSidebarProps> = ({
  searchTerm,
  setSearchTerm,
  filteredChats,
  setSelectedChat,
  selectedChat,
  isLoading
}) => (
  <div className="chat-sidebar" id="chatSidebar">
    <div className="sidebar-header">
      <div className="sidebar-title">Messages</div>
      <div className="search-box">
        <input
          type="text"
          className="search-input"
          placeholder="Search conversations..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
        <Icon icon="mdi:magnify" className="search-icon" />
      </div>
    </div>
    <ChatList
      filteredChats={filteredChats}
      setSelectedChat={setSelectedChat}
      selectedChat={selectedChat}
      isLoading={isLoading}
    />
  </div>
);

export default ChatSidebar;