package model

type ONetOccupation struct {
	SOCCode                string   `json:"soc_code"`
	Title                  string   `json:"title"`
	Description            string   `json:"description"`
	JobZone                int      `json:"job_zone"`
	SVPMinMonths           int      `json:"svp_min_months"`
	SVPMaxMonths           int      `json:"svp_max_months"`
	DefaultEducationLevel  string   `json:"default_education_level"`
	SampleTitles           []string `json:"sample_titles"`
}
