import React from 'react';
import { Icon } from '@iconify/react/dist/iconify.js';
import { getInitials } from '../../../shared/utils/utils';
import { type Message } from '../types/chat';

interface MessageBubbleProps {
  message: Message;
  isOwnMessage: boolean;
  showAvatar: boolean;
}

const MessageBubble: React.FC<MessageBubbleProps> = ({ message, isOwnMessage, showAvatar }) => {
  const messageTime = new Date(message.created_at);
  const senderUser = message.sender;
  const receiverUser = message.receiver;
  
  // Determine which user's avatar to show
  const avatarUser = isOwnMessage ? senderUser : receiverUser;
  
  return (
    <div className={`message ${isOwnMessage ? 'sent' : 'received'}`}>
      {!isOwnMessage && showAvatar && avatarUser && (
        <div className="message-avatar">
          {getInitials(avatarUser.first_name, avatarUser.last_name)}
        </div>
      )}
      <div className="message-bubble">
        {message.type === 'media' && message.media_url ? (
          <div className="message-media">
            <img src={message.media_url} alt="Shared media" className="message-image" />
          </div>
        ) : (
          message.content
        )}
        <div className="message-time">
          {messageTime.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
        </div>
        {isOwnMessage && (
          <div className="message-status">
            {message.is_read ? (
              <Icon icon="mdi:check-all" style={{ color: '#22c55e' }} />
            ) : (
              <Icon icon="mdi:check" />
            )}
          </div>
        )}
      </div>
      {isOwnMessage && showAvatar && avatarUser && (
        <div className="message-avatar">
          {getInitials(avatarUser.first_name, avatarUser.last_name)}
        </div>
      )}
    </div>
  );
};

export default MessageBubble;