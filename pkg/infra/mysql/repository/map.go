package repository

import "game-server-example/pkg/infra/mysql/cachedb"

func NewModelBulkExecutorMap(
	coinEntityRepository *CoinEntityRepository,
	itemEntityRepository *ItemEntityRepository,
	monsterEntityRepository *MonsterEntityRepository,
) map[string]cachedb.BulkExecutor {
	return map[string]cachedb.BulkExecutor{
		"user_coin":    coinEntityRepository,
		"user_item":    itemEntityRepository,
		"user_monster": monsterEntityRepository,
	}
}
