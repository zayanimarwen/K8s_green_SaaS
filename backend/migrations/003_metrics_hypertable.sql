-- TimescaleDB hypertable pour les métriques time-series
-- Partitionnement automatique par 1 jour

CREATE TABLE IF NOT EXISTS pod_metrics (
    time           TIMESTAMPTZ NOT NULL,
    cluster_id     UUID NOT NULL,
    pod_name       VARCHAR(255) NOT NULL,
    container_name VARCHAR(255) NOT NULL,
    namespace      VARCHAR(128) NOT NULL,
    node_name      VARCHAR(255),
    -- Ressources configurées
    cpu_request_m   DOUBLE PRECISION,
    cpu_limit_m     DOUBLE PRECISION,
    mem_request_mi  DOUBLE PRECISION,
    mem_limit_mi    DOUBLE PRECISION,
    -- Usage réel P95/24h
    cpu_usage_p95_m   DOUBLE PRECISION,
    mem_usage_p95_mi  DOUBLE PRECISION,
    -- Gaspillage calculé
    cpu_waste_m     DOUBLE PRECISION,
    mem_waste_mi    DOUBLE PRECISION,
    cost_waste_eur  DOUBLE PRECISION,
    has_limits      BOOLEAN,
    restart_count   INTEGER
);

SELECT create_hypertable('pod_metrics', 'time', chunk_time_interval => INTERVAL '1 day', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS idx_pod_metrics_cluster ON pod_metrics(cluster_id, time DESC);
CREATE INDEX IF NOT EXISTS idx_pod_metrics_ns ON pod_metrics(namespace, time DESC);

-- Compression automatique après 7 jours
ALTER TABLE pod_metrics SET (timescaledb.compress, timescaledb.compress_segmentby = 'cluster_id,namespace');
SELECT add_compression_policy('pod_metrics', INTERVAL '7 days');
-- Rétention 90 jours
SELECT add_retention_policy('pod_metrics', INTERVAL '90 days');
