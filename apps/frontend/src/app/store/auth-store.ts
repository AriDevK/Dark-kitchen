import { create } from "zustand";

type User = {
  id: number;
  email: string;
  role: string;
};

type AuthStore = {
  accessToken: string | null;
  refreshToken: string | null;
  user: User | null;

  setAuth: (
    accessToken: string,
    refreshToken: string,
    user: User
  ) => void;

  logout: () => void;
};

export const useAuthStore = create<AuthStore>((set) => ({
  accessToken: null,
  refreshToken: null,
  user: null,

  setAuth: (accessToken, refreshToken, user) =>
    set({
      accessToken,
      refreshToken,
      user,
    }),

  logout: () =>
    set({
      accessToken: null,
      refreshToken: null,
      user: null,
    }),
}));