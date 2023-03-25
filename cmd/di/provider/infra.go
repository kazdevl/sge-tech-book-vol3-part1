package provider

import (
	"game-server-example/pkg/ifrepository"
	"game-server-example/pkg/infra/mysql"
	"game-server-example/pkg/infra/mysql/cachedb"
	"game-server-example/pkg/infra/mysql/repository"
	"github.com/google/wire"
)

var InfraSet = wire.NewSet(
	repository.NewUUIDGenerator,
	wire.Bind(new(ifrepository.IFUUIDGenerator), new(*repository.UUIDGenerator)),
	mysql.NewDB,
	mysql.RetrieveSqlxDB,
	cachedb.NewCacheDB,
	repository.NewModelBulkExecutorMap,
)
