package types

import (
	"database/sql/driver"
	"fmt"
)

type FtCity string

const (
	CityMsk FtCity = "Москва"
	CityKzn FtCity = "Казань"
	CitySpb FtCity = "СанктПетербург" // TODO
)

func (e *FtCity) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = FtCity(s)
	case string:
		*e = FtCity(s)
	default:
		return fmt.Errorf("unsupported scan type for FtCity: %T", src)
	}
	return nil
}

type NullFtCity struct {
	City  FtCity
	Valid bool // Valid is true if City is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullFtCity) Scan(value interface{}) error {
	if value == nil {
		ns.City, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.City.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullFtCity) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.City), nil
}
