package integrationrepository

import (
	"context"
	"game-server-example/pkg/entity"
	"game-server-example/pkg/ifrepository"
	"game-server-example/pkg/infra/mysql/cacherepository"
	"game-server-example/pkg/infra/mysql/repository"
)

type CoinIntegrationRepository struct {
	coinRepository ifrepository.IFCoinEntityRepository
}

func NewCoinIntegrationRepository(entityRepository *repository.CoinEntityRepository, cacheRepository *cacherepository.CoinEntityCacheRepository, enableCache bool) *CoinIntegrationRepository {
	if enableCache {
		return &CoinIntegrationRepository{coinRepository: cacheRepository}
	}
	return &CoinIntegrationRepository{coinRepository: entityRepository}
}

func (r *CoinIntegrationRepository) Create(ctx context.Context, userId int64) (*entity.CoinEntity, error) {
	return r.coinRepository.Create(ctx, userId)
}

func (r *CoinIntegrationRepository) Save(ctx context.Context, e *entity.CoinEntity) error {
	return r.coinRepository.Save(ctx, e)
}

func (r *CoinIntegrationRepository) Get(ctx context.Context, userId int64) (*entity.CoinEntity, error) {
	return r.coinRepository.Get(ctx, userId)
}
