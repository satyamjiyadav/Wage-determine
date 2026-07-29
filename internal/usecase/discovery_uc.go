package usecase

import (
	"context"

	"prevailing-wage-service/internal/domain/model"
	"prevailing-wage-service/internal/domain/service"
)

type DiscoveryUseCase struct {
	geoRepo  service.IGeoRepository
	onetRepo service.IONetRepository
}

func NewDiscoveryUseCase(geoRepo service.IGeoRepository, onetRepo service.IONetRepository) *DiscoveryUseCase {
	return &DiscoveryUseCase{
		geoRepo:  geoRepo,
		onetRepo: onetRepo,
	}
}

func (uc *DiscoveryUseCase) SearchOccupations(ctx context.Context, query string, limit int) ([]*model.ONetOccupation, error) {
	return uc.geoRepo.SearchOccupations(ctx, query, limit)
}

func (uc *DiscoveryUseCase) GetOccupationDetails(ctx context.Context, socCode string) (*model.ONetOccupation, error) {
	return uc.onetRepo.GetOccupationDetails(ctx, socCode)
}

func (uc *DiscoveryUseCase) ResolveLocation(ctx context.Context, zipCode string) (*model.LocationResolution, error) {
	return uc.geoRepo.ResolveLocation(ctx, zipCode, "", "")
}
