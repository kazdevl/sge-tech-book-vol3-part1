package provider

import (
	"game-server-example/config"
	"game-server-example/pkg/ifrepository"
	"game-server-example/pkg/infra/mysql/cacherepository"
	"game-server-example/pkg/infra/mysql/integrationrepository"
	"game-server-example/pkg/infra/mysql/repository"
	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	config.IsEnableCacheRepository,
	repository.NewItemEntityRepository,
	repository.NewMonsterEntityRepository,
	repository.NewCoinEntityRepository,
	cacherepository.NewItemEntityCacheRepository,
	cacherepository.NewMonsterEntityCacheRepository,
	cacherepository.NewCoinEntityCacheRepository,

	integrationrepository.NewItemIntegrationRepository,
	integrationrepository.NewMonsterIntegrationRepository,
	integrationrepository.NewCoinIntegrationRepository,
	wire.Bind(new(ifrepository.IFItemEntityRepository), new(*integrationrepository.ItemIntegrationRepository)),
	wire.Bind(new(ifrepository.IFMonsterEntityRepository), new(*integrationrepository.MonsterIntegrationRepository)),
	wire.Bind(new(ifrepository.IFCoinEntityRepository), new(*integrationrepository.CoinIntegrationRepository)),
)
