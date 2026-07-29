package usecase

import (
	"context"

	"prevailing-wage-service/internal/domain/model"
	"prevailing-wage-service/internal/domain/service"
)

type DeterminationAuditUseCase struct {
	logRepo service.IDeterminationLogRepository
}

func NewDeterminationAuditUseCase(logRepo service.IDeterminationLogRepository) *DeterminationAuditUseCase {
	return &DeterminationAuditUseCase{
		logRepo: logRepo,
	}
}

func (uc *DeterminationAuditUseCase) GetDeterminationByNumber(ctx context.Context, number string) (*model.DeterminationLog, error) {
	return uc.logRepo.FindByNumber(ctx, number)
}
