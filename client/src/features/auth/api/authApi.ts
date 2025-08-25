import axios, { AxiosError } from "axios";
import { type RegisterData, type AuthData, type AuthError, type LoginData } from "../types/auth";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
const REGISTER_URL = `${API_BASE_URL}/auth/register`;
const LOGIN_URL = `${API_BASE_URL}/auth/login`;



export const registerUser = async (registerData: RegisterData): Promise<AuthData | AuthError> => {
    try {
        const { data } = await axios.post<AuthData | AuthError>(REGISTER_URL, registerData)
        return data;
    } catch (error) {
        if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as AuthError;
    }
    return {
        error: "Registration failed",
        details: "An unexpected error occurred. Please try again later."
        } as AuthError;
  }
}   


export const loginUser = async (loginData: LoginData): Promise<AuthData | AuthError> => {
    try {
        if (!loginData.email || !loginData.password) {
            return { error: "Invalid input", details: "Email and password are required" } as AuthError;
        }
        const { data } = await axios.post<AuthData | AuthError>(LOGIN_URL, loginData)
        return data;
    } catch (error) {
        if (error instanceof AxiosError && error.response?.data) {
      return error.response.data as AuthError;
    }
    return {
        error: "Login failed",
        details: "An unexpected error occurred. Please try again later."
        } as AuthError;
  }
}   