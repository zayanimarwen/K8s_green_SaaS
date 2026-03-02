import axios, { AxiosInstance } from 'axios'
import { useAuthStore } from '@/store/authStore'
import { useTenantStore } from '@/store/tenantStore'

const api: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/v1',
  timeout: 30000,
})

// Injecter le token JWT et le tenant_id dans chaque requête
api.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token
  const tenantId = useTenantStore.getState().currentTenantId

  if (token) config.headers.Authorization = `Bearer ${token}`
  if (tenantId) config.headers['X-Tenant-ID'] = tenantId

  return config
})

// Gérer les 401 → logout, les 429 → retry avec backoff
api.interceptors.response.use(
  (res) => res,
  async (err) => {
    if (err.response?.status === 401) {
      useAuthStore.getState().logout()
      window.location.href = '/login'
    }
    return Promise.reject(err)
  }
)

export default api
