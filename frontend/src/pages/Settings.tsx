import { PageLayout } from '@/components/layout/PageLayout'
import { useAuthStore } from '@/store/authStore'
import { logout } from '@/lib/auth'

export function Settings() {
  const { user } = useAuthStore()
  return (
    <PageLayout title="Paramètres">
      <div className="max-w-xl space-y-6">
        <div className="bg-white rounded-xl border border-slate-100 p-6">
          <h2 className="font-semibold mb-4">Profil</h2>
          <div className="space-y-2 text-sm">
            <div><span className="text-slate-500">Email: </span><strong>{user?.email}</strong></div>
            <div><span className="text-slate-500">Rôle: </span><strong>{user?.roles[0]}</strong></div>
            <div><span className="text-slate-500">Tenant: </span><strong>{user?.tenantId}</strong></div>
          </div>
        </div>
        <div className="bg-white rounded-xl border border-slate-100 p-6">
          <h2 className="font-semibold mb-4">Notifications</h2>
          <p className="text-sm text-slate-500">Configuration via la page admin ou API</p>
        </div>
        <button onClick={() => logout()}
          className="bg-red-50 text-red-600 hover:bg-red-100 px-4 py-2 rounded-lg text-sm font-medium transition-colors">
          Se déconnecter
        </button>
      </div>
    </PageLayout>
  )
}
