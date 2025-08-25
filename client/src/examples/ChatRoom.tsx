/* eslint-disable @typescript-eslint/no-explicit-any */
import React, { useEffect, useState, useRef } from 'react';
import { useWebSocketContext } from '../shared/hooks/useWebSocketContext';
import { EventType, type MessageEvent, type TypingEvent } from '../shared/services/websocket';

interface ChatRoomProps {
  connectionId: string;
  currentUserId: string;
  otherUserId: string;
}

const ChatRoom: React.FC<ChatRoomProps> = ({ connectionId, currentUserId, otherUserId }) => {
  const { connectToChat, connectToTyping, disconnectFromChat, disconnectFromTyping, userStatuses } = useWebSocketContext();
  const [messages, setMessages] = useState<any[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const [isTyping, setIsTyping] = useState(false);
  const [otherUserTyping, setOtherUserTyping] = useState(false);
  const chatClientRef = useRef<any>(null);
  const typingClientRef = useRef<any>(null);
  const typingTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    // Connect to chat WebSocket
    const chatClient = connectToChat(connectionId);
    if (chatClient) {
      chatClientRef.current = chatClient;


      // Listen for new messages
      chatClient.addEventListener(EventType.MESSAGE_NEW, (data: MessageEvent) => {
        setMessages(prev => [...prev, data.message]);
      });

      // Listen for message read events
      chatClient.addEventListener(EventType.MESSAGE_READ, (data: any) => {
        setMessages(prev => prev.map(msg => 
          msg.id === data.message_id ? { ...msg, read: true } : msg
        ));
      });
    }

    // Connect to typing WebSocket
    const typingClient = connectToTyping(connectionId);
    if (typingClient) {
      typingClientRef.current = typingClient;

      // Listen for typing events
      typingClient.addEventListener(EventType.MESSAGE_TYPING, (data: TypingEvent) => {
        if (data.user_id !== currentUserId) {
          setOtherUserTyping(data.is_typing);
        }
      });
    }

    return () => {
      disconnectFromChat(connectionId);
      disconnectFromTyping(connectionId);
    };
  }, [connectionId, currentUserId, connectToChat, connectToTyping, disconnectFromChat, disconnectFromTyping]);

  const sendMessage = () => {
    if (!newMessage.trim() || !chatClientRef.current) return;

    const messageData = {
      content: newMessage.trim(),
      connection_id: connectionId,
      sender_id: currentUserId,
      receiver_id: otherUserId
    };

    chatClientRef.current.sendMessage(EventType.MESSAGE_NEW, messageData);
    setNewMessage('');
    stopTyping();
  };

  const handleTyping = (value: string) => {
    setNewMessage(value);

    if (!isTyping && value.trim()) {
      setIsTyping(true);
      typingClientRef.current?.sendMessage(EventType.MESSAGE_TYPING, {
        connection_id: connectionId,
        user_id: currentUserId,
        is_typing: true
      });
    }

    // Clear existing timeout
    if (typingTimeoutRef.current) {
      clearTimeout(typingTimeoutRef.current);
    }

    // Set timeout to stop typing after 2 seconds of inactivity
    typingTimeoutRef.current = setTimeout(() => {
      stopTyping();
    }, 2000);
  };

  const stopTyping = () => {
    if (isTyping) {
      setIsTyping(false);
      typingClientRef.current?.sendMessage(EventType.MESSAGE_TYPING, {
        connection_id: connectionId,
        user_id: currentUserId,
        is_typing: false
      });
    }
  };

  const otherUserStatus = userStatuses.get(otherUserId) || 'offline';

  return (
    <div style={{ 
      border: '1px solid #ddd', 
      borderRadius: '8px', 
      padding: '16px', 
      maxWidth: '500px',
      margin: '20px 0'
    }}>
      <div style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center',
        borderBottom: '1px solid #eee',
        paddingBottom: '8px',
        marginBottom: '16px'
      }}>
        <h3 style={{ margin: 0 }}>Chat Room</h3>
        <div style={{ 
          fontSize: '12px',
          padding: '4px 8px',
          borderRadius: '12px',
          backgroundColor: otherUserStatus === 'online' ? '#dcfce7' : '#fee2e2',
          color: otherUserStatus === 'online' ? '#16a34a' : '#dc2626'
        }}>
          {otherUserStatus}
        </div>
      </div>

      <div style={{ 
        height: '200px', 
        overflowY: 'auto', 
        border: '1px solid #eee',
        borderRadius: '4px',
        padding: '8px',
        marginBottom: '16px'
      }}>
        {messages.map((message, index) => (
          <div 
            key={index} 
            style={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: message.sender_id === currentUserId ? 'flex-end' : 'flex-start',
              marginBottom: '8px'
            }}
          >
            <div style={{
              backgroundColor: message.sender_id === currentUserId ? '#7c3aed' : '#f3f4f6',
              color: message.sender_id === currentUserId ? 'white' : '#374151',
              padding: '8px 12px',
              borderRadius: '12px',
              maxWidth: '70%'
            }}>
              {message.content}
            </div>
            <div style={{ fontSize: '10px', color: '#9ca3af', marginTop: '2px' }}>
              {new Date(message.created_at).toLocaleTimeString()}
              {message.read && message.sender_id === currentUserId && (
                <span style={{ marginLeft: '4px', color: '#10b981' }}>✓✓</span>
              )}
            </div>
          </div>
        ))}
        
        {otherUserTyping && (
          <div style={{ 
            fontStyle: 'italic', 
            color: '#6b7280', 
            fontSize: '14px',
            padding: '4px 0'
          }}>
            Other user is typing...
          </div>
        )}
      </div>

      <div style={{ display: 'flex', gap: '8px' }}>
        <input
          type="text"
          value={newMessage}
          onChange={(e) => handleTyping(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
          placeholder="Type a message..."
          style={{
            flex: 1,
            padding: '8px 12px',
            border: '1px solid #d1d5db',
            borderRadius: '6px',
            outline: 'none'
          }}
        />
        <button 
          onClick={sendMessage} 
          disabled={!newMessage.trim()}
          style={{
            padding: '8px 16px',
            backgroundColor: !newMessage.trim() ? '#9ca3af' : '#7c3aed',
            color: 'white',
            border: 'none',
            borderRadius: '6px',
            cursor: !newMessage.trim() ? 'not-allowed' : 'pointer'
          }}
        >
          Send
        </button>
      </div>
    </div>
  );
};

export default ChatRoom;