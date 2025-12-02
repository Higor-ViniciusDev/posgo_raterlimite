package database

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Higor-ViniciusDev/posgo_raterlimite/configuration/logger"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/entity/policy_entity"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/entity/request_entity"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/internal_error"
	"github.com/redis/go-redis/v9"
)

type RequestInfoRepository struct {
	RedisCLI *redis.Client
}

func NewRequestInfoRepository(redisCli *redis.Client) *RequestInfoRepository {
	return &RequestInfoRepository{
		RedisCLI: redisCli,
	}
}

func (rp *RequestInfoRepository) GetInfoRequestByKey(ctx context.Context, key string) (*request_entity.RequestInfo, *internal_error.InternalError) {
	redisKey := fmt.Sprintf("request:%v", key)

	result, err := rp.RedisCLI.HGetAll(ctx, redisKey).Result()
	if err != nil {
		logger.Error("error find request info redis", err)
		return nil, internal_error.NewInternalServerError("error request info token")
	}

	if len(result) == 0 {
		return nil, internal_error.NewNotFoundError("request info not found")
	}

	quantity, err1 := strconv.ParseInt(result["quantityRequest"], 10, 64)
	limited, err2 := strconv.ParseInt(result["limitedPerPolicy"], 10, 64)
	startAtint, err3 := strconv.ParseInt(result["start_at"], 10, 64)

	if err1 != nil || err2 != nil || err3 != nil {
		logger.Error("error parse request info fields", fmt.Errorf("%v | %v | %v", err1, err2, err3))
		return nil, internal_error.NewInternalServerError("error parse requestInfo fields")
	}

	// Verificar bloqueio
	blockKey := fmt.Sprintf("bloquead:%v", key)

	exists, err := rp.RedisCLI.Exists(ctx, blockKey).Result()
	if err != nil {
		logger.Error("error checking block status redis", err)
		return nil, internal_error.NewInternalServerError("error checking blocked info")
	}

	status := request_entity.Active
	if exists > 0 {
		status = request_entity.Bloqued
	}

	p := &request_entity.RequestInfo{
		QuantityRequest:         quantity,
		LimitedRequestPerPolicy: limited,
		Status:                  status,
	}
	p.SetStartAt(startAtint)

	return p, nil
}

func (rp *RequestInfoRepository) CreateRequestInfo(ctx context.Context, key string, policy *policy_entity.Policy) *internal_error.InternalError {
	requestInfo := request_entity.NewRequestEntity(policy.RequestPerSecond)
	// Converter struct para map (Redis não aceita struct)
	data := map[string]interface{}{
		"quantityRequest":  requestInfo.QuantityRequest,
		"limitedPerPolicy": policy.RequestPerSecond,
		"start_at":         requestInfo.GetTimeRequestStarted(),
	}

	cmd := rp.RedisCLI.HSet(ctx, fmt.Sprintf("request:%v", key), data)
	if cmd.Err() != nil {
		logger.Error("error save request info redis", cmd.Err())
		return internal_error.NewInternalServerError("error save request info redis")
	}

	return nil
}

func (rp *RequestInfoRepository) BloqueadRequestByKey(ctx context.Context, key string) *internal_error.InternalError {
	data := map[string]interface{}{
		"blocked": 1,
	}

	keyFormatada := fmt.Sprintf("bloquead:%v", key)
	cmd := rp.RedisCLI.HSet(ctx, keyFormatada, data)
	if cmd.Err() != nil {
		logger.Error("error save request bloquead info redis", cmd.Err())
		return internal_error.NewInternalServerError("error save request bloquead info redis")
	}

	return nil
}

func (rp *RequestInfoRepository) UpdateRequestInfo(ctx context.Context, key string, field string, value int64) *internal_error.InternalError {
	keyFormatada := fmt.Sprintf("request:%v", key)
	err := rp.RedisCLI.HSet(ctx, keyFormatada, field, value).Err()
	if err != nil {
		logger.Error("updateRequestInfo error: ", err)
		return internal_error.NewInternalServerError("error update field in redis")
	}

	return nil
}

func (rp *RequestInfoRepository) DeleteBloqueadRequestByKey(ctx context.Context, key string) *internal_error.InternalError {
	keyFormatada := fmt.Sprintf("bloquead:%v", key)
	err := rp.RedisCLI.Del(ctx, keyFormatada).Err()
	if err != nil {
		logger.Error("Error try delete hkey bloquead in redis", err)
		return internal_error.NewInternalServerError("Error try delete bloquead key")
	}

	return nil
}
