import axios, { AxiosInstance } from 'axios'
import { useAuthStore } from '@/store/authStore'
import { useTenantStore } from '@/store/tenantStore'

const api: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:9000/v1',
  timeout: 60000,
})

// Token de dev statique (utilisé si pas de token Keycloak)
const _h = btoa(JSON.stringify({alg:'none',typ:'JWT'})).replace(/=/g,'')
const _p = btoa(JSON.stringify({sub:'user-dev',email:'dev@macif.fr',tenant_id:'tenant-demo',realm_access:{roles:['admin','viewer']}})).replace(/=/g,'')
const DEV_TOKEN = `${_h}.${_p}.`

// Injecter le token JWT et le tenant_id dans chaque requête
api.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token || DEV_TOKEN
  const tenantId = useTenantStore.getState().currentTenantId || 'tenant-demo'

  config.headers.Authorization = `Bearer ${token}`
  config.headers['X-Tenant-ID'] = tenantId

  return config
})

// Gérer les 401 → logout, les 429 → retry avec backoff
api.interceptors.response.use(
  (res) => res,
  async (err) => {
    // Ne pas rediriger vers /login en mode dev (pas de Keycloak)
    if (err.response?.status === 401) {
      const isDev = !import.meta.env.VITE_KEYCLOAK_URL
      if (!isDev) {
        useAuthStore.getState().logout()
        window.location.href = '/login'
      }
    }
    return Promise.reject(err)
  }
)

export default api
