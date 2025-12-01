package database

import (
	"context"

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
	return nil, nil
}

func (rp *RequestInfoRepository) CreateRequestInfo(ctx context.Context, key string) *internal_error.InternalError {
	return nil
}

func (rp *RequestInfoRepository) BloqueadRequestByKey(ctx context.Context, key string) *internal_error.InternalError {

}

func (rp *RequestInfoRepository) UpdateRequestInfo(ctx context.Context, key string, field string, value int64) *internal_error.InternalError
func (rp *RequestInfoRepository) DeleteBloqueadRequestByKey(ctx context.Context, key string) *internal_error.InternalError
