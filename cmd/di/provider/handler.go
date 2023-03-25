package provider

import (
	"game-server-example/pkg/handler"
	"game-server-example/pkg/handler/api"
	"game-server-example/pkg/handler/middleware"
	"github.com/google/wire"
)

var HandlerSet = wire.NewSet(
	handler.NewBaseHandler,
	handler.NewRouter,
	middleware.NewTransactionMiddleware,
	api.NewUserHandler,
	api.NewMonsterHandler,
)
