import { type Connection } from "../../connections/types/connections";
import { type User } from "../../../shared/types/user";

// Main message type matching server response
export interface Message {
  id: string;
  connection_id: string;
  sender_id: string;
  receiver_id: string;
  type: 'text' | 'media' | 'mixed';
  content?: string | null;
  media_url?: string | null;
  media_type?: string | null;
  is_read: boolean;
  created_at: string;
  read_at?: string | null;
  sender?: User | null;
  receiver?: User | null;
  connection?: Connection | null;
}

// Component props interfaces
export interface MessageBubbleProps {
  message: Message;
  isOwnMessage: boolean;
  showAvatar: boolean;
}

export interface MessageInputProps {
  onSendMessage: (content: string) => void;
  disabled?: boolean;
}


export interface MessageError {
    error: string;
    details: string;
}

export interface MessageResponse {
    message: string;
    data: Message;
}

export interface ConnMessagesResponse {
    messages: Message[];
    count: number;
    limit: number;
    offset: number;
}

export interface UnreadCountResponse {
    unread_count: number;
}

export interface MessageReadResponse {
    message: string;
}

export interface SendTextMessageBody {
  connection_id: string;
  content: string;
  sender_id?: string;
  receiver_id?: string;
}

export interface SendMediaMessageBody {
  connection_id: string;
  media: File;
  text?: string;
}

export interface ChatListItem {
  connection_id: string;
  other_user: User | null;
  last_message?: Message | null;
  unread_count: number;
  last_activity: string;
  connection_status: string;
}

export interface ChatList {
  chats: ChatListItem[];
  total_chats: number;
  unread_total: number;
}