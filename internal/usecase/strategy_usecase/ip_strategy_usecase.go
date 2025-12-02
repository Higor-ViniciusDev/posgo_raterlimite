package strategy_usecase

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/Higor-ViniciusDev/posgo_raterlimite/configuration/logger"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/entity/policy_entity"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/entity/request_entity"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/entity/tolken_entity"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/internal_error"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/usecase/expire_usecase"
)

type IPStrategyUsecase struct {
	Expirer          expire_usecase.ExpirerInterface
	TolkenRepository tolken_entity.TolkenRepositoryInterface
	RequestInfo      request_entity.RequestRepository
}

func NewIPStrategyUsecase(expire expire_usecase.ExpirerInterface, requestInfo request_entity.RequestRepository) *IPStrategyUsecase {
	return &IPStrategyUsecase{}
}

func (ts *IPStrategyUsecase) Validate(ctx context.Context, key string) *internal_error.InternalError {
	policyTolken, err := ts.TolkenRepository.FindPolicyByTolken(ctx, key)
	if err != nil {
		if err.Err == "not_found" {
			policyTolken = policy_entity.NewPolicyIP()
			errCreate := ts.TolkenRepository.Save(ctx, key, policyTolken)

			if errCreate != nil {
				return errCreate
			}
		} else {
			return err
		}
	}

	infoRequest, err := ts.RequestInfo.GetInfoRequestByKey(ctx, key)
	if err != nil && err.Err != "not_found" {
		return err
	}

	if infoRequest == nil {
		if err := ts.RequestInfo.CreateRequestInfo(ctx, key, policyTolken); err != nil {
			logger.Error("Error creating request info", err)
		}
		return nil
	}

	if infoRequest.Status == request_entity.Bloqued {
		return internal_error.NewManyRequestError("too many requests for this IP")
	}

	if infoRequest.QuantityRequest >= infoRequest.LimitedRequestPerPolicy {

		err := ts.RequestInfo.BloqueadRequestByKey(ctx, key)
		if err != nil {
			logger.Error("Error blocking request", err)
			return err
		}

		unlockWindow := os.Getenv("REQUEST_PER_WINDOW")
		windowSeconds, _ := strconv.Atoi(unlockWindow)

		ts.Expirer.SetExpiration(
			key,
			time.Duration(windowSeconds)*time.Second,
			func() {
				ts.RequestInfo.DeleteBloqueadRequestByKey(context.Background(), key)
				ts.RequestInfo.UpdateRequestInfo(context.Background(), key, "quantityRequest", 0)
			},
		)

		return internal_error.NewManyRequestError("too many requests for this IP")
	}

	_ = ts.RequestInfo.UpdateRequestInfo(ctx, key, "quantityRequest", infoRequest.QuantityRequest+1)

	return nil
}
