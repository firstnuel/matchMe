import axios, { AxiosError } from "axios";
import { type ConnectionError,
     type Connections,
     type ConnectionRequests,
     type ConnectionRequestResponse,
     type SendConnectionRequestBody,
     type ConnectionAcceptResponse,
     type ConnectionNoContentResponse,
     type SkipConnectionRequestBody
    } from "../types/connections";


const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
const CONNECTIONS_URL = API_BASE_URL + "/connections"
const REQUEST_URL = API_BASE_URL + "/connection-requests"

// Create Axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
});

// Request interceptors
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
      // Token expired or invalid → clear and redirect
      localStorage.removeItem("auth-storage");
      window.location.href = "/login"; // or navigate programmatically
    }
    return Promise.reject(error);
  }
);

export const getConnections = async (): Promise<Connections | ConnectionError> => {
  try {
    const { data } = await api.get<Connections>(
      `${CONNECTIONS_URL}/`, { params: { mode: "full" } },
    );
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as ConnectionError;
    }
    return {
      error: "Failed to fetch user Connections",
      details: "An unexpected error occurred. Please try again later.",
    } as ConnectionError;
  }
};

export const deleteConnections = async (id: string): Promise<ConnectionNoContentResponse | ConnectionError> => {
  try {
    const { data } = await api.delete<ConnectionNoContentResponse>(
      `${CONNECTIONS_URL}/${id}`
    );
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as ConnectionError;
    }
    return {
      error: "Failed to delete connection",
      details: "An unexpected error occurred. Please try again later.",
    } as ConnectionError;
  }
};

export const getConnectionRequests = async (): Promise<ConnectionRequests | ConnectionError> => {
  try {
    const { data } = await api.get<ConnectionRequests>(
       `${REQUEST_URL}/`, 
    );
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as ConnectionError;
    }
    return {
      error: "Failed to fetch user Connections requests",
      details: "An unexpected error occurred. Please try again later.",
    } as ConnectionError;
  }
};

export const sendConnectionRequest = async (reqData: SendConnectionRequestBody): Promise<ConnectionRequestResponse | ConnectionError> => {
  try {
    const { data } = await api.post<ConnectionRequestResponse>(
      `${REQUEST_URL}/`, 
      reqData,
    );
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as ConnectionError;
    }
    return {
      error: "Failed to fetch send Connection request",
      details: "An unexpected error occurred. Please try again later.",
    } as ConnectionError;
  }
};

export const skipConnectionRequest = async (reqData: SkipConnectionRequestBody): Promise<ConnectionRequestResponse | ConnectionError> => {
  try {
    const { data } = await api.post<ConnectionRequestResponse>(
      `${REQUEST_URL}/skip`, 
      reqData,
    );
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as ConnectionError;
    }
    return {
      error: "Failed to fetch skip Connection request",
      details: "An unexpected error occurred. Please try again later.",
    } as ConnectionError;
  }
};

export const acceptConnectionRequest = async (id: string): Promise<ConnectionAcceptResponse | ConnectionError> => {
  try {
    const { data } = await api.put<ConnectionAcceptResponse>(
      `${REQUEST_URL}/${id}/accept`, 
    );
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as ConnectionError;
    }
    return {
      error: "Failed to fetch send Connection request",
      details: "An unexpected error occurred. Please try again later.",
    } as ConnectionError;
  }
};

export const rejectConnectionRequest = async (id: string): Promise<ConnectionNoContentResponse | ConnectionError> => {
  try {
    const { data } = await api.put<ConnectionNoContentResponse>(
      `${REQUEST_URL}/${id}/decline`, 
    );
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as ConnectionError;
    }
    return {
      error: "Failed to fetch decline Connection request",
      details: "An unexpected error occurred. Please try again later.",
    } as ConnectionError;
  }
};