import axios, { type InternalAxiosRequestConfig, type AxiosResponse } from 'axios'

const API_BASE_URL: string = import.meta.env.VITE_API_BASE_URL || '/api/v1'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to attach JWT access token and Admin Secret
api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  if (config.url && config.url.includes('/admin')) {
    config.headers['X-Admin-Secret-Key'] = 'teks-admin-secret-key-2026'
  }

  return config
})

// Response interceptor for automatic token refresh
api.interceptors.response.use(
  (response: AxiosResponse) => response,
  async (error) => {
    const originalRequest = error.config

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      const refreshToken = localStorage.getItem('refresh_token')

      if (refreshToken) {
        try {
          const res = await axios.post(`${API_BASE_URL}/auth/refresh`, {
            refresh_token: refreshToken,
          })

          if (res.data?.tokens?.access_token) {
            const newAccessToken: string = res.data.tokens.access_token
            const newRefreshToken: string = res.data.tokens.refresh_token

            localStorage.setItem('access_token', newAccessToken)
            if (newRefreshToken) {
              localStorage.setItem('refresh_token', newRefreshToken)
            }

            originalRequest.headers.Authorization = `Bearer ${newAccessToken}`
            return api(originalRequest)
          }
        } catch (refreshErr) {
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
          window.location.href = '/user/login?expired=true'
        }
      }
    }

    return Promise.reject(error)
  }
)

export default api
