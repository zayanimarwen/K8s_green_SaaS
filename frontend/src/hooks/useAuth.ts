import { useEffect } from 'react'
import { useAuthStore } from '@/store/authStore'
import { getUser, handleCallback } from '@/lib/auth'

export function useAuth() {
  const { user, isAuthenticated, isLoading, setAuth, logout, setLoading } = useAuthStore()

  useEffect(() => {
    async function init() {
      try {
        // Gérer le retour OIDC callback
        if (window.location.search.includes('code=')) {
          const oidcUser = await handleCallback()
          setAuth({
            id:       oidcUser.profile.sub,
            email:    oidcUser.profile.email ?? '',
            name:     oidcUser.profile.name ?? '',
            roles:    (oidcUser.profile['realm_access'] as any)?.roles ?? [],
            tenantId: (oidcUser.profile['tenant_id'] as string) ?? '',
          }, oidcUser.access_token)
          window.history.replaceState({}, '', '/')
          return
        }
        // Session existante
        const existing = await getUser()
        if (existing && !existing.expired) {
          setAuth({
            id:       existing.profile.sub,
            email:    existing.profile.email ?? '',
            name:     existing.profile.name ?? '',
            roles:    (existing.profile['realm_access'] as any)?.roles ?? [],
            tenantId: (existing.profile['tenant_id'] as string) ?? '',
          }, existing.access_token)
        } else {
          setLoading(false)
        }
      } catch (_) {
        setLoading(false)
      }
    }
    init()
  }, [])

  return { user, isAuthenticated, isLoading, logout }
}
