package http

import "time"

const (
	defaultPage  = 1
	defaultLimit = 10
)

func (par *GetPvzParams) WithDefaults() (fullPar *GetPvzParams) {
	if par != nil {
		fullPar = par
	} else {
		fullPar = new(GetPvzParams)
	}

	if fullPar.StartDate == nil {
		fullPar.StartDate = ptr(time.Time{})
	}

	if fullPar.EndDate == nil {
		fullPar.EndDate = ptr(time.Now())
	}

	if fullPar.Page == nil {
		fullPar.Page = ptr(defaultPage)
	}

	if fullPar.Limit == nil {
		fullPar.Limit = ptr(defaultLimit)
	}

	return fullPar

}

func ptr[T any](v T) *T {
	return &v
}
