package controllers

import (
	"github.com/gemyago/atlacp/internal/api/mcp/server"
	"github.com/gemyago/atlacp/internal/app"
	"github.com/gemyago/atlacp/internal/di"
	"go.uber.org/dig"
)

type controllerResult struct {
	dig.Out

	Controller server.ToolsFactory `group:"mcp-controllers"`
}

func newToolsFactory[T server.ToolsFactory](controller T) controllerResult {
	return controllerResult{
		Controller: controller,
	}
}

func Register(container *dig.Container) error {
	return di.ProvideAll(container,
		NewBitbucketController,
		newToolsFactory[*BitbucketController],
		di.ProvideAs[*app.BitbucketService, bitbucketService],
	)
}
