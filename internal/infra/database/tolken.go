package database

import (
	"context"

	"github.com/Higor-ViniciusDev/posgo_raterlimite/configuration/logger"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/entity/policy_entity"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/internal_error"
	"github.com/redis/go-redis/v9"
)

type TolkenDB struct {
	RediDB *redis.Client
}

func NewTolkenDB(redisCli *redis.Client) *TolkenDB {
	return &TolkenDB{
		RediDB: redisCli,
	}
}

func (tb *TolkenDB) Save(ctx context.Context, tolkenID string, policy *policy_entity.Policy) *internal_error.InternalError {
	// Converter struct para map (Redis não aceita struct)
	data := map[string]interface{}{
		"limit":    policy.RequestPerSecond,
		"window":   policy.WindowPerSecond,
		"ttl":      policy.TTL,
		"start_at": policy.GetTimeStartad(),
	}

	cmd := tb.RediDB.HSet(ctx, tolkenID, data)
	if cmd.Err() != nil {
		logger.Error("error save tolken redis", cmd.Err())
		return internal_error.NewInternalServerError("error save tolken redis")
	}

	return nil
}

func (tb *TolkenDB) FindPolicyByTolken(ctx context.Context, tolkenID string) (*policy_entity.Policy, *internal_error.InternalError) {
	return nil, nil
}

func (tb *TolkenDB) DeleteInfoByTolken(ctx context.Context, tolkenID string) *internal_error.InternalError {
	return nil
}
