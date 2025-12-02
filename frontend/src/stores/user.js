import { defineStore } from 'pinia'
import api, { setAuthToken } from '../api/http'

const persistedToken = typeof window !== 'undefined' ? localStorage.getItem('kpi_token') : ''
if (persistedToken) {
  setAuthToken(persistedToken)
}

export const useUserStore = defineStore('user', {
  state: () => ({
    token: persistedToken || '',
    profile: null,
    loading: false
  }),
  getters: {
    isAuthenticated: (state) => Boolean(state.token),
    fullName: (state) => state.profile?.fullName || '未登录'
  },
  actions: {
    async login(credentials) {
      this.loading = true
      try {
        const { data } = await api.post('/auth/login', credentials)
        this.token = data.data.token
        localStorage.setItem('kpi_token', this.token)
        setAuthToken(this.token)
        await this.fetchProfile()
      } finally {
        this.loading = false
      }
    },
    async fetchProfile() {
      try {
        const { data } = await api.get('/auth/profile')
        this.profile = data.data
        return this.profile
      } catch (error) {
        this.logout()
        throw error
      }
    },
    logout() {
      this.token = ''
      this.profile = null
      localStorage.removeItem('kpi_token')
      setAuthToken(null)
    }
  }
})
