package v1controllers

import (
	"github.com/gemyago/atlacp/internal/di"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return di.ProvideAll(container,
		newEchoController,
		di.ProvideValue(&HealthController{}),
	)
}
