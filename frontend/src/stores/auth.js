import { defineStore } from 'pinia'
import api from '../utils/api'

const currencyMap = {
  USD: '$', EUR: '€', GBP: '£', GHS: 'GH₵',
  NGN: '₦', CAD: 'CA$', AUD: 'A$', JPY: '¥',
  INR: '₹', ZAR: 'R', BRL: 'R$', AED: 'AED'
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    profile: null,
    credits: 0,
    accessToken: localStorage.getItem('access_token') || null,
    refreshToken: localStorage.getItem('refresh_token') || null,
    loading: false,
    error: null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.accessToken && !!state.user,
    isProfileComplete: (state) => state.user?.is_profile_complete || false,
    currencySymbol: (state) => {
      const code = state.profile?.default_currency || 'USD'
      return currencyMap[code] || code
    }
  },

  actions: {
    setTokens(accessToken, refreshToken) {
      this.accessToken = accessToken
      this.refreshToken = refreshToken
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
      if (!this.accessToken) return null
      this.loading = true
      try {
        const res = await api.get('/auth/me')
        if (res.data?.user) {
          this.user = res.data.user
          this.profile = res.data.user.profile || null
          this.credits = res.data.user.credits || 0
        }
        return this.user
      } catch (err) {
        this.clearAuth()
        return null
      } finally {
        this.loading = false
      }
    },

    async login(email, password) {
      this.loading = true
      this.error = null
      try {
        const res = await api.post('/auth/login', { email, password })
        if (res.data?.tokens) {
          this.setTokens(res.data.tokens.access_token, res.data.tokens.refresh_token)
          this.user = res.data.user
          this.profile = res.data.user.profile || null
          this.credits = res.data.user.credits || 0
        }
        return res.data
      } catch (err) {
        if (!err.response) {
          this.error = 'Unable to connect to the server. Please try again later.'
        } else {
          this.error = err.response?.data?.error || 'Invalid credentials'
        }
        throw err
      } finally {
        this.loading = false
      }
    },

    async register(name, email, password, confirmPassword) {
      this.loading = true
      this.error = null
      try {
        const res = await api.post('/auth/register', {
          name,
          email,
          password,
          confirm_password: confirmPassword,
        })
        return res.data
      } catch (err) {
        if (!err.response) {
          this.error = 'Unable to connect to the server. Please try again later.'
        } else {
          this.error = err.response?.data?.error || 'Registration failed'
        }
        throw err
      } finally {
        this.loading = false
      }
    },

    async logout() {
      try {
        await api.post('/auth/logout')
      } catch (e) {
        // Ignore logout network errors
      } finally {
        this.clearAuth()
      }
    },
  },
})
