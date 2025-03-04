-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
		gen_random_uuid(),
		NOW(),
		NOW(),
		$1,
		$2
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUserPasswordAndEmail :one
UPDATE USERS
SET 
email = $2,
hashed_password = $3,
updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserToChirpyRed :one
UPDATE USERS
SET 
is_chirpy_red = true
WHERE id = $1
RETURNING *;
