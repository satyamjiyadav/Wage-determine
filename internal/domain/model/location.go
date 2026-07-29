package model

type FIPSCounty struct {
	FIPSCode   string `json:"fips_code"`
	StateAbbr  string `json:"state_abbr"`
	CountyName string `json:"county_name"`
}

type BLSArea struct {
	AreaCode  string `json:"area_code"`
	AreaTitle string `json:"area_title"`
	AreaType  string `json:"area_type"`
	StateAbbr string `json:"state_abbr"`
}

type ZIPMapping struct {
	ZIPCode     string  `json:"zip_code"`
	FIPSCode    string  `json:"fips_code"`
	PrimaryCity string  `json:"primary_city"`
	StateAbbr   string  `json:"state_abbr"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type LocationResolution struct {
	ZIPCode    string     `json:"zip_code"`
	FIPSCounty FIPSCounty `json:"fips_county"`
	BLSArea    BLSArea    `json:"bls_area"`
}
