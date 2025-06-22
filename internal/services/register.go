package services

import (
	"time"

	"github.com/gemyago/atlacp/internal/di"
	"github.com/gemyago/atlacp/internal/services/bitbucket"
	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/jira"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return di.ProvideAll(container,
		NewTimeProvider,
		di.ProvideValue(time.NewTicker),
		NewShutdownHooks,
		httpservices.NewClientFactory,
		bitbucket.NewClient,
		jira.NewClient,
	)
}
