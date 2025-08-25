import React, { useState } from 'react'; // Import useState
import { Icon } from '@iconify/react/dist/iconify.js';
import { getInitials } from '../../../shared/utils/utils';
import { type Message } from '../types/chat';

interface MessageBubbleProps {
  message: Message;
  isOwnMessage: boolean;
  showAvatar: boolean;
}

const MessageBubble: React.FC<MessageBubbleProps> = ({ message, isOwnMessage, showAvatar }) => {
  // State to manage the image modal
  const [isModalOpen, setIsModalOpen] = useState(false);

  const messageTime = new Date(message.created_at);
  const senderUser = message.sender;
  const avatarUser = senderUser;

  const handleImageClick = (e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent event bubbling
    setIsModalOpen(true);
  };

  const renderMessageContent = () => {
    switch (message.type) {
      case 'mixed':
        return (
          <>
            {message.media_url && (
              <div className="message-media">
                <img 
                  src={message.media_url} 
                  alt="Shared media" 
                  className="message-image" 
                  onClick={handleImageClick} // Add click handler
                />
              </div>
            )}
            {message.content && <div className="message-caption">{message.content}</div>}
          </>
        );
      case 'media':
        return (
          message.media_url && (
            <div className="message-media">
              <img 
                src={message.media_url} 
                alt="Shared media" 
                className="message-image" 
                onClick={handleImageClick} // Add click handler
              />
            </div>
          )
        );
      default:
        return message.content;
    }
  };

  return (
    <>
      {/* The Modal for displaying the expanded image */}
      {isModalOpen && (
        <div className="image-modal-overlay" onClick={() => setIsModalOpen(false)}>
          <button className="image-modal-close" onClick={() => setIsModalOpen(false)}>
            <Icon icon="mdi:close" />
          </button>
          <img 
            src={message.media_url!} 
            alt="Expanded media" 
            className="image-modal-content"
            onClick={(e) => e.stopPropagation()} // Prevent closing when clicking the image itself
          />
        </div>
      )}

      {/* The original message bubble component */}
      <div className={`message ${isOwnMessage ? 'sent' : 'received'}`}>
        {!isOwnMessage && showAvatar && avatarUser && (
          <div className="message-avatar">
            {getInitials(avatarUser.first_name, avatarUser.last_name)}
          </div>
        )}
        <div className="message-bubble">
          {renderMessageContent()}
          <div className="message-time">
            {messageTime.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
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
        </div>
        {isOwnMessage && showAvatar && avatarUser && (
          <div className="message-avatar">
            {getInitials(avatarUser.first_name, avatarUser.last_name)}
          </div>
        )}
      </div>
    </>
  );
};

export default MessageBubble;