import React from 'react';
import { getInitials, formatTime } from '../../../shared/utils/utils';
import { type ChatListItem } from '../types/chat';
import { useWebSocketContext } from '../../../shared/hooks/useWebSocketContext';

interface ChatListProps {
  filteredChats: ChatListItem[];
  setSelectedChat: (chat: ChatListItem | null) => void;
  selectedChat: ChatListItem | null;
  isLoading: boolean;
}

const ChatList: React.FC<ChatListProps> = ({
  filteredChats,
  setSelectedChat,
  selectedChat,
  isLoading
}) => {
  const { onlineUsers } = useWebSocketContext();
  
  if (isLoading) {
    return (
      <div className="chat-list">
        {[...Array(5)].map((_, index) => (
          <div key={index} className="chat-item-skeleton">
            <div className="chat-avatar-skeleton"></div>
            <div className="chat-info-skeleton">
              <div className="chat-name-skeleton"></div>
              <div className="chat-message-skeleton"></div>
            </div>
          </div>
        ))}
      </div>
    );
  }

  if (filteredChats.length === 0) {
    return (
      <div className="chat-list-empty">
        <p>No conversations yet</p>
      </div>
    );
  }

  return (
    <div className="chat-list">
      {filteredChats.map((chat) => {
        const otherUser = chat.other_user;
        if (!otherUser || !otherUser.first_name || !otherUser.last_name) return null;

        const profilePhoto =
          otherUser?.profile_photo ||
          (otherUser?.photos && otherUser.photos.length > 0 ? otherUser.photos[0].photo_url : null);
        
        const initials = getInitials(otherUser?.first_name ?? "", otherUser?.last_name ?? "");

        return (
          <div
            key={chat.connection_id}
            className={`chat-item ${selectedChat?.connection_id === chat.connection_id ? 'active' : ''} ${chat.unread_count > 0 ? 'unread' : ''}`}
            onClick={() => setSelectedChat(chat)}
          >
            <div className="chat-avatar" style={
              profilePhoto
                ? {
                    backgroundImage: `url(${profilePhoto})`,
                    backgroundSize: "cover",
                    backgroundPosition: "center",
                    color: "transparent",
                  }
                : {}
            }>
            {profilePhoto ? "" : initials}
              {onlineUsers.has(otherUser.id) && <div className="online-indicator"></div>}
            </div>
            <div className="chat-info">
              <div className="chat-header">
                <div className="chat-name">{otherUser.first_name} {otherUser.last_name}</div>
                <div className="chat-time">{formatTime(new Date(chat.last_activity))}</div>
              </div>
              <div className="chat-preview">
                <div className="chat-message-preview">
                  {chat.last_message?.content || 'No messages yet'}
                </div>
                {chat.unread_count > 0 && (
                  <div className="unread-badge">{chat.unread_count}</div>
                )}
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default ChatList;