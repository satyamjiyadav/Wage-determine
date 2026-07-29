package model

import (
	"time"
)

type EducationLevel string

const (
	DegreeHighSchool EducationLevel = "HighSchool"
	DegreeAssociate  EducationLevel = "Associate"
	DegreeBachelor   EducationLevel = "Bachelor"
	DegreeMaster     EducationLevel = "Master"
	DegreeDoctorate  EducationLevel = "Doctorate"
)

// DegreeRank maps education levels to ordinal values for evaluation
func DegreeRank(level EducationLevel) int {
	switch level {
	case DegreeHighSchool:
		return 1
	case DegreeAssociate:
		return 2
	case DegreeBachelor:
		return 3
	case DegreeMaster:
		return 4
	case DegreeDoctorate:
		return 5
	default:
		return 1
	}
}

type EducationRequirement struct {
	RequiredDegree EducationLevel `json:"required_degree"`
	FieldOfStudy   string         `json:"field_of_study,omitempty"`
}

type JobRequirementPayload struct {
	SOCCode              string               `json:"soc_code"`
	ZIPCode              string               `json:"zip_code,omitempty"`
	FIPSCode             string               `json:"fips_code,omitempty"`
	AreaCode             string               `json:"area_code,omitempty"`
	JobTitle             string               `json:"job_title"`
	Education            EducationRequirement `json:"education"`
	ExperienceMonths     int                  `json:"experience_months"`
	SpecialSkills        []string             `json:"special_skills,omitempty"`
	SupervisesEmployees  bool                 `json:"supervises_employees"`
	NumberOfSubordinates int                  `json:"number_of_subordinates,omitempty"`
	ProgramYear          string               `json:"program_year,omitempty"`
}

type RationaleBreakdown struct {
	BaseLevel             int      `json:"base_level"`
	EducationPoints       int      `json:"education_points"`
	ExperiencePoints      int      `json:"experience_points"`
	SkillsPoints          int      `json:"skills_points"`
	SupervisionPoints     int      `json:"supervision_points"`
	TotalCalculatedPoints int      `json:"total_calculated_points"`
	FinalLevel            int      `json:"final_level"`
	Explanation           []string `json:"explanation"`
}

type DeterminationResult struct {
	DeterminationNumber string             `json:"determination_number"`
	SOCCode             string             `json:"soc_code"`
	SOCTitle            string             `json:"soc_title"`
	AreaCode            string             `json:"area_code"`
	AreaTitle           string             `json:"area_title"`
	ProgramYear         string             `json:"program_year"`
	AssignedLevel       int                `json:"assigned_level"`
	DeterminedWage      WageTier           `json:"determined_wage"`
	Rationale           RationaleBreakdown `json:"rationale_breakdown"`
	CreatedAt           time.Time          `json:"created_at"`
}

type DeterminationLog struct {
	ID                  string                `json:"id"`
	DeterminationNumber string                `json:"determination_number"`
	ClientID            string                `json:"client_id,omitempty"`
	DatasetVersionID    string                `json:"dataset_version_id"`
	SOCCode             string                `json:"soc_code"`
	AreaCode            string                `json:"area_code"`
	InputPayload        JobRequirementPayload `json:"input_payload"`
	Result              DeterminationResult   `json:"result"`
	CreatedAt           time.Time             `json:"created_at"`
}
