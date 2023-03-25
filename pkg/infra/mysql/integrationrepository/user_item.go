package integrationrepository

import (
	"context"
	"game-server-example/pkg/entity"
	"game-server-example/pkg/ifrepository"
	"game-server-example/pkg/infra/mysql/cacherepository"
	"game-server-example/pkg/infra/mysql/repository"
)

type ItemIntegrationRepository struct {
	itemRepository ifrepository.IFItemEntityRepository
}

func NewItemIntegrationRepository(entityRepository *repository.ItemEntityRepository, cacheRepository *cacherepository.ItemEntityCacheRepository, enableCache bool) *ItemIntegrationRepository {
	if enableCache {
		return &ItemIntegrationRepository{itemRepository: cacheRepository}
	}
	return &ItemIntegrationRepository{itemRepository: entityRepository}
}

func (r *ItemIntegrationRepository) Create(ctx context.Context, userId, itemId int64) (*entity.ItemEntity, error) {
	return r.itemRepository.Create(ctx, userId, itemId)
}

func (r *ItemIntegrationRepository) Save(ctx context.Context, e *entity.ItemEntity) error {
	return r.itemRepository.Save(ctx, e)
}

func (r *ItemIntegrationRepository) Get(ctx context.Context, userId, itemId int64) (*entity.ItemEntity, error) {
	return r.itemRepository.Get(ctx, userId, itemId)
}

func (r *ItemIntegrationRepository) FindByItemIds(ctx context.Context, userId int64, itemIds []int64) ([]*entity.ItemEntity, error) {
	return r.itemRepository.FindByItemIds(ctx, userId, itemIds)
}

func (r *ItemIntegrationRepository) FindByUserId(ctx context.Context, userId int64) ([]*entity.ItemEntity, error) {
	return r.itemRepository.FindByUserId(ctx, userId)
}
