package handler

import (
	"game-server-example/pkg/handler/api"
	"game-server-example/pkg/handler/schema"
	"github.com/labstack/echo/v4"
)

type BaseHandler struct {
	monsterHandler *api.MonsterHandler
	userHandler    *api.UserHandler
}

func NewBaseHandler(monsterHandler *api.MonsterHandler, userHandler *api.UserHandler) *BaseHandler {
	return &BaseHandler{monsterHandler: monsterHandler, userHandler: userHandler}
}

func (h *BaseHandler) MonsterEnhance(ctx echo.Context, params schema.MonsterEnhanceParams) error {
	return h.monsterHandler.Enhance(ctx, params)
}

func (h *BaseHandler) UserGetData(ctx echo.Context, params schema.UserGetDataParams) error {
	return h.userHandler.UserGetData(ctx, params)
}

func (h *BaseHandler) UserRegister(ctx echo.Context) error {
	return h.userHandler.UserRegister(ctx)
}
