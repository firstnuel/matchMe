import { create } from "zustand";
import { type User } from "../../../shared/types/user";

interface ProfileState {
  user: User | Partial<User> | null;
  setUser: (u: User | Partial<User> | null) => void;

}

export const useUserProfile = create<ProfileState>((set) => ({
  user: null,
  setUser: (u: User | Partial<User> | null) => set({ user: u }),
}));