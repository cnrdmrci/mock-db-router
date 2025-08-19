CREATE TABLE IF NOT EXISTS mock_responses (
    id SERIAL PRIMARY KEY,
    path VARCHAR(500) NOT NULL,
    method VARCHAR(10) NOT NULL,
    response_body TEXT,
    response_status_code INTEGER DEFAULT 200,
    headers TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(path, method)
);
