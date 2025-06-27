CREATE TABLE IF NOT EXISTS telemetry (
    timestamp DATETIME, 
    level VARCHAR(10),
    service_name VARCHAR(50),
    message TEXT,
    op VARCHAR(50)
)
ENGINE= MergeTree()
ORDER BY (service_name, timestamp);