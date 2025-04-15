-- name: Create :exec
INSERT INTO products (id, created_at, type, reception_id)
VALUES (@id, @created_at, @type, @reception_id)
;

-- name: DeleteLastInPvz :execrows
WITH last_product AS (
    SELECT p.id, p.reception_id
    FROM products p
        JOIN receptions r ON p.reception_id = r.id
    WHERE r.pvz_id = @pvz_id
        AND r.status = 'in_progress'::reception_progress
    ORDER BY p.created_at DESC
    LIMIT 1
)
DELETE FROM products
WHERE id IN (SELECT id FROM last_product)
RETURNING *
;