// +build wireinject

package di

import (
	"context"

	"github.com/google/wire"
)

func NewApp(ctx context.Context) (*App, func(), error) {
	wire.Build(wire.Struct(new(App), "*"), ServerSet, ConfigSet)
	return nil, nil, nil
}
