import React, { useState } from 'react';
import { WebSocketProvider } from '../shared/contexts/WebSocketContext';
import { useWebSocketContext } from '../shared/hooks/useWebSocketContext';
import ChatRoom from './ChatRoom';

// Example component showing how to use WebSocket features
const WebSocketDemo: React.FC = () => {
  const { isStatusConnected, onlineUsers, userStatuses } = useWebSocketContext();
  const [isConnectedToChat, setIsConnectedToChat] = useState(false);
  const [connectionId, setConnectionId] = useState('example-connection-123');
  const [currentUserId, setCurrentUserId] = useState('user-1');
  const [otherUserId, setOtherUserId] = useState('user-2');

  const handleConnect = () => {
    if (!connectionId.trim() || !currentUserId.trim() || !otherUserId.trim()) {
      alert('Please fill in all connection details');
      return;
    }
    setIsConnectedToChat(true);
  };

  const handleDisconnect = () => {
    setIsConnectedToChat(false);
  };

  return (
    <div style={{ padding: '20px' }}>
      <h2>WebSocket Connection Status</h2>
      <p>Status Connection: {isStatusConnected ? '🟢 Connected' : '🔴 Disconnected'}</p>
      
      <h3>Online Users ({onlineUsers.size})</h3>
      {onlineUsers.size === 0 ? (
        <p style={{ color: '#6b7280', fontStyle: 'italic' }}>No users currently online</p>
      ) : (
        <ul>
          {Array.from(onlineUsers).map(userId => (
            <li key={userId} style={{ marginBottom: '4px' }}>
              <span style={{ fontFamily: 'monospace' }}>{userId}</span> - 
              <span style={{ 
                marginLeft: '8px',
                padding: '2px 6px',
                borderRadius: '4px',
                fontSize: '12px',
                backgroundColor: userStatuses.get(userId) === 'online' ? '#dcfce7' : '#fee2e2',
                color: userStatuses.get(userId) === 'online' ? '#16a34a' : '#dc2626'
              }}>
                {userStatuses.get(userId) || 'unknown'}
              </span>
            </li>
          ))}
        </ul>
      )}

      <div style={{ marginTop: '24px', padding: '16px', backgroundColor: '#f9fafb', borderRadius: '8px' }}>
        <h3>Chat Demo Configuration</h3>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px', marginBottom: '16px' }}>
          <div>
            <label style={{ display: 'block', marginBottom: '4px', fontSize: '14px', fontWeight: '500' }}>
              Connection ID:
            </label>
            <input
              type="text"
              value={connectionId}
              onChange={(e) => setConnectionId(e.target.value)}
              style={{
                width: '100%',
                padding: '8px',
                border: '1px solid #d1d5db',
                borderRadius: '4px',
                fontSize: '14px'
              }}
            />
          </div>
          <div>
            <label style={{ display: 'block', marginBottom: '4px', fontSize: '14px', fontWeight: '500' }}>
              Current User ID:
            </label>
            <input
              type="text"
              value={currentUserId}
              onChange={(e) => setCurrentUserId(e.target.value)}
              style={{
                width: '100%',
                padding: '8px',
                border: '1px solid #d1d5db',
                borderRadius: '4px',
                fontSize: '14px'
              }}
            />
          </div>
        </div>
        <div style={{ marginBottom: '16px' }}>
          <label style={{ display: 'block', marginBottom: '4px', fontSize: '14px', fontWeight: '500' }}>
            Other User ID:
          </label>
          <input
            type="text"
            value={otherUserId}
            onChange={(e) => setOtherUserId(e.target.value)}
            style={{
              width: '100%',
              padding: '8px',
              border: '1px solid #d1d5db',
              borderRadius: '4px',
              fontSize: '14px'
            }}
          />
        </div>
        
        <div style={{ display: 'flex', gap: '12px' }}>
          <button
            onClick={isConnectedToChat ? handleDisconnect : handleConnect}
            style={{
              padding: '10px 20px',
              backgroundColor: isConnectedToChat ? '#dc2626' : '#16a34a',
              color: 'white',
              border: 'none',
              borderRadius: '6px',
              cursor: 'pointer',
              fontWeight: '500'
            }}
          >
            {isConnectedToChat ? 'Disconnect' : 'Connect'}
          </button>
          
          {isConnectedToChat && (
            <div style={{
              padding: '10px 16px',
              backgroundColor: '#dcfce7',
              color: '#16a34a',
              borderRadius: '6px',
              fontWeight: '500',
              fontSize: '14px'
            }}>
              🟢 Connected to chat room
            </div>
          )}
        </div>
      </div>

      {isConnectedToChat && (
        <div>
          <h3>Chat Room</h3>
          <p style={{ color: '#6b7280', fontSize: '14px', marginBottom: '16px' }}>
            Connected to room: <strong>{connectionId}</strong> | You: <strong>{currentUserId}</strong> | Other: <strong>{otherUserId}</strong>
          </p>
          <ChatRoom 
            connectionId={connectionId}
            currentUserId={currentUserId}
            otherUserId={otherUserId}
          />
        </div>
      )}

      <div style={{ marginTop: '32px', padding: '16px', backgroundColor: '#fef3c7', borderRadius: '8px' }}>
        <h4 style={{ marginTop: 0, color: '#92400e' }}>💡 Integration Notes</h4>
        <ul style={{ color: '#92400e', fontSize: '14px', lineHeight: '1.6' }}>
          <li>Wrap your main App component with <code>WebSocketProvider</code></li>
          <li>The status connection automatically connects when a user is authenticated</li>
          <li>Chat connections are created on-demand for specific conversation IDs</li>
          <li>All WebSocket connections require valid JWT authentication</li>
          <li>Real-time events include: messages, typing indicators, user status, and connection requests</li>
        </ul>
      </div>
    </div>
  );
};

// Main app wrapper showing how to set up WebSocket provider
const AppWithWebSocket: React.FC = () => {
  return (
    <WebSocketProvider>
      <div className="app">
        <WebSocketDemo />
        {/* Your other app components */}
      </div>
    </WebSocketProvider>
  );
};

export default AppWithWebSocket;