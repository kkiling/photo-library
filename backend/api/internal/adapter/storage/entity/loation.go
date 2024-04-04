package entity

import (
	"time"

	"github.com/google/uuid"
)

type Location struct {
	PhotoID          uuid.UUID
	CreatedAt        time.Time
	Latitude         float64
	Longitude        float64
	FormattedAddress string
	Street           string
	HouseNumber      string
	Suburb           string
	Postcode         string
	State            string
	StateCode        string
	StateDistrict    string
	County           string
	Country          string
	CountryCode      string
	City             string
}
