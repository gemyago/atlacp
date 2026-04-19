package services

import "time"

type TimeProvider interface {
	Now() time.Time
}

type timeProviderFn func() time.Time

func (fn timeProviderFn) Now() time.Time {
	return fn()
}

//nolint:ireturn // clock port; concrete type is unexported func wrapper
func NewTimeProvider() TimeProvider {
	return timeProviderFn(time.Now)
}
