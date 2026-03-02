import { create } from 'zustand'
import type { Tenant } from '@/types/api'

interface TenantState {
  tenants: Tenant[]
  currentTenantId: string | null
  setTenants: (tenants: Tenant[]) => void
  switchTenant: (id: string) => void
}

export const useTenantStore = create<TenantState>((set) => ({
  tenants: [],
  currentTenantId: null,
  setTenants: (tenants) => set({ tenants, currentTenantId: tenants[0]?.id ?? null }),
  switchTenant: (id) => set({ currentTenantId: id }),
}))
