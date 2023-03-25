package handler

import (
	amiddleware "game-server-example/pkg/handler/middleware"
	"game-server-example/pkg/handler/schema"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct {
	baseHandler           *BaseHandler
	transactionMiddleware *amiddleware.TransactionMiddleware
}

func NewRouter(baseHandler *BaseHandler, transactionMiddleware *amiddleware.TransactionMiddleware) *Router {
	return &Router{baseHandler: baseHandler, transactionMiddleware: transactionMiddleware}
}

func (r *Router) RegisterRoute(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(r.transactionMiddleware.Intercept)
	schema.RegisterHandlers(e, r.baseHandler)
}
