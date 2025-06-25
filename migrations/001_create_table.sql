CREATE TABLE IF NOT EXISTS urls (
    id TEXT PRIMARY KEY,
    original TEXT NOT NULL,
    short_code TEXT UNIQUE NOT NULL,
    custom_code TEXT,
    domain TEXT,
    visibility TEXT,
    created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS visits (
    id SERIAL PRIMARY KEY,
    code TEXT,
    timestamp TIMESTAMP,
    ip TEXT,
    country TEXT,
    browser TEXT,
    device TEXT,
    FOREIGN KEY (code) REFERENCES urls(id)
);
