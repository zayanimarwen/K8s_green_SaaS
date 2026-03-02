import { PageLayout } from '@/components/layout/PageLayout'
import { useQuery } from '@tanstack/react-query'
import api from '@/lib/api'
import type { Tenant } from '@/types/api'

export function Tenants() {
  const { data, isLoading } = useQuery<{ tenants: Tenant[] }>({
    queryKey: ['admin', 'tenants'],
    queryFn: async () => { const { data } = await api.get('/v1/admin/tenants'); return data },
  })
  return (
    <PageLayout title="Gestion des tenants">
      <div className="max-w-4xl">
        {isLoading ? <div className="text-slate-400">Chargement...</div> : (
          <div className="bg-white rounded-xl border border-slate-100 overflow-hidden">
            <table className="w-full text-sm">
              <thead className="bg-slate-50 border-b border-slate-100">
                <tr>{['ID', 'Nom', 'Plan', 'Clusters', 'Statut'].map(h => (
                  <th key={h} className="px-4 py-3 text-left text-xs font-semibold text-slate-500">{h}</th>
                ))}</tr>
              </thead>
              <tbody>
                {data?.tenants.map(t => (
                  <tr key={t.id} className="border-b border-slate-50 hover:bg-slate-50">
                    <td className="px-4 py-3 font-mono text-xs">{t.id}</td>
                    <td className="px-4 py-3 font-medium">{t.name}</td>
                    <td className="px-4 py-3"><span className="bg-brand-50 text-brand-700 text-xs px-2 py-1 rounded">{t.plan}</span></td>
                    <td className="px-4 py-3">{t.clusters}</td>
                    <td className="px-4 py-3"><span className={`text-xs px-2 py-1 rounded ${t.active ? 'bg-green-50 text-green-700' : 'bg-slate-50 text-slate-400'}`}>{t.active ? 'Actif' : 'Inactif'}</span></td>
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
