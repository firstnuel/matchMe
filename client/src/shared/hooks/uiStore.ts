import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";

interface UIState {
  view: "chat" | "home" | "connections" | "profile" | "edit-profile";
  setView: (v: "chat" | "home" | "connections" | "profile" | "edit-profile") => void;
  isChatMessageViewActive: boolean;
  setIsChatMessageViewActive: (active: boolean) => void;
  errorMsg: string;
  infoMsg: string;
  clearMsgs: () => void;
  setInfo: (info: string) => void;
  setError: (err: string) => void;
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      view: "home",
      setView: (v) => set({ view: v }),
      isChatMessageViewActive: false,
      setIsChatMessageViewActive: (active) => set({ isChatMessageViewActive: active }),
      errorMsg: "",
      infoMsg: "",
      clearMsgs: () => set({ infoMsg: "", errorMsg: "" }),
      setInfo: (info) => set({ infoMsg: info }),
      setError: (err) => set({ errorMsg: err }),
    }),
    {
      name: "ui-storage", // key in sessionStorage
      storage: createJSONStorage(() => sessionStorage),
    }
  )
);
