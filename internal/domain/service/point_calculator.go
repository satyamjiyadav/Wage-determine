package service

import (
	"fmt"

	"prevailing-wage-service/internal/domain/model"
)

type PointCalculator struct{}

func NewPointCalculator() *PointCalculator {
	return &PointCalculator{}
}

// EvaluateLevel calculates the prevailing wage level (I through IV) based on job requirements and O*NET metadata
func (c *PointCalculator) EvaluateLevel(payload model.JobRequirementPayload, onet *model.ONetOccupation) model.RationaleBreakdown {
	rationale := model.RationaleBreakdown{
		BaseLevel:             1,
		EducationPoints:       0,
		ExperiencePoints:      0,
		SkillsPoints:          0,
		SupervisionPoints:     0,
		TotalCalculatedPoints: 0,
		FinalLevel:            1,
		Explanation:           make([]string, 0),
	}

	// 1. Education Evaluation
	reqEduRank := model.DegreeRank(payload.Education.RequiredDegree)
	defaultEduRank := model.DegreeRank(model.EducationLevel(onet.DefaultEducationLevel))
	
	if reqEduRank > defaultEduRank {
		diff := reqEduRank - defaultEduRank
		if diff == 1 {
			rationale.EducationPoints = 1
			rationale.Explanation = append(rationale.Explanation, 
				fmt.Sprintf("+1 Level for Education: Required degree (%s) is 1 tier above O*NET default (%s)", 
					payload.Education.RequiredDegree, onet.DefaultEducationLevel))
		} else if diff >= 2 {
			rationale.EducationPoints = 2
			rationale.Explanation = append(rationale.Explanation, 
				fmt.Sprintf("+2 Levels for Education: Required degree (%s) is %d tiers above O*NET default (%s)", 
					payload.Education.RequiredDegree, diff, onet.DefaultEducationLevel))
		}
	}

	// 2. Experience Evaluation
	reqExpMonths := payload.ExperienceMonths
	maxSVPMonths := onet.SVPMaxMonths

	if reqExpMonths > maxSVPMonths {
		excessMonths := reqExpMonths - maxSVPMonths
		if excessMonths <= 24 { // Up to 2 years excess
			rationale.ExperiencePoints = 1
			rationale.Explanation = append(rationale.Explanation, 
				fmt.Sprintf("+1 Level for Experience: Required experience (%d mos) exceeds O*NET SVP upper limit (%d mos) by %d mos", 
					reqExpMonths, maxSVPMonths, excessMonths))
		} else { // More than 2 years excess
			rationale.ExperiencePoints = 2
			rationale.Explanation = append(rationale.Explanation, 
				fmt.Sprintf("+2 Levels for Experience: Required experience (%d mos) exceeds O*NET SVP upper limit (%d mos) by %d mos (>2 yrs)", 
					reqExpMonths, maxSVPMonths, excessMonths))
		}
	}

	// 3. Special Skills / Licensure Evaluation
	if len(payload.SpecialSkills) > 0 {
		rationale.SkillsPoints = 1
		rationale.Explanation = append(rationale.Explanation, 
			fmt.Sprintf("+1 Level for Special Skills: Job requires specialized skills/credentials (%v)", payload.SpecialSkills))
	}

	// 4. Supervisory Duties Evaluation
	if payload.SupervisesEmployees {
		rationale.SupervisionPoints = 1
		subordinateInfo := ""
		if payload.NumberOfSubordinates > 0 {
			subordinateInfo = fmt.Sprintf(" (%d subordinates)", payload.NumberOfSubordinates)
		}
		rationale.Explanation = append(rationale.Explanation, 
			fmt.Sprintf("+1 Level for Supervisory Duties: Position involves direct management/supervision of staff%s", subordinateInfo))
	}

	// Total points calculation
	totalPoints := rationale.EducationPoints + rationale.ExperiencePoints + rationale.SkillsPoints + rationale.SupervisionPoints
	rationale.TotalCalculatedPoints = totalPoints

	// Final Level = Base (1) + Points, clamped to Max Level 4
	finalLevel := 1 + totalPoints
	if finalLevel > 4 {
		finalLevel = 4
	}
	rationale.FinalLevel = finalLevel

	if len(rationale.Explanation) == 0 {
		rationale.Explanation = append(rationale.Explanation, "Position meets baseline Level 1 entry requirements with no point adjustments.")
	}

	return rationale
}
