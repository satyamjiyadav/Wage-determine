package service_test

import (
	"testing"

	"prevailing-wage-service/internal/domain/model"
	"prevailing-wage-service/internal/domain/service"
)

func TestPointCalculator_EvaluateLevel(t *testing.T) {
	calc := service.NewPointCalculator()

	onetDev := &model.ONetOccupation{
		SOCCode:               "15-1252.00",
		Title:                 "Software Developers",
		JobZone:               4,
		SVPMinMonths:          24,
		SVPMaxMonths:          48,
		DefaultEducationLevel: "Bachelor",
	}

	tests := []struct {
		name          string
		payload       model.JobRequirementPayload
		expectedLevel int
	}{
		{
			name: "Entry level candidate - Bachelor degree & 1 year exp -> Level 1",
			payload: model.JobRequirementPayload{
				SOCCode:             "15-1252.00",
				Education:           model.EducationRequirement{RequiredDegree: model.DegreeBachelor},
				ExperienceMonths:    12,
				SupervisesEmployees: false,
			},
			expectedLevel: 1,
		},
		{
			name: "Qualified candidate - Master degree (+1) & 3 years exp -> Level 2",
			payload: model.JobRequirementPayload{
				SOCCode:             "15-1252.00",
				Education:           model.EducationRequirement{RequiredDegree: model.DegreeMaster},
				ExperienceMonths:    36,
				SupervisesEmployees: false,
			},
			expectedLevel: 2,
		},
		{
			name: "Experienced candidate - Master (+1) & 5 years exp (+1) & Special skills (+1) -> Level 4",
			payload: model.JobRequirementPayload{
				SOCCode:             "15-1252.00",
				Education:           model.EducationRequirement{RequiredDegree: model.DegreeMaster},
				ExperienceMonths:    60,
				SpecialSkills:       []string{"Go", "Kubernetes"},
				SupervisesEmployees: false,
			},
			expectedLevel: 4,
		},
		{
			name: "Lead candidate - Master (+1) & 6 years exp (>2 yrs excess = +2) & Supervision (+1) -> Level 4 (clamped)",
			payload: model.JobRequirementPayload{
				SOCCode:              "15-1252.00",
				Education:            model.EducationRequirement{RequiredDegree: model.DegreeMaster},
				ExperienceMonths:     78, // > 48 + 24
				SupervisesEmployees:  true,
				NumberOfSubordinates: 5,
			},
			expectedLevel: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rationale := calc.EvaluateLevel(tt.payload, onetDev)
			if rationale.FinalLevel != tt.expectedLevel {
				t.Errorf("Expected Level %d, got Level %d. Rationale: %+v", tt.expectedLevel, rationale.FinalLevel, rationale.Explanation)
			}
		})
	}
}
