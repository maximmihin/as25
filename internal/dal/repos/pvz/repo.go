package pvz

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maximmihin/as25/internal/dal/repos/pvz/models"
)

//go:generate go tool sqlc generate

type ExtDBTX interface {
	models.DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Repo struct {
	pgxPool ExtDBTX
	queries *models.Queries
	// TODO add timeouts
}

func New(pool *pgxpool.Pool) *Repo {
	return &Repo{
		pgxPool: pool,
		queries: models.New(pool),
	}
}

func (r Repo) WithTx(tx pgx.Tx) *Repo {
	return &Repo{
		pgxPool: tx,
		queries: r.queries.WithTx(tx),
	}
}

var ErrUnavailableCityType = errors.New("unavailable city type")

func (r Repo) Create(ctx context.Context, pvz models.Pvz) error {
	err := r.queries.Create(ctx, models.CreateParams(pvz))
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// TODO check how look constraint error
		return ErrUnavailableCityType
	}

	return err
}
