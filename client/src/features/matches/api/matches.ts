import axios, { AxiosError } from "axios";
import { type RecommendationsResponse, type MatchError, type DistanceResponse } from './../types/match';
import { type UserResponse } from "../../userProfile/types/user";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
const USER_RECOMMENDATIONS_URL =  API_BASE_URL + "/api/me/recommendations";
const USER_BASE_URL = API_BASE_URL + "/users"
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

export const getRecommendations = async (): Promise<RecommendationsResponse | MatchError> => {
  try {
    const { data } = await api.get<RecommendationsResponse>(
      USER_RECOMMENDATIONS_URL
    );
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as MatchError;
    }
    return {
      error: "Failed to fetch user recommendations",
      details: "An unexpected error occurred. Please try again later.",
    } as MatchError;
  }
};

export const getUserProfile = async (id: string): Promise<UserResponse | MatchError> => {
  try {
    const { data } = await api.get<UserResponse>(
      `${USER_BASE_URL}/${id}/profile`
    );
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as MatchError;
    }
    return {
      error: "Failed to fetch user profile",
      details: "An unexpected error occurred. Please try again later.",
    } as MatchError;
  }
};

export const getDistanceFromUser = async (id: string): Promise<DistanceResponse | MatchError> => {
  try {
    const { data } = await api.get<DistanceResponse>(
      `${USER_BASE_URL}/${id}/distance`
    );
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as MatchError;
    }
    return {
      error: "Failed to fetch distance from user",
      details: "An unexpected error occurred. Please try again later.",
    } as MatchError;
  }
};


