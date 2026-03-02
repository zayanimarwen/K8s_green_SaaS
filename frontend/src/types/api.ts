// Types API partagés — alignés sur les réponses du backend Go

export interface Cluster {
  id: string
  name: string
  provider: 'aws' | 'gcp' | 'azure' | 'on-prem'
  region: string
  environment: 'production' | 'staging' | 'dev'
  k8s_version?: string
  last_seen_at?: string
  active: boolean
}

export interface ScoreBreakdown {
  cpu_efficiency: number
  mem_efficiency: number
  node_packing: number
  hpa_coverage: number
  limit_compliance: number
}

export interface GreenScore {
  cluster_id: string
  score: number
  grade: 'A+' | 'A' | 'B+' | 'B' | 'C' | 'D' | 'F' | 'N/A'
  label?: string
  breakdown: ScoreBreakdown
  timestamp?: string
}

export interface WasteReport {
  cluster_id: string
  pods_overprovisioned: number
  total_cpu_waste_cores: number
  total_mem_waste_gb: number
  annual_waste_eur: number
  waste_reports: PodWaste[]
}

export interface PodWaste {
  pod_name: string
  container_name: string
  namespace: string
  cpu_request_m: number
  cpu_usage_p95_m: number
  cpu_optimal_m: number
  cpu_waste_pct: number
  mem_request_mi: number
  mem_usage_p95_mi: number
  mem_optimal_mi: number
  mem_waste_pct: number
  annual_cost_waste_eur: number
  priority: 'HIGH' | 'MEDIUM' | 'LOW'
  confidence: number
  has_limits: boolean
  is_hpa_managed: boolean
}

export interface CarbonReport {
  cluster_id: string
  co2_kg_annual: number
  kwh_annual: number
  equivalent_km_car: number
  equivalent_trees: number
  provider: string
  carbon_intensity: number
}

export interface SavingsReport {
  cluster_id: string
  annual_savings_eur: number
  monthly_savings_eur: number
  breakdown: {
    rightsizing_eur: number
    node_consolidation_eur: number
    hpa_automation_eur: number
  }
}

export interface Recommendation {
  id?: string
  priority: 'HIGH' | 'MEDIUM' | 'LOW'
  type: 'rightsizing' | 'missing_limits' | 'add_hpa' | 'node_consolidation'
  title: string
  description: string
  savings_eur_annual: number
  target: string
  confidence: number
  patch_yaml?: string
  status?: 'open' | 'applied' | 'dismissed'
}

export interface SimulateRequest {
  deployment: string
  namespace: string
  new_cpu_req_m?: number
  new_mem_req_mi?: number
}

export interface SimulateResult {
  cluster_id: string
  deployment: string
  projected_savings_eur_annual: number
  projected_score_delta: number
  projected_co2_reduction_kg: number
  safe_to_apply: boolean
  confidence: number
}

export interface Tenant {
  id: string
  name: string
  plan: 'starter' | 'pro' | 'enterprise'
  clusters: number
  active: boolean
}

export interface User {
  id: string
  email: string
  display_name?: string
  role: 'superadmin' | 'admin' | 'viewer'
  tenant_id: string
}

export interface ApiError {
  error: string
  code?: string
}
