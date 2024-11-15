
-- +goose Up
ALTER TABLE users 
ADD column is_chirpy_red bool DEFAULT false NOT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN is_chirpy_red;
