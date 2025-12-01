package expire_usecase

import "time"

type ExpirerInterface interface {
	SetExpiration(Key string, duration time.Duration, callback func())
}

type DefaultExpirer struct{}

func NewDefaultExpirer() *DefaultExpirer {
	return &DefaultExpirer{}
}

func (e *DefaultExpirer) SetExpiration(key string, duration time.Duration, callback func()) {
	time.AfterFunc(duration, callback)
}
