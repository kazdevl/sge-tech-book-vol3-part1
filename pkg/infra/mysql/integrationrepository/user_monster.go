package integrationrepository

import (
	"context"
	"game-server-example/pkg/entity"
	"game-server-example/pkg/ifrepository"
	"game-server-example/pkg/infra/mysql/cacherepository"
	"game-server-example/pkg/infra/mysql/repository"
)

type MonsterIntegrationRepository struct {
	monsterRepository ifrepository.IFMonsterEntityRepository
}

func NewMonsterIntegrationRepository(entityRepository *repository.MonsterEntityRepository, cacheRepository *cacherepository.MonsterEntityCacheRepository, enableCache bool) *MonsterIntegrationRepository {
	if enableCache {
		return &MonsterIntegrationRepository{monsterRepository: cacheRepository}
	}
	return &MonsterIntegrationRepository{monsterRepository: entityRepository}
}

func (r *MonsterIntegrationRepository) Create(ctx context.Context, userId, monsterId int64) (*entity.MonsterEntity, error) {
	return r.monsterRepository.Create(ctx, userId, monsterId)
}

func (r *MonsterIntegrationRepository) Save(ctx context.Context, e *entity.MonsterEntity) error {
	return r.monsterRepository.Save(ctx, e)
}

func (r *MonsterIntegrationRepository) Get(ctx context.Context, userId, monsterId int64) (*entity.MonsterEntity, error) {
	return r.monsterRepository.Get(ctx, userId, monsterId)
}

func (r *MonsterIntegrationRepository) FindByUserId(ctx context.Context, userId int64) ([]*entity.MonsterEntity, error) {
	return r.monsterRepository.FindByUserId(ctx, userId)
}
