-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    token, user_id, expires_at, created_at, updated_at
) VALUES (
    $1,
    $2,
    NOW() + INTERVAL '60 days',
    NOW(),
    NOW()
    )
    RETURNING *;

-- name: GetRefreshToken :one
select * from refresh_tokens rt 
where rt."token" = $1
and rt.revoked_at is null 
and rt.expires_at > now();

-- name: RevokeRefreshToken :exec
update refresh_tokens 
set revoked_at = now(), updated_at = now()
where "token" =$1;
