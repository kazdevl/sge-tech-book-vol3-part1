package ifrepository

import (
	"context"
	"game-server-example/pkg/entity"
)

type IFItemEntityRepository interface {
	Create(ctx context.Context, userId, itemId int64) (*entity.ItemEntity, error)
	Save(ctx context.Context, e *entity.ItemEntity) error
	Get(ctx context.Context, userId, itemId int64) (*entity.ItemEntity, error)
	FindByItemIds(ctx context.Context, userId int64, itemIds []int64) ([]*entity.ItemEntity, error)
	FindByUserId(ctx context.Context, userId int64) ([]*entity.ItemEntity, error)
}
