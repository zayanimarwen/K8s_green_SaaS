import { NavLink } from 'react-router-dom'
import { useAuthStore } from '../../store/authStore'

const navItems = [
  { path: '/dashboard',       label: 'Dashboard',       icon: 'H'   },
  { path: '/waste',           label: 'Gaspillage',      icon: 'W'   },
  { path: '/carbon',          label: 'Carbone',         icon: 'C'   },
  { path: '/savings',         label: 'Economies',       icon: 'S'   },
  { path: '/recommendations', label: 'Recommandations', icon: 'R'   },
  { path: '/diagnostics',     label: 'AI Diagnostics',  icon: 'AI', badge: 'NEW' },
  { path: '/ai/models',       label: 'AI Models',       icon: 'M'   },
  { path: '/simulator',       label: 'Simulateur',      icon: 'Sim' },
  { path: '/history',         label: 'Historique',      icon: 'Hst' },
  { path: '/settings',        label: 'Parametres',      icon: 'Set' },
]

const adminItems = [
  { path: '/admin/tenants', label: 'Tenants', icon: 'T' },
  { path: '/admin/users',   label: 'Users',   icon: 'U' },
]

export function Sidebar() {
  const { user } = useAuthStore()
  const isSuperAdmin = user?.roles?.includes('superadmin')

  return (
    <aside className="w-64 bg-slate-900 text-white flex flex-col min-h-screen">
      <div className="px-6 py-5 border-b border-slate-700">
        <h1 className="text-xl font-bold text-green-400">K8s GreenOps</h1>
        <p className="text-slate-400 text-xs mt-1">Powered by Ollama (local)</p>
      </div>

      <nav className="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
        {navItems.map(item => (
          <NavLink
            key={item.path}
            to={item.path}
            className={({ isActive }) =>
              `flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors ${
                isActive
                  ? 'bg-green-600 text-white'
                  : 'text-slate-300 hover:bg-slate-800 hover:text-white'
              }`
            }
          >
            <span className="w-8 h-8 flex items-center justify-center bg-slate-700 rounded-md text-xs font-bold flex-shrink-0">
              {item.icon}
            </span>
            <span className="flex-1">{item.label}</span>
            {item.badge && (
              <span className="text-xs bg-indigo-500 text-white px-1.5 py-0.5 rounded-full">
                {item.badge}
              </span>
            )}
          </NavLink>
        ))}

        {isSuperAdmin && (
          <>
            <div className="pt-4 pb-2 px-3">
              <p className="text-xs text-slate-500 uppercase tracking-wider font-semibold">Admin</p>
            </div>
            {adminItems.map(item => (
              <NavLink
                key={item.path}
                to={item.path}
                className={({ isActive }) =>
                  `flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors ${
                    isActive
                      ? 'bg-green-600 text-white'
                      : 'text-slate-300 hover:bg-slate-800 hover:text-white'
                  }`
                }
              >
                <span className="w-8 h-8 flex items-center justify-center bg-slate-700 rounded-md text-xs">
                  {item.icon}
                </span>
                {item.label}
              </NavLink>
            ))}
          </>
        )}
      </nav>

      <div className="px-6 py-4 border-t border-slate-700">
        <p className="text-xs text-slate-500">GreenOps v1.0 + AI local</p>
      </div>
    </aside>
  )
}
