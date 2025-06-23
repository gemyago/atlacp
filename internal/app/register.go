package app

import (
	"github.com/gemyago/atlacp/internal/di"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return di.ProvideAll(container,
		NewEchoService,
		NewTimeService,
		NewMathService,
		NewBitbucketService,
	)
}
