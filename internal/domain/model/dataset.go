package model

import "time"

type DatasetVersion struct {
	ID                 string    `json:"id"`
	ProgramYear        string    `json:"program_year"`
	IsActive           bool      `json:"is_active"`
	ReleaseDate        time.Time `json:"release_date"`
	EffectiveStartDate time.Time `json:"effective_start_date"`
	EffectiveEndDate   time.Time `json:"effective_end_date"`
	Status             string    `json:"status"` // e.g. READY, INGESTING, ARCHIVED
	CreatedAt          time.Time `json:"created_at"`
}
