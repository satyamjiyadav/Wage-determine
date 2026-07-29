package usecase

import (
	"context"

	"prevailing-wage-service/internal/domain/model"
)

type IWageLookupUseCase interface {
	LookupWage(ctx context.Context, socCode, zipCode, fipsCode, areaCode, programYear string) (*model.WageMatrix, error)
	BatchLookupWages(ctx context.Context, requests []model.WageLookupRequest) ([]*model.WageMatrix, error)
}

type ILevelCalculatorUseCase interface {
	DetermineWageLevel(ctx context.Context, payload model.JobRequirementPayload) (*model.DeterminationResult, error)
}

type IDiscoveryUseCase interface {
	SearchOccupations(ctx context.Context, query string, limit int) ([]*model.ONetOccupation, error)
	GetOccupationDetails(ctx context.Context, socCode string) (*model.ONetOccupation, error)
	ResolveLocation(ctx context.Context, zipCode string) (*model.LocationResolution, error)
}

type IDeterminationAuditUseCase interface {
	GetDeterminationByNumber(ctx context.Context, number string) (*model.DeterminationLog, error)
}
