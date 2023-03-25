//go:build wireinject
// +build wireinject

package di

import (
	"game-server-example/cmd/di/provider"
	"github.com/google/wire"
)

func Inject() (*provider.GameServer, func(), error) {
	wire.Build(
		provider.NewGameServer,
		provider.HandlerSet,
		provider.UsecaseSet,
		provider.ServiceSet,
		provider.RepositorySet,
		provider.ConfigSet,
		provider.InfraSet,
	)
	return nil, nil, nil
}
