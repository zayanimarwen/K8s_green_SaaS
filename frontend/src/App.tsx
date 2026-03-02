import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

import Login           from './pages/Login'
import Dashboard       from './pages/Dashboard'
import ClusterDetail   from './pages/ClusterDetail'
import WasteExplorer   from './pages/WasteExplorer'
import CarbonReport    from './pages/CarbonReport'
import Savings         from './pages/Savings'
import Recommendations from './pages/Recommendations'
import Simulator       from './pages/Simulator'
import History         from './pages/History'
import Diagnostics     from './pages/Diagnostics'
import AIModels        from './pages/AIModels'
import Settings        from './pages/Settings'
import Tenants         from './pages/admin/Tenants'
import Users           from './pages/admin/Users'

const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: 1, staleTime: 30_000 } }
})

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/login"           element={<Login />} />
          <Route path="/"                element={<Navigate to="/dashboard" replace />} />
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
