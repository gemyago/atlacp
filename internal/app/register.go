package app

import (
	"github.com/gemyago/atlacp/internal/di"
	"github.com/gemyago/atlacp/internal/services/bitbucket"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return di.ProvideAll(container,
		NewEchoService,
		NewBitbucketService,
		newBitbucketAuthFactory,
		di.ProvideAs[*bitbucket.Client, bitbucketClient],
	)
}
