package tolken_entity

import (
	"context"
	"os"

	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/entity/policy_entity"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/internal_error"
	"github.com/go-chi/jwtauth"
)

type Tolken struct {
	JWTSecret string
	TokenAuth *jwtauth.JWTAuth
}

func NewTolken() *Tolken {
	secret := os.Getenv("JWT_SECRET")

	return &Tolken{
		JWTSecret: secret,
		TokenAuth: jwtauth.New("HS256", []byte(secret), nil),
	}
}

func (t *Tolken) GetTolkenString() (string, *internal_error.InternalError) {
	_, retorno, err := t.TokenAuth.Encode(map[string]interface{}{
		"sub": 1,
	})

	if err != nil {
		return "", internal_error.NewInternalServerError("Error get tolken string")
	}

	return retorno, nil
}

type TolkenRepositoryInterface interface {
	Save(ctx context.Context, tolkenID string, policy *policy_entity.Policy) *internal_error.InternalError
	FindPolicyByTolken(ctx context.Context, tolkenID string) (*policy_entity.Policy, *internal_error.InternalError)
	DeleteInfoByTolken(ctx context.Context, tolkenID string) *internal_error.InternalError
}
