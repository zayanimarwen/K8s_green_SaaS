import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Login }           from './pages/Login'
import { Dashboard }       from './pages/Dashboard'
import { ClusterDetail }   from './pages/ClusterDetail'
import { WasteExplorer }   from './pages/WasteExplorer'
import { CarbonReport }    from './pages/CarbonReport'
import { Savings }         from './pages/Savings'
import { Recommendations } from './pages/Recommendations'
import { Simulator }       from './pages/Simulator'
import { History }         from './pages/History'
import { Diagnostics }     from './pages/Diagnostics'
import { AIModels }        from './pages/AIModels'
import { Settings }        from './pages/Settings'
import { Tenants }         from './pages/admin/Tenants'
import { Users }           from './pages/admin/Users'

const qc = new QueryClient()

// Dev mode: initialiser les stores avec un token et tenant de dev
// En production, ces valeurs viennent de Keycloak
const DEV_TOKEN = (() => {
  const h = btoa(JSON.stringify({alg:'none',typ:'JWT'})).replace(/=/g,'')
  const p = btoa(JSON.stringify({sub:'user-dev',email:'dev@macif.fr',tenant_id:'tenant-demo',realm_access:{roles:['admin','viewer']}})).replace(/=/g,'')
  return `${h}.${p}.`
})()

import { useAuthStore } from './store/authStore'
import { useTenantStore } from './store/tenantStore'
import { useClusterStore } from './store/clusterStore'
import { useEffect } from 'react'

function DevInit() {
  useEffect(() => {
    const { token, setAuth } = useAuthStore.getState()
    if (!token) {
      setAuth({ id: 'user-dev', email: 'dev@macif.fr', name: 'Dev User', roles: ['admin'], tenantId: 'tenant-demo' }, DEV_TOKEN)
      useTenantStore.getState().setTenants([{ id: 'tenant-demo', name: 'MACIF Demo', plan: 'enterprise' }])
      useClusterStore.getState().setClusters([{ id: 'eb3f8f3b-202f-4ac9-a80a-530a925338cc', name: 'docker-desktop', provider: 'on-prem', region: 'local', environment: 'development', active: true }])
    }
  }, [])
  return null
}

export default function App() {
  return (
    <QueryClientProvider client={qc}>
      <BrowserRouter>
        <DevInit />
        <Routes>
          <Route path="/login"           element={<Login />} />
          <Route path="/dashboard"       element={<Dashboard />} />
          <Route path="/clusters/:id"    element={<ClusterDetail />} />
          <Route path="/waste"           element={<WasteExplorer />} />
          <Route path="/carbon"          element={<CarbonReport />} />
          <Route path="/savings"         element={<Savings />} />
          <Route path="/recommendations" element={<Recommendations />} />
          <Route path="/simulator"       element={<Simulator />} />
          <Route path="/history"         element={<History />} />
          <Route path="/diagnostics"     element={<Diagnostics />} />
          <Route path="/ai/models"       element={<AIModels />} />
          <Route path="/settings"        element={<Settings />} />
          <Route path="/admin/tenants"   element={<Tenants />} />
          <Route path="/admin/users"     element={<Users />} />
          <Route path="*"               element={<Navigate to="/dashboard" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  )
}
