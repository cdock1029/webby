-- name: GetProperty :one
SELECT *
from properties
WHERE id = $1
LIMIT 1;
-- name: ListProperties :many
SELECT *
from properties
order by name;
-- name: CreateProperty :one
INSERT INTO properties (name)
VALUES($1)
RETURNING *;
-- name: DeleteProperty :exec
DELETE from properties
WHERE id = $1;
-- name: UpdateProperty :one
UPDATE properties
SET name = $2
WHERE id = $1
RETURNING *;