package provider

import (
	"game-server-example/pkg/usecase"
	"github.com/google/wire"
)

var UsecaseSet = wire.NewSet(
	usecase.NewUserUsecase,
	usecase.NewMonsterUsecase,
)
