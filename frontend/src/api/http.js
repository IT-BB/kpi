import axios from 'axios'

const base = import.meta.env.VITE_API_BASE_URL?.replace(/\/$/, '') || ''
const api = axios.create({
  baseURL: `${base}/api/v1`,
  timeout: 10000
})

api.interceptors.request.use((config) => {
  if (config.url && config.url.startsWith('/') && !config.url.startsWith('//')) {
    config.url = config.url.slice(1)
  }
  return config
})

export const setAuthToken = (token) => {
  if (token) {
    api.defaults.headers.common.Authorization = `Bearer ${token}`
  } else {
    delete api.defaults.headers.common.Authorization
  }
}

export default api
