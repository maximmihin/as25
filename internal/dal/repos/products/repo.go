package pvz

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maximmihin/as25/internal/dal/repos/products/models"
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

// TODO rm ctx
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

var ErrNonexistentReceptionId = errors.New("nonexistent reception id")

func (r Repo) Create(ctx context.Context, product models.Product) error {
	err := r.queries.Create(ctx, models.CreateParams(product))
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.ConstraintName {
		case "products_fk_reception_id":
			return ErrNonexistentReceptionId
		}
	}
	return err
}

var ErrNothingToDelete = errors.New("nothing to delete")

func (r Repo) DeleteLastInPvz(ctx context.Context, pvzId uuid.UUID) error {
	ra, err := r.queries.DeleteLastInPvz(ctx, pvzId)
	if err != nil {
		return err
	}
	if ra == 0 {
		return ErrNothingToDelete
	}
	return nil
}
