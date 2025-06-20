package services

import (
	"time"

	"github.com/gemyago/atlacp/internal/di"
	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return di.ProvideAll(container,
		NewTimeProvider,
		di.ProvideValue(time.NewTicker),
		NewShutdownHooks,
		httpservices.NewClientFactory,
	)
}
