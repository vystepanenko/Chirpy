
-- +goose Up
CREATE TABLE refresh_tokens (
    token VARCHAR(64) PRIMARY KEY,
    user_id UUID NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT refresh_tokens_user_foregin FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE 
);

-- +goose Down
DROP TABLE refresh_tokens;

