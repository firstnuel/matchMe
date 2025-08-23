import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AuthState {
  authToken: string | null;
  setAuthToken: (token: string) => void;
  clearAuth: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      authToken: null,
      setAuthToken: (token: string) => set({ authToken: token }),
      clearAuth: () => set({ authToken: null }),
    }),
    {
      name: 'auth-storage', // key in localStorage
    }
  )
);