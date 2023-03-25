package provider

import "game-server-example/pkg/handler"

type GameServer struct {
	Router *handler.Router
}

func NewGameServer(router *handler.Router) *GameServer {
	return &GameServer{Router: router}
}
