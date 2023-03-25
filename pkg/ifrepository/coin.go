package ifrepository

import (
	"context"
	"game-server-example/pkg/entity"
)

type IFCoinEntityRepository interface {
	Create(ctx context.Context, userId int64) (*entity.CoinEntity, error)
	Save(ctx context.Context, e *entity.CoinEntity) error
	Get(ctx context.Context, userId int64) (*entity.CoinEntity, error)
}
