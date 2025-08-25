/* eslint-disable @typescript-eslint/no-explicit-any */
export const EventType = {
  // Message events
  MESSAGE_NEW: 'message_new',
  MESSAGE_READ: 'message_read',
  MESSAGE_TYPING: 'message_typing',

  // User status events
  USER_ONLINE: 'user_online',
  USER_OFFLINE: 'user_offline',
  USER_AWAY: 'user_away',

  // Connection events
  CONNECTION_REQUEST: 'connection_request',
  CONNECTION_ACCEPTED: 'connection_accepted',
  CONNECTION_DROPPED: 'connection_dropped',

  // System events
  ERROR: 'error',
  PING: 'ping',
  PONG: 'pong'
} as const;

export type EventType = typeof EventType[keyof typeof EventType];

export interface WebSocketMessage {
  type: EventType;
  data?: any;
  timestamp: string;
  message_id: string;
}

export interface MessageEvent {
  message: any;
  connection_id: string;
  sender_id: string;
  receiver_id: string;
}

export interface TypingEvent {
  connection_id: string;
  user_id: string;
  is_typing: boolean;
  updated_at: string;
}

export interface UserStatusEvent {
  user_id: string;
  status: 'online' | 'offline' | 'away';
  last_activity: string;
}

export interface ConnectionRequestEvent {
  request: any;
  action: 'new' | 'accepted' | 'declined';
}

export interface ConnectionEvent {
  connection: any;
  action: 'established' | 'dropped';
}

export interface ErrorEvent {
  code: number;
  message: string;
}

export class WebSocketClient {
  private ws: WebSocket | null = null;
  private url: string;
  private authToken: string;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectInterval = 1000;
  private isConnecting = false;
  private eventListeners = new Map<EventType, ((data: any) => void)[]>();
  private messageQueue: WebSocketMessage[] = [];

  constructor(endpoint: string, authToken: string) {
    const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://192.168.100.44:3000';
    // Convert http to ws protocol
    const wsBaseUrl = baseUrl.replace(/^http/, 'ws');
    this.url = `${wsBaseUrl}${endpoint}`;
    this.authToken = authToken;
  }

  connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        resolve();
        return;
      }

      if (this.isConnecting) {
        reject(new Error('Connection already in progress'));
        return;
      }

      this.isConnecting = true;

      try {
        // Add auth token as query parameter
        const urlWithAuth = `${this.url}?token=${encodeURIComponent(this.authToken)}`;
        this.ws = new WebSocket(urlWithAuth);

        this.ws.onopen = () => {
          console.log('WebSocket connected to:', this.url);
          this.isConnecting = false;
          this.reconnectAttempts = 0;
          
          // Send any queued messages
          this.flushMessageQueue();
          
          resolve();
        };

        this.ws.onmessage = (event) => {
          console.log('🔌 Raw WebSocket message received:', event.data);
          try {
            const message: WebSocketMessage = JSON.parse(event.data);
            console.log('📨 Parsed WebSocket message:', message);
            this.handleMessage(message);
          } catch (error) {
            console.error('Failed to parse WebSocket message:', error, 'Raw data:', event.data);
          }
        };

        this.ws.onclose = (event) => {
          console.log('WebSocket disconnected:', event.code, event.reason);
          this.isConnecting = false;
          this.handleReconnect();
        };

        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error);
          this.isConnecting = false;
          reject(error);
        };

      } catch (error) {
        this.isConnecting = false;
        reject(error);
      }
    });
  }

  private handleMessage(message: WebSocketMessage) {
    console.log('🎯 Handling WebSocket message type:', message.type);
    const listeners = this.eventListeners.get(message.type);
    console.log('👂 Found', listeners?.length || 0, 'listeners for', message.type);
    if (listeners) {
      listeners.forEach((listener, index) => {
        console.log(`📢 Calling listener ${index + 1}/${listeners.length} for ${message.type}`);
        listener(message.data);
      });
    } else {
      console.log('⚠️ No listeners registered for message type:', message.type);
    }
  }

  private handleReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);
      
      setTimeout(() => {
        this.connect().catch(error => {
          console.error('Reconnection failed:', error);
        });
      }, this.reconnectInterval * this.reconnectAttempts);
    } else {
      console.error('Max reconnection attempts reached');
    }
  }

  addEventListener(eventType: EventType, callback: (data: any) => void) {
    if (!this.eventListeners.has(eventType)) {
      this.eventListeners.set(eventType, []);
    }
    this.eventListeners.get(eventType)!.push(callback);
  }

  removeEventListener(eventType: EventType, callback: (data: any) => void) {
    const listeners = this.eventListeners.get(eventType);
    if (listeners) {
      const index = listeners.indexOf(callback);
      if (index > -1) {
        listeners.splice(index, 1);
      }
    }
  }

  private flushMessageQueue() {
    while (this.messageQueue.length > 0 && this.ws?.readyState === WebSocket.OPEN) {
      const message = this.messageQueue.shift()!;
      this.ws.send(JSON.stringify(message));
    }
  }

  sendMessage(type: EventType, data?: any) {
    const message: WebSocketMessage = {
      type,
      data,
      timestamp: new Date().toISOString(),
      message_id: Math.random().toString(36).substring(2, 15)
    };

    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      // Queue the message to be sent when connection is established
      console.log('WebSocket not ready, queueing message:', type);
      this.messageQueue.push(message);
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.eventListeners.clear();
    this.messageQueue = [];
  }

  get isConnected() {
    return this.ws?.readyState === WebSocket.OPEN;
  }
}