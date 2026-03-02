-- Historique des Green Scores par cluster
CREATE TABLE IF NOT EXISTS green_scores (
    time         TIMESTAMPTZ NOT NULL,
    cluster_id   UUID NOT NULL,
    score        DOUBLE PRECISION NOT NULL,
    grade        VARCHAR(4) NOT NULL,
    cpu_eff      DOUBLE PRECISION,
    mem_eff      DOUBLE PRECISION,
    node_packing DOUBLE PRECISION,
    hpa_coverage DOUBLE PRECISION,
    limit_comp   DOUBLE PRECISION,
    -- Contexte
    pod_count    INTEGER,
    node_count   INTEGER,
    waste_eur_annual DOUBLE PRECISION,
    co2_kg_annual    DOUBLE PRECISION
);

SELECT create_hypertable('green_scores', 'time', chunk_time_interval => INTERVAL '7 days', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS idx_scores_cluster ON green_scores(cluster_id, time DESC);

-- Aggregate journalier (continuous aggregate TimescaleDB)
CREATE MATERIALIZED VIEW IF NOT EXISTS daily_scores
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 day', time) AS day,
    cluster_id,
    AVG(score)  AS avg_score,
    MIN(score)  AS min_score,
    MAX(score)  AS max_score,
    last(grade, time) AS last_grade
FROM green_scores
GROUP BY day, cluster_id;

SELECT add_continuous_aggregate_policy('daily_scores',
    start_offset => INTERVAL '3 days',
    end_offset   => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour');
