import { create } from "zustand";
import { type User } from "../../../shared/types/user";

interface ProfileState {
  user: User | Partial<User> | null;
  pendingUpdate: Partial<User> | null; // Store the update data separately
  setUser: (u: User | Partial<User> | null) => void;
  setPendingUpdate: (data: Partial<User>) => void;
  clearPendingUpdate: () => void;
  hasPendingChanges: () => boolean;
}

export const useUserProfile = create<ProfileState>((set, get) => ({
  user: null,
  pendingUpdate: null,
  
  setUser: (u: User | Partial<User> | null) => set({ user: u }),
  
  setPendingUpdate: (data: Partial<User>) => set({ pendingUpdate: data }),
  
  clearPendingUpdate: () => set({ pendingUpdate: null }),
  
  // Helper function to check if there are pending changes
  hasPendingChanges: () => {
    const state = get();
    return state.pendingUpdate !== null && Object.keys(state.pendingUpdate).length > 0;
  },
}));