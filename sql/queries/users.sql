-- name: CreateUser :one
INSERT INTO users (
   id, email, hashed_password, created_at, updated_at
) VALUES ( 
    gen_random_uuid(),
    $1,
    $2,
    NOW(),
    NOW()
) returning *;

-- name: DeleteAllUsers :exec
DELETE FROM USERS;

-- name: GetUserByEmail :one
SELECT * FROM users u
where u.email = $1; 
