CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'want to read'
        CHECK (status IN ('reading', 'finished', 'want to read'))
);
