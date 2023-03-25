package provider

import (
	"game-server-example/config"
	"github.com/google/wire"
)

var ConfigSet = wire.NewSet(
	config.NewMysqlConfig,
)
