package http

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
)

func (req *PostDummyLoginJSONBody) Validate() error {
	return v.ValidateStruct(req,

		v.Field(&req.Role,
			v.In(Employee, Moderator)))
}

// TODO handle internal errror validator
func (req *PostPvzJSONRequestBody) Validate() error {
	return v.ValidateStruct(req,

		v.Field(&req.City,
			v.Required, v.In(Казань, Москва, СанктПетербург)),
	)
}

func (par *GetPvzParams) Validate() error {
	return v.ValidateStruct(&par,

		// TODO check case when equals
		v.Field(&par.StartDate,
			v.Required, v.Min(par.EndDate)),

		v.Field(&par.EndDate,
			v.Required, v.Max(par.StartDate)),

		v.Field(&par.Page,
			v.Required, v.Max(0)), // TODO ???

		v.Field(&par.Limit,
			v.Required, v.Min(1), v.Max(30)),
	)
}

func (req PostProductsJSONRequestBody) Validate() error {
	return v.ValidateStruct(&req,
		v.Field(&req.Type,
			v.Required, v.In(Обувь, Одежда, Электроника)),
	)
}

//func (req *PostReceptionsJSONRequestBody) Validate() error {
//
//}
