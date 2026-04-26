package services

import "time"

type TimeProvider interface {
	Now() time.Time
}

type timeProviderFn func() time.Time

func (fn timeProviderFn) Now() time.Time {
	return fn()
}

//nolint:ireturn,nolintlint // clock port returns the package interface by design
func NewTimeProvider() TimeProvider {
	return timeProviderFn(time.Now)
}
