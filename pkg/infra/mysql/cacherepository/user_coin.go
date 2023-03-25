package cacherepository

import (
	"context"
	"game-server-example/pkg/entity"
	"game-server-example/pkg/infra/mysql/cachedb"
	"game-server-example/pkg/infra/mysql/cachemodel"
	"game-server-example/pkg/infra/mysql/datamodel"
)

type CoinEntityCacheRepository struct {
	cacheDB *cachedb.CacheDB
}

func NewCoinEntityCacheRepository(cacheDB *cachedb.CacheDB) *CoinEntityCacheRepository {
	return &CoinEntityCacheRepository{cacheDB: cacheDB}
}

func (r *CoinEntityCacheRepository) Create(ctx context.Context, userId int64) (*entity.CoinEntity, error) {
	m := &datamodel.UserCoin{UserID: userId}
	if err := r.cacheDB.Insert(ctx, cachemodel.NewUserCoinCacheModel(*m, cachedb.Insert)); err != nil {
		return nil, err
	}
	return entity.NewCoinEntity(m), nil
}

func (r *CoinEntityCacheRepository) Save(ctx context.Context, e *entity.CoinEntity) error {
	m := datamodel.UserCoin{UserID: e.UserId(), Num: e.Num()}
	if err := r.cacheDB.Update(ctx, cachemodel.NewUserCoinCacheModel(m, cachedb.Update)); err != nil {
		return err
	}
	return nil
}

func (r *CoinEntityCacheRepository) Get(ctx context.Context, userId int64) (*entity.CoinEntity, error) {
	m := datamodel.UserCoin{UserID: userId}
	cachedM := cachemodel.NewUserCoinCacheModel(m, cachedb.None)
	if err := r.cacheDB.GetAndSet(ctx, cachedM); err != nil {
		return nil, err
	}

	return entity.NewCoinEntity(&datamodel.UserCoin{
		UserID: cachedM.UserID,
		Num:    cachedM.Num,
	}), nil
}
