package database

import (
	"context"

	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/entity/policy_entity"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/entity/tolken_entity"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/internal_error"
	"github.com/redis/go-redis/v9"
)

type TolkenDB struct {
	RediDB *redis.Client
}

func (tb *TolkenDB) Save(ctx context.Context, tolkenEntity *tolken_entity.Tolken, policy *policy_entity.Policy) *internal_error.InternalError {
	tolkenString, err := tolkenEntity.GetTolkenString()

	if err != nil {
		return err
	}

	tb.RediDB.HSet(ctx, tolkenString, policy)
	return nil
}
