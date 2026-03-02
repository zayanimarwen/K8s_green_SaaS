import { login } from '@/lib/auth'

export function Login() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-navy-900 to-navy-800 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-2xl p-10 w-full max-w-md text-center">
        <div className="text-5xl mb-4">🌿</div>
        <h1 className="text-2xl font-bold text-slate-900 mb-2">K8s Green Optimizer</h1>
        <p className="text-slate-500 mb-8">Optimisez vos clusters Kubernetes.<br/>Réduisez votre empreinte carbone.</p>
        <button
          onClick={() => login()}
          className="w-full bg-brand-500 hover:bg-brand-600 text-white font-semibold py-3 px-6 rounded-xl transition-colors"
        >
          Se connecter avec SSO
        </button>
        <p className="text-xs text-slate-400 mt-4">Authentification via Keycloak · Compatible PingFederate</p>
      </div>
    </div>
  )
}
