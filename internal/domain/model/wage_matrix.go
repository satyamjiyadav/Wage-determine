package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type WageTier struct {
	Hourly decimal.Decimal `json:"hourly"`
	Annual decimal.Decimal `json:"annual"`
}

type WageMatrix struct {
	ID               string    `json:"id"`
	DatasetVersionID string    `json:"dataset_version_id"`
	SOCCode          string    `json:"soc_code"`
	SOCTitle         string    `json:"soc_title,omitempty"`
	AreaCode         string    `json:"area_code"`
	AreaTitle        string    `json:"area_title,omitempty"`
	ProgramYear      string    `json:"program_year,omitempty"`
	Level1           WageTier  `json:"level_1"`
	Level2           WageTier  `json:"level_2"`
	Level3           WageTier  `json:"level_3"`
	Level4           WageTier  `json:"level_4"`
	Mean             WageTier  `json:"mean"`
	GeoLevel         string    `json:"geo_level"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type WageLookupRequest struct {
	SOCCode     string `json:"soc_code"`
	ZIPCode     string `json:"zip_code,omitempty"`
	FIPSCode    string `json:"fips_code,omitempty"`
	AreaCode    string `json:"area_code,omitempty"`
	ProgramYear string `json:"program_year,omitempty"`
}
