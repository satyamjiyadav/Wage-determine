package service

import (
	"context"

	"prevailing-wage-service/internal/domain/model"
)

type IWageRepository interface {
	GetWage(ctx context.Context, datasetVersionID, socCode, areaCode string) (*model.WageMatrix, error)
	SaveWageBatch(ctx context.Context, wages []*model.WageMatrix) error
}

type IGeoRepository interface {
	ResolveLocation(ctx context.Context, zipCode, fipsCode, areaCode string) (*model.LocationResolution, error)
	GetBLSArea(ctx context.Context, areaCode string) (*model.BLSArea, error)
	SearchOccupations(ctx context.Context, query string, limit int) ([]*model.ONetOccupation, error)
}

type IONetRepository interface {
	GetOccupationDetails(ctx context.Context, socCode string) (*model.ONetOccupation, error)
}

type IDatasetRepository interface {
	GetActiveDatasetVersion(ctx context.Context, programYear string) (*model.DatasetVersion, error)
	ListDatasets(ctx context.Context) ([]*model.DatasetVersion, error)
	ActivateDataset(ctx context.Context, datasetID string) error
}

type IDeterminationLogRepository interface {
	SaveLog(ctx context.Context, log *model.DeterminationLog) error
	FindByNumber(ctx context.Context, number string) (*model.DeterminationLog, error)
}

type ICacheManager interface {
	Get(ctx context.Context, key string, target interface{}) (bool, error)
	Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error
	Invalidate(ctx context.Context, pattern string) error
}
