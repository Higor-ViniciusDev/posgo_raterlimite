package strategy_usecase

import (
	"context"

	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/entity/request_entity"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/internal_error"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/usecase/expire_usecase"
)

type IPStrategyUsecase struct {
	Expire expire_usecase.ExpirerInterface
}

func NewIPStrategyUsecase(expire expire_usecase.ExpirerInterface, requestInfo request_entity.RequestRepository) *IPStrategyUsecase {
	return &IPStrategyUsecase{}
}

func (ts *IPStrategyUsecase) Validate(ctx context.Context, key string) *internal_error.InternalError {
	return nil
}
