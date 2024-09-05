import { create } from "zustand";

const defaultLanguage = "en-US";

const initialState = {
  role: "",
  skills: "",
  language: defaultLanguage,
  interviewId: "",
  interviewSecret: "",
  initialAudio: "",
  initialText: "",
  initialSSML: "",
  messages: [],
  isIntroDone: false,
  hasEnded: false,
};

export interface Message {
  text: string;
  isUser: boolean;
  isAnimated?: boolean;
  hasAnimated?: boolean;
}

interface InterviewState {
  role: string;
  skills: string;
  language: string;
  interviewId: string;
  interviewSecret: string;
  initialAudio: string;
  initialText: string;
  initialSSML: string;
  messages: Array<Message>;
  isIntroDone: boolean;
  hasEnded: boolean;

  setRole: (role: string) => void;
  setSkills: (skills: string) => void;
  setLanguage: (language: string) => void;
  setInterviewId: (id: string) => void;
  setInterviewSecret: (secret: string) => void;
  setInitialAudio: (audio: string) => void;
  setInitialText: (text: string) => void;
  setInitialSSML: (text: string) => void;
  setMessages: (messages: Array<Message>) => void;
  addMessage: (message: Message) => void;
  setIsIntroDone: (isIntroDone: boolean) => void;
  setHasEnded: (hasEnded: boolean) => void;

  resetStore: () => void;
}

export const useInterviewStore = create<InterviewState>((set) => ({
  ...initialState,

  setRole: (role) => set({ role }),
  setSkills: (skills) => set({ skills }),
  setLanguage: (language) => set({ language }),
  setInterviewId: (id) => set({ interviewId: id }),
  setInterviewSecret: (secret) => set({ interviewSecret: secret }),
  setInitialAudio: (audio) => set({ initialAudio: audio }),
  setInitialText: (text) => set({ initialText: text }),
  setInitialSSML: (text) => set({ initialSSML: text }),
  setMessages: (messages) => set({ messages }),
  addMessage: (message) =>
    set((state) => ({ messages: [...state.messages, message] })),
  setIsIntroDone: (isIntroDone) => set({ isIntroDone }),
  setHasEnded: (hasEnded) => set({ hasEnded }),

  resetStore: () => set({ ...initialState }),
}));
