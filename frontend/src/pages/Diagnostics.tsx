import { useState } from 'react'
import { useQuery, useMutation } from '@tanstack/react-query'
import { apiClient } from '../lib/api'
import { useClusterStore } from '../store/clusterStore'
import PageLayout from '../components/layout/PageLayout'
import MetricCard from '../components/ui/MetricCard'

interface Issue {
  id: string
  type: string
  severity: 'critical' | 'warning' | 'info'
  namespace: string
  resource_kind: string
  resource_name: string
  message: string
  ai_explanation?: string
  detected_at: string
}

interface DiagnosticsResult {
  cluster_id: string
  health_score: number
  issues: Issue[]
  summary: string
  llm_backend?: string
  analyzed_at: string
}

const severityColor = (s: string) => {
  switch (s) {
    case 'critical': return 'bg-red-50 border-red-200 text-red-800'
    case 'warning':  return 'bg-yellow-50 border-yellow-200 text-yellow-800'
    default:         return 'bg-blue-50 border-blue-200 text-blue-800'
  }
}

const severityBadge = (s: string) => {
  switch (s) {
    case 'critical': return 'bg-red-100 text-red-700 ring-red-600/20'
    case 'warning':  return 'bg-yellow-100 text-yellow-700 ring-yellow-600/20'
    default:         return 'bg-blue-100 text-blue-700 ring-blue-600/20'
  }
}

const healthColor = (score: number) => {
  if (score >= 80) return 'text-green-600'
  if (score >= 50) return 'text-yellow-600'
  return 'text-red-600'
}

export default function Diagnostics() {
  const { selectedCluster } = useClusterStore()
  const [expandedIssue, setExpandedIssue] = useState<string | null>(null)
  const [aiAnalysis, setAiAnalysis] = useState<Record<string, string>>({})

  const { data, isLoading, refetch } = useQuery<DiagnosticsResult>({
    queryKey: ['diagnostics', selectedCluster],
    queryFn: () => apiClient.get(`/clusters/${selectedCluster}/diagnostics`).then(r => r.data),
    enabled: !!selectedCluster,
    refetchInterval: 60_000,
  })

  const analyzeIssueMutation = useMutation({
    mutationFn: (issue: Issue) =>
      apiClient.post(`/clusters/${selectedCluster}/analyze-issue`, {
        type: issue.type,
        severity: issue.severity,
        namespace: issue.namespace,
        resource_kind: issue.resource_kind,
        resource_name: issue.resource_name,
        message: issue.message,
      }).then(r => r.data),
    onSuccess: (result, issue) => {
      setAiAnalysis(prev => ({ ...prev, [issue.id]: result.explanation }))
    }
  })

  const criticalCount = data?.issues?.filter(i => i.severity === 'critical').length ?? 0
  const warningCount  = data?.issues?.filter(i => i.severity === 'warning').length ?? 0

  return (
    <PageLayout title="AI Diagnostics" subtitle="Analyse intelligente de la sante du cluster">

      {/* Cartes de resume */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
        <MetricCard
          title="Score Sante"
          value={`${Math.round(data?.health_score ?? 0)}%`}
          trend="neutral"
          valueClass={healthColor(data?.health_score ?? 0)}
        />
        <MetricCard
          title="Problemes Critiques"
          value={String(criticalCount)}
          trend={criticalCount > 0 ? 'down' : 'neutral'}
          valueClass={criticalCount > 0 ? 'text-red-600' : 'text-green-600'}
        />
        <MetricCard
          title="Avertissements"
          value={String(warningCount)}
          trend={warningCount > 0 ? 'down' : 'neutral'}
          valueClass={warningCount > 0 ? 'text-yellow-600' : 'text-green-600'}
        />
      </div>

      {/* Resume LLM */}
      {data?.summary && (
        <div className="bg-gradient-to-r from-indigo-50 to-purple-50 border border-indigo-200 rounded-xl p-5 mb-6">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-2xl">AI</span>
            <h3 className="font-semibold text-indigo-900">Analyse IA</h3>
            {data.llm_backend && (
              <span className="ml-auto text-xs bg-indigo-100 text-indigo-700 px-2 py-1 rounded-full">
                {data.llm_backend}
              </span>
            )}
          </div>
          <p className="text-slate-700 text-sm whitespace-pre-wrap leading-relaxed">{data.summary}</p>
          <p className="text-xs text-slate-400 mt-3">
            Analyse du {new Date(data.analyzed_at).toLocaleString('fr-FR')}
          </p>
        </div>
      )}

      {/* Bouton refresh */}
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-lg font-semibold text-slate-800">
          {data?.issues?.length ?? 0} probleme(s) detecte(s)
        </h2>
        <button
          onClick={() => refetch()}
          disabled={isLoading}
          className="px-4 py-2 bg-indigo-600 text-white text-sm rounded-lg hover:bg-indigo-700 disabled:opacity-50 transition-colors"
        >
          {isLoading ? 'Analyse...' : 'Actualiser'}
        </button>
      </div>

      {/* Liste des issues */}
      {isLoading ? (
        <div className="text-center py-16 text-slate-400">Analyse en cours...</div>
      ) : !data?.issues?.length ? (
        <div className="text-center py-16">
          <div className="text-5xl mb-4">OK</div>
          <p className="text-green-600 font-semibold text-lg">Aucun probleme detecte</p>
          <p className="text-slate-400 text-sm mt-1">Cluster en bonne sante</p>
        </div>
      ) : (
        <div className="space-y-3">
          {data.issues.map(issue => (
            <div
              key={issue.id}
              className={`border rounded-xl p-4 cursor-pointer transition-all ${severityColor(issue.severity)}`}
              onClick={() => setExpandedIssue(expandedIssue === issue.id ? null : issue.id)}
            >
              {/* Header issue */}
              <div className="flex items-start justify-between gap-3">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 flex-wrap">
                    <span className={`text-xs font-medium px-2 py-0.5 rounded-full ring-1 ring-inset ${severityBadge(issue.severity)}`}>
                      {issue.severity.toUpperCase()}
                    </span>
                    <span className="text-xs bg-white/60 px-2 py-0.5 rounded-full font-mono">
                      {issue.type}
                    </span>
                    <span className="text-xs text-slate-500">
                      {issue.resource_kind} / {issue.namespace || 'cluster'} / {issue.resource_name}
                    </span>
                  </div>
                  <p className="mt-2 text-sm font-medium">{issue.message}</p>
                </div>
                <span className="text-slate-400 text-lg">{expandedIssue === issue.id ? 'v' : '>'}</span>
              </div>

              {/* Detail expandable */}
              {expandedIssue === issue.id && (
                <div className="mt-4 border-t border-current/20 pt-4 space-y-3">
                  {/* Analyse IA individuelle */}
                  {aiAnalysis[issue.id] ? (
                    <div className="bg-white/70 rounded-lg p-3">
                      <p className="text-xs font-semibold text-indigo-700 mb-2">Analyse IA detaillee :</p>
                      <p className="text-sm whitespace-pre-wrap text-slate-700">{aiAnalysis[issue.id]}</p>
                    </div>
                  ) : (
                    <button
                      onClick={e => { e.stopPropagation(); analyzeIssueMutation.mutate(issue) }}
                      disabled={analyzeIssueMutation.isPending}
                      className="w-full py-2 bg-white/80 hover:bg-white text-indigo-700 text-sm font-medium rounded-lg border border-indigo-300 transition-colors disabled:opacity-50"
                    >
                      {analyzeIssueMutation.isPending ? 'Analyse en cours...' : 'Analyser avec IA'}
                    </button>
                  )}

                  {/* Commandes kubectl suggeries */}
                  <div className="bg-slate-900 rounded-lg p-3">
                    <p className="text-xs text-slate-400 mb-2">Commandes rapides :</p>
                    <code className="text-xs text-green-400 block">
                      kubectl describe {issue.resource_kind.toLowerCase()} {issue.resource_name} -n {issue.namespace || 'default'}
                    </code>
                    <code className="text-xs text-green-400 block mt-1">
                      kubectl logs {issue.resource_name} -n {issue.namespace || 'default'} --previous
                    </code>
                  </div>

                  <p className="text-xs text-slate-500">
                    Detecte le {new Date(issue.detected_at).toLocaleString('fr-FR')}
                  </p>
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Config LLM manquante */}
      {data && !data.llm_backend && (
        <div className="mt-6 bg-amber-50 border border-amber-200 rounded-xl p-4">
          <p className="text-amber-800 text-sm font-semibold">LLM non configure</p>
          <p className="text-amber-700 text-sm mt-1">
            Pour activer l analyse IA, ajoutez dans vos secrets Kubernetes :
          </p>
          <pre className="bg-amber-100 rounded p-2 mt-2 text-xs text-amber-900 overflow-x-auto">{`AI_BACKEND=openai      # openai | claude | ollama
AI_API_KEY=sk-...      # Votre cle API
AI_MODEL=gpt-4o        # Modele a utiliser`}</pre>
        </div>
      )}
    </PageLayout>
  )
}
