-- name: Create :exec
INSERT INTO receptions (id, pvz_id, status, created_at)
VALUES (@id, @pvz_id, @status, @created_at)
;

-- name: GetActiveByPVZ :one
SELECT * FROM receptions
WHERE pvz_id = @pvz_id AND status = 'in_progress'
;

-- name: CloseActive :one
WITH last_in_progress_reception AS (
    SELECT r.id
    FROM receptions r
    WHERE r.pvz_id = @pvz_id
      AND r.status = 'in_progress'::reception_progress
    ORDER BY r.created_at DESC
    LIMIT 1
)
UPDATE receptions
SET status = 'close'::reception_progress
WHERE id IN (SELECT id FROM last_in_progress_reception)
RETURNING *
;