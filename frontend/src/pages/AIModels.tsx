import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import apiClient from '../lib/api'
import { PageLayout } from '../components/layout/PageLayout'

interface ModelInfo {
  name: string
  size: number
  description: string
  pulled: boolean
}

interface ModelsResponse {
  current_model: string
  ollama_url: string
  models: ModelInfo[]
}

function formatSize(bytes: number): string {
  if (!bytes) return 'inconnu'
  const gb = bytes / 1e9
  return gb >= 1 ? `${gb.toFixed(1)} Go` : `${(bytes / 1e6).toFixed(0)} Mo`
}

const MODEL_COLORS: Record<string, string> = {
  'llama3.2':    'bg-blue-50 border-blue-200',
  'llama3.1:8b': 'bg-indigo-50 border-indigo-200',
  'mistral':     'bg-purple-50 border-purple-200',
  'codellama':   'bg-green-50 border-green-200',
  'phi3:mini':   'bg-orange-50 border-orange-200',
  'gemma2:2b':   'bg-pink-50 border-pink-200',
}

export function AIModels() {
  const qc = useQueryClient()
  const [pulling, setPulling] = useState<string | null>(null)

  const { data, isLoading } = useQuery<ModelsResponse>({
    queryKey: ['ai-models'],
    queryFn: () => apiClient.get('/ai/models').then(r => r.data),
    refetchInterval: 5000,
  })

  const switchMutation = useMutation({
    mutationFn: (model: string) =>
      apiClient.post('/ai/models/switch', { model }).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['ai-models'] }),
  })

  const pullMutation = useMutation({
    mutationFn: (model: string) => {
      setPulling(model)
      return apiClient.post('/ai/models/pull', { model }).then(r => r.data)
    },
    onSuccess: () => {
      setTimeout(() => {
        setPulling(null)
        qc.invalidateQueries({ queryKey: ['ai-models'] })
      }, 3000)
    },
    onError: () => setPulling(null),
  })

  const deleteMutation = useMutation({
    mutationFn: (model: string) =>
      apiClient.delete(`/ai/models/${encodeURIComponent(model)}`).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['ai-models'] }),
  })

  return (
    <PageLayout title="AI Models" subtitle="Gestion des modeles Ollama locaux — 100% prive, aucune donnee externe">

      {/* Info Ollama */}
      <div className="bg-gradient-to-r from-green-50 to-emerald-50 border border-green-200 rounded-xl p-4 mb-6">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center text-green-700 font-bold">
            OK
          </div>
          <div>
            <p className="font-semibold text-green-900">Ollama local actif</p>
            <p className="text-sm text-green-700">{data?.ollama_url ?? 'http://localhost:11434'}</p>
          </div>
          <div className="ml-auto text-right">
            <p className="text-xs text-green-600">Modele actif</p>
            <p className="font-mono font-bold text-green-800">{data?.current_model ?? '...'}</p>
          </div>
        </div>
      </div>

      {/* Confidentialite */}
      <div className="bg-slate-50 border border-slate-200 rounded-xl p-4 mb-6">
        <p className="text-sm font-semibold text-slate-700 mb-1">Confidentialite garantie</p>
        <p className="text-xs text-slate-500">
          Tous les modeles tournent localement sur votre machine. Vos donnees Kubernetes,
          logs et evenements ne quittent jamais votre infrastructure.
        </p>
      </div>

      {/* Liste des modeles */}
      <h2 className="text-lg font-semibold text-slate-800 mb-4">
        Modeles disponibles
        <span className="ml-2 text-sm font-normal text-slate-400">
          ({data?.models?.filter(m => m.pulled).length ?? 0} installes)
        </span>
      </h2>

      {isLoading ? (
        <div className="text-center py-12 text-slate-400">Chargement...</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {data?.models?.map(model => {
            const isActive  = model.name === data.current_model
            const isPulling = pulling === model.name
            const cardColor = MODEL_COLORS[model.name] ?? 'bg-slate-50 border-slate-200'

            return (
              <div
                key={model.name}
                className={`border rounded-xl p-4 transition-all ${
                  isActive ? 'ring-2 ring-indigo-500 ' + cardColor : cardColor
                }`}
              >
                {/* Header */}
                <div className="flex items-start justify-between mb-2">
                  <div>
                    <div className="flex items-center gap-2">
                      <h3 className="font-mono font-bold text-slate-800">{model.name}</h3>
                      {isActive && (
                        <span className="text-xs bg-indigo-600 text-white px-2 py-0.5 rounded-full">
                          ACTIF
                        </span>
                      )}
                      {model.pulled && !isActive && (
                        <span className="text-xs bg-green-100 text-green-700 px-2 py-0.5 rounded-full">
                          installe
                        </span>
                      )}
                    </div>
                    <p className="text-xs text-slate-500 mt-1">{model.description}</p>
                  </div>
                  {model.pulled && (
                    <span className="text-xs text-slate-400 font-mono whitespace-nowrap ml-2">
                      {formatSize(model.size)}
                    </span>
                  )}
                </div>

                {/* Actions */}
                <div className="flex gap-2 mt-3">
                  {model.pulled ? (
                    <>
                      {!isActive && (
                        <button
                          onClick={() => switchMutation.mutate(model.name)}
                          disabled={switchMutation.isPending}
                          className="flex-1 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-medium rounded-lg transition-colors disabled:opacity-50"
                        >
                          Activer
                        </button>
                      )}
                      {!isActive && (
                        <button
                          onClick={() => {
                            if (confirm('Supprimer ' + model.name + ' ?')) {
                              deleteMutation.mutate(model.name)
                            }
                          }}
                          className="px-3 py-1.5 bg-red-50 hover:bg-red-100 text-red-600 text-xs font-medium rounded-lg transition-colors border border-red-200"
                        >
                          Supprimer
                        </button>
                      )}
                      {isActive && (
                        <p className="text-xs text-indigo-600 font-medium py-1.5">
                          Ce modele est utilise pour les diagnostics
                        </p>
                      )}
                    </>
                  ) : (
                    <button
                      onClick={() => pullMutation.mutate(model.name)}
                      disabled={isPulling || !!pulling}
                      className="flex-1 py-1.5 bg-slate-700 hover:bg-slate-800 text-white text-xs font-medium rounded-lg transition-colors disabled:opacity-50"
                    >
                      {isPulling ? 'Telechargement...' : 'Telecharger'}
                    </button>
                  )}
                </div>

                {/* Progress pull */}
                {isPulling && (
                  <div className="mt-2">
                    <div className="h-1.5 bg-slate-200 rounded-full overflow-hidden">
                      <div className="h-full bg-indigo-500 rounded-full animate-pulse w-3/4" />
                    </div>
                    <p className="text-xs text-slate-400 mt-1">
                      Telechargement en cours — cela peut prendre plusieurs minutes...
                    </p>
                  </div>
                )}
              </div>
            )
          })}
        </div>
      )}

      {/* Guide choix modele */}
      <div className="mt-8 bg-slate-50 border border-slate-200 rounded-xl p-5">
        <h3 className="font-semibold text-slate-700 mb-3">Guide de choix</h3>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 text-sm">
          <div className="bg-white rounded-lg p-3 border border-slate-100">
            <p className="font-medium text-slate-700">Machine legere</p>
            <p className="text-slate-500 text-xs mt-1">RAM 8 Go</p>
            <p className="text-indigo-600 font-mono text-xs mt-2">phi3:mini ou llama3.2</p>
          </div>
          <div className="bg-white rounded-lg p-3 border border-slate-100">
            <p className="font-medium text-slate-700">Usage general</p>
            <p className="text-slate-500 text-xs mt-1">RAM 16 Go</p>
            <p className="text-indigo-600 font-mono text-xs mt-2">llama3.2 ou mistral</p>
          </div>
          <div className="bg-white rounded-lg p-3 border border-slate-100">
            <p className="font-medium text-slate-700">Analyse YAML/code</p>
            <p className="text-slate-500 text-xs mt-1">RAM 16 Go+</p>
            <p className="text-indigo-600 font-mono text-xs mt-2">codellama</p>
          </div>
        </div>
      </div>
    </PageLayout>
  )
}
