// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const closeActive = `-- name: CloseActive :one
WITH last_in_progress_reception AS (
    SELECT r.id
    FROM receptions r
    WHERE r.pvz_id = $1
      AND r.status = 'in_progress'::reception_progress
    ORDER BY r.created_at DESC
    LIMIT 1
)
UPDATE receptions
SET status = 'close'::reception_progress
WHERE id IN (SELECT id FROM last_in_progress_reception)
RETURNING id, pvz_id, status, created_at
`

func (q *Queries) CloseActive(ctx context.Context, pvzID uuid.UUID) (Reception, error) {
	row := q.db.QueryRow(ctx, closeActive, pvzID)
	var i Reception
	err := row.Scan(
		&i.ID,
		&i.PvzID,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const create = `-- name: Create :exec
INSERT INTO receptions (id, pvz_id, status, created_at)
VALUES ($1, $2, $3, $4)
`

type CreateParams struct {
	ID        uuid.UUID
	PvzID     uuid.UUID
	Status    ReceptionProgress
	CreatedAt time.Time
}

func (q *Queries) Create(ctx context.Context, arg CreateParams) error {
	_, err := q.db.Exec(ctx, create,
		arg.ID,
		arg.PvzID,
		arg.Status,
		arg.CreatedAt,
	)
	return err
}

const getActiveByPVZ = `-- name: GetActiveByPVZ :one
SELECT id, pvz_id, status, created_at FROM receptions
WHERE pvz_id = $1 AND status = 'in_progress'
`

func (q *Queries) GetActiveByPVZ(ctx context.Context, pvzID uuid.UUID) (Reception, error) {
	row := q.db.QueryRow(ctx, getActiveByPVZ, pvzID)
	var i Reception
	err := row.Scan(
		&i.ID,
		&i.PvzID,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}
