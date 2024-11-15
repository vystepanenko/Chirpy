-- name: CreateChirps :one
INSERT INTO chirps (
   id, user_id, body ,created_at, updated_at
) VALUES ( 
    gen_random_uuid(),
    $1,
    $2,
    NOW(),
    NOW()
)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps c
order by c.created_at ASC;

-- name: GetAllChirpsByAuthor :many
SELECT * FROM chirps c
where c.user_id = $1
order by c.created_at ASC;

-- name: GetChirp :one
select * from chirps c
where c.id = $1;

-- name: DeleteChirp :exec
delete from chirps c 
where c.id = $1
and c.user_id = $2;
