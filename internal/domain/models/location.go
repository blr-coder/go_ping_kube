package models

type Location struct {
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	City        string `json:"city"`
	Region      string `json:"region"`
	Err         error  `json:"-"`
}
