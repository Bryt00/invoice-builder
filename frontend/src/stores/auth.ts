import { defineStore } from 'pinia'
import api from '@/utils/api'
import type { User, BusinessProfile, AuthTokenPair, LoginResponse } from '../types/user'

const currencyMap: Record<string, string> = {
  USD: '$', EUR: '€', GBP: '£', GHS: 'GH₵',
  NGN: '₦', CAD: 'CA$', AUD: 'A$', JPY: '¥',
  INR: '₹', ZAR: 'R', BRL: 'R$', AED: 'AED'
}

export interface AuthState {
  user: User | null
  profile: BusinessProfile | null
  credits: number
  accessToken: string | null
  refreshToken: string | null
  loading: boolean
  error: string | null
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    profile: null,
    credits: 0,
    accessToken: localStorage.getItem('access_token') || null,
    refreshToken: localStorage.getItem('refresh_token') || null,
    loading: false,
    error: null,
  }),

  getters: {
    isAuthenticated: (state): boolean => !!state.accessToken && !!state.user,
    isProfileComplete: (state): boolean => state.user?.is_profile_complete || false,
    currencySymbol: (state): string => {
      const code = state.profile?.default_currency || 'USD'
      return currencyMap[code] || code
    }
  },

  actions: {
    setTokens(accessToken: string, refreshToken?: string) {
      this.accessToken = accessToken
      if (refreshToken) {
        this.refreshToken = refreshToken
      }
      localStorage.setItem('access_token', accessToken)
      if (refreshToken) {
        localStorage.setItem('refresh_token', refreshToken)
      }
    },

    clearAuth() {
      this.user = null
      this.profile = null
      this.credits = 0
      this.accessToken = null
      this.refreshToken = null
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
    },

    async fetchCurrentUser() {
      if (!this.accessToken) return
      this.loading = true
      this.error = null

      try {
        const res = await api.get<{ status: string; user: User; profile: BusinessProfile; credits: number }>('/auth/me')
        if (res.data?.user) {
          this.user = res.data.user
          this.profile = res.data.profile || null
          this.credits = res.data.credits || 0
        }
      } catch (err: any) {
        this.error = err.response?.data?.error || 'Failed to fetch current user'
        this.clearAuth()
      } finally {
        this.loading = false
      }
    },

    async login(email: string, password: string): Promise<LoginResponse> {
      this.loading = true
      this.error = null

      try {
        const res = await api.post<LoginResponse>('/auth/login', { email, password })
        const data = res.data

        if (data?.tokens?.access_token) {
          this.setTokens(data.tokens.access_token, data.tokens.refresh_token)
          this.user = data.user
          this.profile = data.profile || null
          this.credits = data.credits || 0
        }
        return data
      } catch (err: any) {
        this.error = err.response?.data?.error || 'Invalid credentials'
        throw err
      } finally {
        this.loading = false
      }
    },

    async register(name: string, email: string, password: string): Promise<any> {
      this.loading = true
      this.error = null

      try {
        const res = await api.post('/auth/register', { name, email, password })
        return res.data
      } catch (err: any) {
        this.error = err.response?.data?.error || 'Registration failed'
        throw err
      } finally {
        this.loading = false
      }
    },

    async logout() {
      try {
        await api.post('/auth/logout')
      } catch (err: any) {
        // Ignore network errors on logout
      } finally {
        this.clearAuth()
      }
    }
  }
})
