-- +goose Up
CREATE TABLE chirps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    body TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT chirps_user_foregin FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE 
);

-- +goose Down
DROP TABLE chirps;

