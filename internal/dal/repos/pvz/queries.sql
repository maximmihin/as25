-- name: Create :exec
INSERT INTO pvz (id, created_at, located_at)
VALUES (@id, @created_at, @located_at)
;
