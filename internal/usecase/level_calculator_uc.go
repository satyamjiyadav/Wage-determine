package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"prevailing-wage-service/internal/domain/model"
	"prevailing-wage-service/internal/domain/service"
)

type LevelCalculatorUseCase struct {
	wageRepo       service.IWageRepository
	geoRepo        service.IGeoRepository
	onetRepo       service.IONetRepository
	datasetRepo    service.IDatasetRepository
	logRepo        service.IDeterminationLogRepository
	calculator     *service.PointCalculator
}

func NewLevelCalculatorUseCase(
	wageRepo service.IWageRepository,
	geoRepo service.IGeoRepository,
	onetRepo service.IONetRepository,
	datasetRepo service.IDatasetRepository,
	logRepo service.IDeterminationLogRepository,
) *LevelCalculatorUseCase {
	return &LevelCalculatorUseCase{
		wageRepo:    wageRepo,
		geoRepo:     geoRepo,
		onetRepo:    onetRepo,
		datasetRepo: datasetRepo,
		logRepo:     logRepo,
		calculator:  service.NewPointCalculator(),
	}
}

func (uc *LevelCalculatorUseCase) DetermineWageLevel(ctx context.Context, payload model.JobRequirementPayload) (*model.DeterminationResult, error) {
	// 1. Get active dataset version
	ds, err := uc.datasetRepo.GetActiveDatasetVersion(ctx, payload.ProgramYear)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve dataset version: %w", err)
	}

	// 2. Fetch O*NET metadata for SOC Code
	onet, err := uc.onetRepo.GetOccupationDetails(ctx, payload.SOCCode)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch O*NET details for SOC '%s': %w", payload.SOCCode, err)
	}

	// 3. Resolve location
	locRes, err := uc.geoRepo.ResolveLocation(ctx, payload.ZIPCode, payload.FIPSCode, payload.AreaCode)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve location: %w", err)
	}

	// 4. Run 4-Tier Point Assessment Rules Engine
	rationale := uc.calculator.EvaluateLevel(payload, onet)

	// 5. Fetch Wage Matrix for the resolved location
	wageMatrix, err := uc.wageRepo.GetWage(ctx, ds.ID, payload.SOCCode, locRes.BLSArea.AreaCode)
	if err != nil {
		return nil, fmt.Errorf("wage lookup failed: %w", err)
	}

	// Select assigned tier wage rate
	var assignedWage model.WageTier
	switch rationale.FinalLevel {
	case 1:
		assignedWage = wageMatrix.Level1
	case 2:
		assignedWage = wageMatrix.Level2
	case 3:
		assignedWage = wageMatrix.Level3
	case 4:
		assignedWage = wageMatrix.Level4
	default:
		assignedWage = wageMatrix.Level1
	}

	// Generate Unique Determination Tracking Number (e.g. PWD-2026-874312)
	trackingNumber := fmt.Sprintf("PWD-%d-%06d", time.Now().Year(), rand.Intn(1000000))

	result := &model.DeterminationResult{
		DeterminationNumber: trackingNumber,
		SOCCode:             payload.SOCCode,
		SOCTitle:            onet.Title,
		AreaCode:            locRes.BLSArea.AreaCode,
		AreaTitle:           locRes.BLSArea.AreaTitle,
		ProgramYear:         ds.ProgramYear,
		AssignedLevel:       rationale.FinalLevel,
		DeterminedWage:      assignedWage,
		Rationale:           rationale,
		CreatedAt:           time.Now(),
	}

	// 6. Immutably Log Determination for Audit Compliance
	log := &model.DeterminationLog{
		ID:                  trackingNumber,
		DeterminationNumber: trackingNumber,
		DatasetVersionID:    ds.ID,
		SOCCode:             payload.SOCCode,
		AreaCode:            locRes.BLSArea.AreaCode,
		InputPayload:        payload,
		Result:              *result,
		CreatedAt:           time.Now(),
	}
	_ = uc.logRepo.SaveLog(ctx, log)

	return result, nil
}
