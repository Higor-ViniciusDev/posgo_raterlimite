package request_entity

import (
	"context"
	"time"

	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/entity/policy_entity"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/internal_error"
)

type RequestInfo struct {
	timeRequestStarted      int64
	QuantityRequest         int64
	LimitedRequestPerPolicy int64
	Status                  RequestCondition
}

type RequestCondition int64

const (
	Active RequestCondition = iota
	Bloqued
)

func (r *RequestInfo) GetTimeRequestStarted() int64 {
	return r.timeRequestStarted
}

// SetStartAt define o tempo inicial (usado ao ler do armazenamento)
func (p *RequestInfo) SetStartAt(t int64) {
	p.timeRequestStarted = t
}

func NewRequestEntity(requestPerSecond int64) *RequestInfo {
	// limetedRequest
	return &RequestInfo{
		timeRequestStarted:      time.Now().Unix(),
		QuantityRequest:         0,
		LimitedRequestPerPolicy: requestPerSecond,
		Status:                  Active,
	}
}

type RequestRepository interface {
	GetInfoRequestByKey(ctx context.Context, key string) (*RequestInfo, *internal_error.InternalError)
	CreateRequestInfo(ctx context.Context, key string, policy *policy_entity.Policy) *internal_error.InternalError
	BloqueadRequestByKey(ctx context.Context, key string) *internal_error.InternalError
	UpdateRequestInfo(ctx context.Context, key string, field string, value int64) *internal_error.InternalError
	DeleteBloqueadRequestByKey(ctx context.Context, key string) *internal_error.InternalError
}
