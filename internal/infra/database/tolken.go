package database

import (
	"context"
	"fmt"
	"strconv"

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
	result, err := tb.RediDB.HGetAll(ctx, tolkenID).Result()
	if err != nil {
		logger.Error("error find tolken redis", err)
		return nil, internal_error.NewInternalServerError("error search tolken")
	}

	if len(result) == 0 {
		return nil, internal_error.NewNotFoundError("tolken not found")
	}

	limitStr := result["limit"]
	windowStr := result["window"]
	ttlStr := result["ttl"]
	startStr := result["start_at"]

	limit, err1 := strconv.ParseInt(limitStr, 10, 64)
	window, err2 := strconv.ParseInt(windowStr, 10, 64)
	ttl, err3 := strconv.ParseInt(ttlStr, 10, 64)
	startAt, err4 := strconv.ParseInt(startStr, 10, 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		logger.Error("error parse tolken fields", fmt.Errorf("%v %v %v %v", err1, err2, err3, err4))
		return nil, internal_error.NewInternalServerError("error parse tolken fields")
	}

	p := &policy_entity.Policy{
		RequestPerSecond: limit,
		WindowPerSecond:  window,
		TTL:              ttl,
		Fonte:            policy_entity.FONTE_TOLKEN,
	}
	p.SetStartAt(startAt)

	return p, nil
}

func (tb *TolkenDB) DeleteInfoByTolken(ctx context.Context, tolkenID string) *internal_error.InternalError {
	err := tb.RediDB.Del(ctx, tolkenID).Err()
	if err != nil {
		logger.Error("Error try delete hkey in redis", err)
		return internal_error.NewInternalServerError("Error try delete key")
	}

	return nil
}
