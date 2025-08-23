import { create } from "zustand";

interface UIState {
  view: "chat" | "home" | "connections" | "profile" | "edit-profile";
  setView: (v: "chat" | "home" | "connections" | "profile" | "edit-profile") => void;
  errorMsg: string;
  infoMsg: string;
  clearMsgs: () => void;
  setInfo: (info: string) => void;
  setError: (err: string) => void;
}

export const useUIStore = create<UIState>((set) => ({
  view: "home",
  setView: (v: "chat" | "home" | "connections" | "profile" | "edit-profile") => set({ view: v }),
  errorMsg: "",
  infoMsg: "",
  clearMsgs: () => set({ infoMsg: "", errorMsg: "" }),
  setInfo: (info: string) => set({ infoMsg: info }),
  setError: (err: string ) => set({ errorMsg: err })
}));
