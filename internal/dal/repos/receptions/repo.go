package pvz

import (
	"context"
	"errors"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maximmihin/as25/internal/dal/repos/receptions/models"
)

//go:generate go tool sqlc generate

type ExtDBTX interface {
	models.DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Repo struct {
	ExtDBTX
	queries *models.Queries
	// TODO add timeouts
}

func New(pool *pgxpool.Pool) *Repo {
	return &Repo{
		ExtDBTX: pool,
		queries: models.New(pool),
	}
}

func (r Repo) WithTx(tx pgx.Tx) *Repo {
	return &Repo{
		ExtDBTX: tx,
		queries: r.queries.WithTx(tx),
	}
}

var ErrNonexistentPvzId = errors.New("nonexistent Pvz ID")
var ErrThereAreOpenReceptions = errors.New("there are open receptions")

func (r Repo) Create(ctx context.Context, pvz models.Reception) error {
	err := r.queries.Create(ctx, models.CreateParams(pvz))
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.ConstraintName {
		case "receptions_fk_pvz_id":
			return ErrNonexistentPvzId
		case "receptions_idx_one_in_progress_per_pvz":
			return ErrThereAreOpenReceptions
		}
	}

	return err
}

func (r Repo) GetActiveByPVZ(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error) {
	rec, err := r.queries.GetActiveByPVZ(ctx, pvzId)
	if err == nil {
		return &rec, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

var ErrNothingToClose = errors.New("nothing to close")

func (r Repo) CloseActive(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error) {
	rec, err := r.queries.CloseActive(ctx, pvzId)
	if err == nil {
		return &rec, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNothingToClose
	}
	return nil, err
}
