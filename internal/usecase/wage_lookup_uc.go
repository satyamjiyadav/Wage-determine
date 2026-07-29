package usecase

import (
	"context"
	"fmt"

	"prevailing-wage-service/internal/domain/model"
	"prevailing-wage-service/internal/domain/service"
)

type WageLookupUseCase struct {
	wageRepo    service.IWageRepository
	geoRepo     service.IGeoRepository
	datasetRepo service.IDatasetRepository
	cache       service.ICacheManager
}

func NewWageLookupUseCase(
	wageRepo service.IWageRepository,
	geoRepo service.IGeoRepository,
	datasetRepo service.IDatasetRepository,
	cache service.ICacheManager,
) *WageLookupUseCase {
	return &WageLookupUseCase{
		wageRepo:    wageRepo,
		geoRepo:     geoRepo,
		datasetRepo: datasetRepo,
		cache:       cache,
	}
}

func (uc *WageLookupUseCase) LookupWage(ctx context.Context, socCode, zipCode, fipsCode, areaCode, programYear string) (*model.WageMatrix, error) {
	// 1. Fetch active dataset version
	ds, err := uc.datasetRepo.GetActiveDatasetVersion(ctx, programYear)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve dataset version: %w", err)
	}

	// 2. Resolve location to BLS Area Code
	locRes, err := uc.geoRepo.ResolveLocation(ctx, zipCode, fipsCode, areaCode)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve location: %w", err)
	}

	resolvedAreaCode := locRes.BLSArea.AreaCode

	// 3. Cache lookup
	cacheKey := fmt.Sprintf("wage:%s:%s:%s", ds.ID, socCode, resolvedAreaCode)
	var cachedMatrix model.WageMatrix
	if found, _ := uc.cache.Get(ctx, cacheKey, &cachedMatrix); found {
		return &cachedMatrix, nil
	}

	// 4. Database / Persistence Lookup
	wageMatrix, err := uc.wageRepo.GetWage(ctx, ds.ID, socCode, resolvedAreaCode)
	if err != nil {
		return nil, fmt.Errorf("wage lookup failed: %w", err)
	}

	// Populate additional metadata
	wageMatrix.AreaTitle = locRes.BLSArea.AreaTitle
	wageMatrix.ProgramYear = ds.ProgramYear

	// 5. Store in Cache (1 hour TTL)
	_ = uc.cache.Set(ctx, cacheKey, wageMatrix, 3600)

	return wageMatrix, nil
}

func (uc *WageLookupUseCase) BatchLookupWages(ctx context.Context, requests []model.WageLookupRequest) ([]*model.WageMatrix, error) {
	results := make([]*model.WageMatrix, len(requests))
	for i, req := range requests {
		matrix, err := uc.LookupWage(ctx, req.SOCCode, req.ZIPCode, req.FIPSCode, req.AreaCode, req.ProgramYear)
		if err != nil {
			return nil, fmt.Errorf("batch item [%d] failed for SOC '%s': %w", i, req.SOCCode, err)
		}
		results[i] = matrix
	}
	return results, nil
}
