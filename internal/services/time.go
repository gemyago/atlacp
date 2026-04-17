package services

import "time"

type TimeProvider interface {
	Now() time.Time
}

type timeProviderFn func() time.Time

func (fn timeProviderFn) Now() time.Time {
	return fn()
}

func NewTimeProvider() TimeProvider { //nolint:ireturn // clock port; concrete type is unexported func wrapper
	return timeProviderFn(time.Now)
}
