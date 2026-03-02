import { useTenantStore } from '@/store/tenantStore'

export function useTenant() {
  const { tenants, currentTenantId, switchTenant } = useTenantStore()
  const current = tenants.find(t => t.id === currentTenantId) ?? null
  return { tenants, current, currentTenantId, switchTenant }
}
