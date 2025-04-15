package queries

import (
	"context"
	"errors"
	"time"

	"github.com/maximmihin/as25/internal/dal/queries/models"
)

//go:generate go tool sqlc generate

type FtQueries struct {
	*models.Queries
}

var ErrPageLessOne = errors.New("page must be more then 1")
var ErrPageSizeBadDiaposone = errors.New("pageSize must be between 1 and 30")

func (q FtQueries) GetPvz(ctx context.Context, start, end time.Time, pageSize, page int32) ([]models.GetDateFilteredPVZWithReceptionsAndProductsRow, error) {

	var errs []error
	if page < 1 {
		errs = append(errs, ErrPageLessOne)
	}
	if pageSize < 1 || pageSize > 30 {
		errs = append(errs, ErrPageSizeBadDiaposone)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	offset := (page - 1) * pageSize
	limit := offset + pageSize

	res, err := q.GetDateFilteredPVZWithReceptionsAndProducts(ctx, models.GetDateFilteredPVZWithReceptionsAndProductsParams{
		FtOffset:  offset,
		FtLimit:   limit,
		StartDate: start,
		EndDate:   end,
	})
	if err == nil {
		if len(res) == 0 {
			return []models.GetDateFilteredPVZWithReceptionsAndProductsRow{}, nil
		}
		return res, nil
	}
	return nil, err
}
