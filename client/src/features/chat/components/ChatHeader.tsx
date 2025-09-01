import React, { useState, useRef, useEffect } from 'react';
import { Icon } from '@iconify/react/dist/iconify.js';
import { getInitials } from '../../../shared/utils/utils';
import { type ChatListItem } from '../types/chat';
import { useNavigate } from 'react-router';
import ConfirmModal from '../../../shared/components/ConfirmModal';
import { useDeleteConnection } from '../../connections/hooks/useConnections';

interface ChatHeaderProps {
  selectedChat: ChatListItem;
  isMobile?: boolean;
  setIsSidebarOpen: (isOpen: boolean) => void;
  onlineUsers: Set<string>;
  userStatuses: Map<string, string>;
}

const ChatHeader: React.FC<ChatHeaderProps> = ({
  selectedChat,
  isMobile = false,
  setIsSidebarOpen,
  onlineUsers,
  userStatuses,
}) => {
  const [showActionsMenu, setShowActionsMenu] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate()
  const deleteConnection = useDeleteConnection();

  const otherUser = selectedChat.other_user;
  const profilePhoto =
    otherUser?.profile_photo ||
    (otherUser?.photos && otherUser.photos.length > 0 ? otherUser.photos[0].photo_url : null);

  const initials = getInitials(otherUser?.first_name ?? "", otherUser?.last_name ?? "");


  const handleDelete = () => {
    setShowActionsMenu(false);
    setShowDeleteModal(true);
  };

  const confirmDelete = () => {
      deleteConnection.mutate((selectedChat.connection_id));
    setShowDeleteModal(false);
  };

  const getStatusText = (userId: string) => {
    const status = userStatuses.get(userId);
    if (status === 'away') return 'Away';
    return 'Offline';
  };

  // Close actions menu when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setShowActionsMenu(false);
      }
    };

    if (showActionsMenu) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [showActionsMenu]);
  
  if (!otherUser) return null;

  return (
    <div className="chat-header-main">
      <button className="mobile-back-btn" onClick={() => setIsSidebarOpen(true)}>
        <Icon icon={isMobile ? 'mdi:arrow-left' : 'mdi:menu'} />
      </button>

      <div className="chat-partner-info" onClick={() => navigate(`/users/${otherUser.id}`)}>
        <div className="chat-partner-avatar"
         style={
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
        </div>
        <div className="chat-partner-details">
          <div className="chat-partner-name">
            {otherUser.first_name} {otherUser.last_name}
          </div>
          <div className={`chat-partner-status ${!onlineUsers.has(otherUser.id)? "offline" : ""}`} >
            {onlineUsers.has(otherUser.id) ? 'Online now' : getStatusText(otherUser.id)}
          </div>
        </div>
      </div>

      <div className="chat-actions">
        <div className="chat-actions-dropdown" ref={dropdownRef}>
          <button
            className="chat-action-btn"
            title="More options"
            onClick={() => setShowActionsMenu(!showActionsMenu)}
          >
            <Icon icon="mdi:dots-horizontal" />
          </button>

          {showActionsMenu && (
            <div className="chat-actions-menu">
              <button
                className="chat-menu-item chat-menu-profile"
                onClick={() => {
                  navigate(`/users/${otherUser.id}`)
                  setShowActionsMenu(false);
                }}
              >
                <Icon icon="mdi:account-outline" />
                View Profile
              </button>
              <button
                className="chat-menu-item chat-menu-remove"
                onClick={handleDelete}
              >
                <Icon icon="mdi:delete-outline" />
                Remove
              </button>
            </div>
          )}
        </div>
      </div>

      <ConfirmModal
        isOpen={showDeleteModal}
        title="Remove Connection"
        message={`Are you sure you want to remove your connection with ${otherUser.first_name} ${otherUser.last_name}? This action cannot be undone.`}
        confirmText="Remove"
        cancelText="Cancel"
        onConfirm={confirmDelete}
        onCancel={() => setShowDeleteModal(false)}
        isLoading={deleteConnection.isPending}
        variant="danger"
      />
    </div>
  );
};

export default ChatHeader;
