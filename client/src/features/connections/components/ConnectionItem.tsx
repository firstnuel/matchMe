import React, { useState } from 'react';
import { useAcceptConnectionRequest, useRejectConnectionRequest, useDeleteConnection } from '../hooks/useConnections';
import { type User } from '../../../shared/types/user';
import ConfirmModal from '../../../shared/components/ConfirmModal';
import { useNavigate } from 'react-router';


interface ConnectionItemProps {
  id: string;
  type: 'request' | 'connection';
  user?: User;
  initials: string;
  name: string;
  description: string;
  createdAt: string;
}

const ConnectionItem: React.FC<ConnectionItemProps> = ({ 
  id, 
  type, 
  user, 
  initials, 
  name, 
  description, 
  createdAt 
}) => {
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const navigate = useNavigate()
  const acceptRequest = useAcceptConnectionRequest();
  const rejectRequest = useRejectConnectionRequest();
  const deleteConnection = useDeleteConnection();

  const handleAccept = () => {
    acceptRequest.mutate(id);
  };

  const handleReject = () => {
    rejectRequest.mutate(id);
  };

  const handleDelete = () => {
    setShowDeleteModal(true);
  };

  const confirmDelete = () => {
    deleteConnection.mutate(id);
    setShowDeleteModal(false);
  };

  const cancelDelete = () => {
    setShowDeleteModal(false);
  };

  const displayPhoto = user?.profile_photo || user?.photos?.[0]?.photo_url;

  return (
    <>
      <div className="connection-item">
        <div className="connection-avatar" onClick={() => navigate(`/users/${user?.id}`)}>
          {displayPhoto ? (
            <img 
              src={displayPhoto} 
              alt={`${name} profile`}
              className="avatar-image"
            />
          ) : (
            initials
          )}
        </div>
        
        <div className="connection-info">
          <div className="connection-name" onClick={() => navigate(`/users/${user?.id}`)}>{name}</div>
          <div className="connection-description">{description}</div>
          <div className="connection-time">
            {new Date(createdAt).toLocaleDateString('en-US', {
              month: 'short',
              day: 'numeric',
              year: new Date().getFullYear() !== new Date(createdAt).getFullYear() ? 'numeric' : undefined
            })}
          </div>
        </div>
        
        <div className="connection-actions">
          {type === 'request' ? (
            <>
              <button 
                className="accept-btn"
                onClick={handleAccept}
                disabled={acceptRequest.isPending}
              >
                {acceptRequest.isPending ? 'Accepting...' : 'Accept'}
              </button>
              <button 
                className="decline-btn"
                onClick={handleReject}
                disabled={rejectRequest.isPending}
              >
                {rejectRequest.isPending ? 'Declining...' : 'Decline'}
              </button>
            </>
          ) : (
            <button 
              className="delete-btn"
              onClick={handleDelete}
              disabled={deleteConnection.isPending}
              title="Remove connection"
            >
              {deleteConnection.isPending ? 'Removing...' : 'Remove'}
            </button>
          )}
        </div>
      </div>

      <ConfirmModal
        isOpen={showDeleteModal}
        title="Remove Connection"
        message={`Are you sure you want to remove your connection with ${name}? This action cannot be undone.`}
        confirmText="Remove"
        cancelText="Cancel"
        onConfirm={confirmDelete}
        onCancel={cancelDelete}
        isLoading={deleteConnection.isPending}
        variant="danger"
      />
    </>
  );
};

export default ConnectionItem;