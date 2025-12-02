package strategy_usecase

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/Higor-ViniciusDev/posgo_raterlimite/configuration/logger"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/entity/request_entity"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/entity/tolken_entity"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/internal_error"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/usecase/expire_usecase"
)

type TolkenStrategyUsecase struct {
	Expirer          expire_usecase.ExpirerInterface
	TolkenRepository tolken_entity.TolkenRepositoryInterface
	RequestInfo      request_entity.RequestRepository
}

func NewTolkenStrategyUsecase(expire expire_usecase.ExpirerInterface, tolkenRepository tolken_entity.TolkenRepositoryInterface, requestInfo request_entity.RequestRepository) *TolkenStrategyUsecase {
	return &TolkenStrategyUsecase{
		Expirer:          expire,
		TolkenRepository: tolkenRepository,
		RequestInfo:      requestInfo,
	}
}

func (ts *TolkenStrategyUsecase) Validate(ctx context.Context, key string) *internal_error.InternalError {
	tolkeRequest, err := ts.TolkenRepository.FindPolicyByTolken(ctx, key)

	if err != nil {
		return err
	}

	if tolkeRequest == nil {
		return internal_error.NewBadRequestError("invalid tolken request")
	}

	infoRequest, err := ts.RequestInfo.GetInfoRequestByKey(ctx, key)

	if err != nil {
		return err
	}

	if infoRequest == nil {
		err := ts.RequestInfo.CreateRequestInfo(ctx, key, tolkeRequest)

		if err != nil {
			logger.Error("Error in created request info validation", err)
		}

		return nil
	}

	if infoRequest.Status == request_entity.Bloqued {
		return internal_error.NewManyRequestError("you have reached the maximum number of requests or actions allowed within a certain time frame")
	}

	if infoRequest.QuantityRequest >= tolkeRequest.RequestPerSecond {
		err := ts.RequestInfo.BloqueadRequestByKey(ctx, key)

		if err != nil {
			logger.Error("Error bloquead request", err)
			return err
		}
		desbloquead := os.Getenv("REQUEST_PER_WINDOW")
		desbloqueadInt, _ := strconv.Atoi(desbloquead)

		ts.Expirer.SetExpiration(
			key,
			time.Duration(desbloqueadInt)*time.Second,
			func() {
				ts.RequestInfo.DeleteBloqueadRequestByKey(context.Background(), key)
				ts.RequestInfo.UpdateRequestInfo(ctx, key, "quantityRequest", 0)
			},
		)

		return internal_error.NewManyRequestError("you have reached the maximum number of requests or actions allowed within a certain time frame")
	}
	_ = ts.RequestInfo.UpdateRequestInfo(ctx, key, "quantityRequest", infoRequest.QuantityRequest+1)

	return nil
}
