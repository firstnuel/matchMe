import axios, { AxiosError } from "axios";
import { type UserResponse, type UserError, type PhotoUploadResponse } from "../types/user";
import { type UpdateUserRequest } from "../types/user";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://192.168.100.44:3000";
const CURRENT_USER_URL =  API_BASE_URL + "/api/me";
const UPLOAD_PHOTOS_URL = CURRENT_USER_URL + "/photos";

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


export const getLocationCity = async ({
  latitude,
  longitude,
}: {
  latitude: number;
  longitude: number;
}): Promise<{ lat: number; lon: number; city: string }> => {
  const { data } = await axios.get(
    `https://nominatim.openstreetmap.org/reverse?lat=${latitude}&lon=${longitude}&format=json`
  );

  const city =
    data.address.city + ", " +  data.address.country ||
    data.address.town ||
    data.address.village ||
    data.address.county ||
    "Unknown";

  return {
    lat: latitude,
    lon: longitude,
    city,
  };
};

export const getCurrentUser = async (): Promise<UserResponse | UserError> => {
  try {
    const { data } = await api.get<UserResponse>(
      CURRENT_USER_URL
    );
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as UserError;
    }
    return {
      error: "Failed to fetch user",
      details: "An unexpected error occurred. Please try again later.",
    } as UserError;
  }
};

export const updateUser = async (userData: Partial<UpdateUserRequest>): Promise<UserResponse | UserError> => {
  try {
    const { data } = await api.put<UserResponse>(
      CURRENT_USER_URL,
      userData
    );
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as UserError;
    }
    return {
      error: "Failed to update user",
      details: "An unexpected error occurred. Please try again later.",
    } as UserError;
  }
};

export const uploadUserPhotos = async (photos: File[]): Promise<PhotoUploadResponse | UserError> => {
  try {
    const formData = new FormData();
    
    // Add all photos to the form data with the key "photos"
    photos.forEach((photo) => {
      formData.append('photos', photo);
    });

    const { data } = await api.post<PhotoUploadResponse>(
      UPLOAD_PHOTOS_URL,
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      }
    );
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as UserError;
    }
    return {
      error: "Failed to upload photos",
      details: "An unexpected error occurred. Please try again later.",
    } as UserError;
  }
};

export const deleteUserPhoto = async (photoId: string): Promise<{ message: string; photo_id: string } | UserError> => {
  try {
    const { data } = await api.delete<{ message: string; photo_id: string }>(
      `${UPLOAD_PHOTOS_URL}/${photoId}`
    );
    return data;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as UserError;
    }
    return {
      error: "Failed to delete photo",
      details: "An unexpected error occurred. Please try again later.",
    } as UserError;
  }
};
