CREATE TABLE IF NOT EXISTS public.mock_responses (
    id SERIAL PRIMARY KEY,
    path VARCHAR(500) NOT NULL,
    method VARCHAR(10) NOT NULL,
    request_body JSONB,
    response_body JSONB NOT NULL,
    response_status_code INTEGER DEFAULT 200,
    headers TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_mock_responses_lookup
ON public.mock_responses (path, method, (md5(request_body::jsonb::text)));
