package services

import "time"

type TimeProvider interface {
	Now() time.Time
}

type timeProviderFn func() time.Time

func (fn timeProviderFn) Now() time.Time {
	return fn()
}

//nolint:ireturn // TimeProvider is the clock port; implementation stays a small function type.
func NewTimeProvider() TimeProvider {
	return timeProviderFn(time.Now)
}
