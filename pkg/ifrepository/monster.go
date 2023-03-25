package ifrepository

import (
	"context"
	"game-server-example/pkg/entity"
)

type IFMonsterEntityRepository interface {
	Create(ctx context.Context, userId, monsterId int64) (*entity.MonsterEntity, error)
	Get(ctx context.Context, userId, monsterId int64) (*entity.MonsterEntity, error)
	Save(ctx context.Context, e *entity.MonsterEntity) error
	FindByUserId(ctx context.Context, userId int64) ([]*entity.MonsterEntity, error)
}
