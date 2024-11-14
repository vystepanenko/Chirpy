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


-- name: GetChirp :one
select * from chirps c
where c.id = $1;
