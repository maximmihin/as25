// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/maximmihin/as25/internal/dal/types"
)

const getDateFilteredPVZWithReceptionsAndProducts = `-- name: GetDateFilteredPVZWithReceptionsAndProducts :many
WITH filtered_receptions AS (
    SELECT
        r.id,
        r.created_at,
        r.pvz_id,
        r.status
    FROM
        receptions r
    WHERE
        r.created_at BETWEEN $3 AND $4
)
SELECT
    p.id AS pvz_id,
    p.created_at AS pvz_created_at,
    p.located_at AS pvz_city,
    COALESCE(
        json_agg(
            json_build_object(
                'reception_id', fr.id,
                'reception_created_at', fr.created_at,
                'reception_status', fr.status
            ) ORDER BY fr.created_at
        ) FILTER (WHERE fr.id IS NOT NULL),
        '[]'::json
    ) AS receptions
FROM pvz p
    LEFT JOIN filtered_receptions fr ON p.id = fr.pvz_id
GROUP BY
    p.id
ORDER BY
    p.created_at
LIMIT $2
OFFSET $1
`

type GetDateFilteredPVZWithReceptionsAndProductsParams struct {
	FtOffset  int32
	FtLimit   int32
	StartDate time.Time
	EndDate   time.Time
}

type GetDateFilteredPVZWithReceptionsAndProductsRow struct {
	PvzID        uuid.UUID
	PvzCreatedAt time.Time
	PvzCity      types.FtCity
	Receptions   interface{}
}

func (q *Queries) GetDateFilteredPVZWithReceptionsAndProducts(ctx context.Context, arg GetDateFilteredPVZWithReceptionsAndProductsParams) ([]GetDateFilteredPVZWithReceptionsAndProductsRow, error) {
	rows, err := q.db.Query(ctx, getDateFilteredPVZWithReceptionsAndProducts,
		arg.FtOffset,
		arg.FtLimit,
		arg.StartDate,
		arg.EndDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDateFilteredPVZWithReceptionsAndProductsRow
	for rows.Next() {
		var i GetDateFilteredPVZWithReceptionsAndProductsRow
		if err := rows.Scan(
			&i.PvzID,
			&i.PvzCreatedAt,
			&i.PvzCity,
			&i.Receptions,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
