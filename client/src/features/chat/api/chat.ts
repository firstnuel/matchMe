import axios, { AxiosError } from "axios";
import {
  type ChatList,
  type MessageError,
  type MessageResponse,
  type ConnMessagesResponse,
  type UnreadCountResponse,
  type MessageReadResponse,
  type SendTextMessageBody,
  type SendMediaMessageBody,
} from "../types/chat";

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || "http://192.168.100.44:3000";
const MESSAGES_URL = API_BASE_URL + "/messages";

// Create Axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
});


// Get chat list
export const getChatList = async (): Promise<ChatList | MessageError> => {
  try {
    const { data } = await api.get<ChatList>(`${MESSAGES_URL}/chat-list`);
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as MessageError;
    }
    return {
      error: "Failed to fetch chat list",
      details: "An unexpected error occurred. Please try again later.",
    };
  }
};

// Request interceptor (auth token)
api.interceptors.request.use(
  (config) => {
    const authData = localStorage.getItem("auth-storage");
    const parsed = authData ? JSON.parse(authData) : null;
    const authToken = parsed?.state?.authToken ?? null;

    if (authToken) {
      config.headers.Authorization = `Bearer ${authToken}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor → handle expired/invalid tokens
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("auth-storage");
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);


// Send text message
export const sendTextMessage = async (
  body: SendTextMessageBody
): Promise<MessageResponse | MessageError> => {
  try {
    console.log('🚀 Making API call to:', `${MESSAGES_URL}/text`);
    console.log('🚀 Request body:', body);
    const { data } = await api.post<MessageResponse>(`${MESSAGES_URL}/text`, body);
    console.log('✅ API Response:', data);
    return data;
  } catch (error) {
    console.error('❌ API call failed:', error);
    if (error instanceof AxiosError) {
      console.error('❌ Axios error details:', {
        status: error.response?.status,
        statusText: error.response?.statusText,
        data: error.response?.data,
        url: error.config?.url
      });
      if (error.response?.data) {
        return error.response.data as MessageError;
      }
    }
    return {
      error: "Failed to send text message",
      details: "An unexpected error occurred. Please try again later.",
    };
  }
};

// Send media message
export const sendMediaMessage = async (
  body: SendMediaMessageBody
): Promise<MessageResponse | MessageError> => {
  try {
    const formData = new FormData();
    formData.append("connection_id", body.connection_id);
    formData.append("media", body.media);

    const { data } = await api.post<MessageResponse>(
      `${MESSAGES_URL}/media`,
      formData,
      { headers: { "Content-Type": "multipart/form-data" } }
    );

    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as MessageError;
    }
    return {
      error: "Failed to send media message",
      details: "An unexpected error occurred. Please try again later.",
    };
  }
};


// Get messages for a connection
export const getConnectionMessages = async (
  connectionId: string,
  params?: { limit?: number; offset?: number }
): Promise<ConnMessagesResponse | MessageError> => {
  try {
    const { data } = await api.get<ConnMessagesResponse>(
      `${MESSAGES_URL}/connection/${connectionId}`,
      { params }
    );
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as MessageError;
    }
    return {
      error: "Failed to fetch messages",
      details: "An unexpected error occurred. Please try again later.",
    };
  }
};

// Mark messages as read
export const markMessagesAsRead = async (
  connectionId: string
): Promise<MessageReadResponse | MessageError> => {
  try {
    const { data } = await api.put<MessageReadResponse>(
      `${MESSAGES_URL}/connection/${connectionId}/read`
    );
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as MessageError;
    }
    return {
      error: "Failed to mark messages as read",
      details: "An unexpected error occurred. Please try again later.",
    };
  }
};

// Get unread count
export const getUnreadCount = async (): Promise<
  UnreadCountResponse | MessageError
> => {
  try {
    const { data } = await api.get<UnreadCountResponse>(
      `${MESSAGES_URL}/unread-count`
    );
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as MessageError;
    }
    return {
      error: "Failed to fetch unread count",
      details: "An unexpected error occurred. Please try again later.",
    };
  }
};
