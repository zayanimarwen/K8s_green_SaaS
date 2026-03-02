import { clsx, type ClassValue } from 'clsx'

export function cn(...inputs: ClassValue[]) {
  return clsx(inputs)
}

/** Formate un score en couleur Tailwind */
export function scoreColor(score: number): string {
  if (score >= 80) return 'text-green-500'
  if (score >= 60) return 'text-yellow-500'
  if (score >= 40) return 'text-orange-500'
  return 'text-red-500'
}

/** Formate un montant en euros */
export function formatEur(amount: number): string {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR', maximumFractionDigits: 0 }).format(amount)
}

/** Formate du CO₂ */
export function formatCO2(kg: number): string {
  if (kg >= 1000) return `${(kg / 1000).toFixed(1)} tCO₂`
  return `${kg.toFixed(0)} kgCO₂`
}

/** Couleur de priorité */
export function priorityColor(p: 'HIGH' | 'MEDIUM' | 'LOW'): string {
  return { HIGH: 'text-red-500', MEDIUM: 'text-yellow-500', LOW: 'text-green-500' }[p]
}

/** Couleur de grade */
export function gradeColor(grade: string): string {
  const map: Record<string, string> = {
    'A+': '#22c55e', 'A': '#4ade80', 'B+': '#84cc16',
    'B': '#eab308', 'C': '#f97316', 'D': '#ef4444', 'F': '#dc2626'
  }
  return map[grade] ?? '#94a3b8'
}

/** Date relative */
export function relativeTime(isoDate: string): string {
  const diff = Date.now() - new Date(isoDate).getTime()
  const minutes = Math.floor(diff / 60000)
  if (minutes < 1) return 'à l'instant'
  if (minutes < 60) return `il y a ${minutes}min`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `il y a ${hours}h`
  return `il y a ${Math.floor(hours / 24)}j`
}
