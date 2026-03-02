import { create } from 'zustand'

interface User {
  id: string
  email: string
  name: string
  roles: string[]
  tenantId: string
}

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  setAuth: (user: User, token: string) => void
  logout: () => void
  setLoading: (loading: boolean) => void
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: null,
  isAuthenticated: false,
  isLoading: true,

  setAuth: (user, token) => set({ user, token, isAuthenticated: true, isLoading: false }),

  logout: () => set({ user: null, token: null, isAuthenticated: false }),

  setLoading: (isLoading) => set({ isLoading }),
}))
