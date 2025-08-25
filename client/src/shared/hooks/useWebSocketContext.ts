import { useContext } from 'react';
import { WebSocketContext, type WebSocketContextType } from '../contexts/WebSocketContext';

export const useWebSocketContext = (): WebSocketContextType => {
  const context = useContext(WebSocketContext);
  if (context === undefined) {
    throw new Error('useWebSocketContext must be used within a WebSocketProvider');
  }
  return context;
};