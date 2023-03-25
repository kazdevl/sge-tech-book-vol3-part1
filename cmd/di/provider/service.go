package provider

import (
	"game-server-example/pkg/service"
	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(
	service.NewUserRegisterService,
	service.NewUserGetDataService,
	service.NewMonsterEnhanceService,
	service.NewCheckPerformanceService,
)
