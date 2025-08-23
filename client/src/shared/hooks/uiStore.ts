import { create } from "zustand";

interface UIState {
  view: "chat" | "home" | "connections" | "profile";
  setView: (v: "chat" | "home" | "connections" | "profile") => void;
}

export const useUIStore = create<UIState>((set) => ({
  view: "home",
  setView: (v: "chat" | "home" | "connections" | "profile") => set({ view: v }),
}));
