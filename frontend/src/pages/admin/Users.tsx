import { PageLayout } from '@/components/layout/PageLayout'
import { useQuery } from '@tanstack/react-query'
import api from '@/lib/api'
import type { User } from '@/types/api'

export function Users() {
  const { data, isLoading } = useQuery<{ users: User[] }>({
    queryKey: ['admin', 'users'],
    queryFn: async () => { const { data } = await api.get('/v1/admin/users'); return data },
  })
  return (
    <PageLayout title="Gestion des utilisateurs">
      <div className="max-w-3xl">
        {isLoading ? <div className="text-slate-400">Chargement...</div> : (
          <div className="bg-white rounded-xl border border-slate-100 overflow-hidden">
            <table className="w-full text-sm">
              <thead className="bg-slate-50 border-b border-slate-100">
                <tr>{['Email', 'Nom', 'Rôle', 'Tenant'].map(h => (
                  <th key={h} className="px-4 py-3 text-left text-xs font-semibold text-slate-500">{h}</th>
                ))}</tr>
              </thead>
              <tbody>
                {data?.users.map(u => (
                  <tr key={u.id} className="border-b border-slate-50 hover:bg-slate-50">
                    <td className="px-4 py-3">{u.email}</td>
                    <td className="px-4 py-3">{u.display_name ?? '—'}</td>
                    <td className="px-4 py-3"><span className="bg-navy-50 text-navy-700 text-xs px-2 py-1 rounded font-medium">{u.role}</span></td>
                    <td className="px-4 py-3 text-slate-500 text-xs font-mono">{u.tenant_id}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </PageLayout>
  )
}
