package cacherepository

import (
	"context"
	"errors"
	"game-server-example/pkg/entity"
	"game-server-example/pkg/infra/mysql/cachedb"
	"game-server-example/pkg/infra/mysql/cachemodel"
	"game-server-example/pkg/infra/mysql/datamodel"
)

type MonsterEntityCacheRepository struct {
	cacheDB *cachedb.CacheDB
}

func NewMonsterEntityCacheRepository(cacheDB *cachedb.CacheDB) *MonsterEntityCacheRepository {
	return &MonsterEntityCacheRepository{cacheDB: cacheDB}
}

func (r *MonsterEntityCacheRepository) Create(ctx context.Context, userId, monsterId int64) (*entity.MonsterEntity, error) {
	m := &datamodel.UserMonster{UserID: userId, MonsterID: monsterId}
	if err := r.cacheDB.Insert(ctx, cachemodel.NewUserMonsterCacheModel(*m, cachedb.Insert)); err != nil {
		return nil, err
	}
	return entity.NewMonsterEntity(m)
}

func (r *MonsterEntityCacheRepository) Save(ctx context.Context, e *entity.MonsterEntity) error {
	m := datamodel.UserMonster{UserID: e.UserId(), MonsterID: e.MonsterId(), Exp: e.Exp()}
	if err := r.cacheDB.Update(ctx, cachemodel.NewUserMonsterCacheModel(m, cachedb.Update)); err != nil {
		return err
	}
	return nil
}

func (r *MonsterEntityCacheRepository) Get(ctx context.Context, userId, monsterId int64) (*entity.MonsterEntity, error) {
	m := datamodel.UserMonster{UserID: userId, MonsterID: monsterId}
	cachedM := cachemodel.NewUserMonsterCacheModel(m, cachedb.None)
	if err := r.cacheDB.GetAndSet(ctx, cachedM); err != nil {
		return nil, err
	}

	return entity.NewMonsterEntity(&datamodel.UserMonster{
		UserID:    cachedM.UserID,
		MonsterID: cachedM.MonsterID,
		Exp:       cachedM.Exp,
	})
}

func (r *MonsterEntityCacheRepository) FindByUserId(ctx context.Context, userId int64) ([]*entity.MonsterEntity, error) {
	cachedM := make(cachemodel.UserMonsterCacheModels, 0)
	if err := r.cacheDB.FindAndSetByConditions(
		ctx,
		&cachedM,
		[]*cachedb.ConditionValue{
			{TargetColumn: "user_id", ConditionType: cachedb.Eq, Values: []any{userId}},
		},
	); err != nil {
		return nil, err
	}

	return r.toEntities(&cachedM)
}

func (r *MonsterEntityCacheRepository) FindByMonsterIds(ctx context.Context, userId int64, monsterIds []int64) ([]*entity.MonsterEntity, error) {
	cachedM := make(cachemodel.UserMonsterCacheModels, 0)
	monsterAnyIds := make([]any, 0, len(monsterIds))
	for _, monsterId := range monsterIds {
		monsterAnyIds = append(monsterAnyIds, monsterId)
	}

	if err := r.cacheDB.FindAndSetByConditions(
		ctx,
		&cachedM,
		[]*cachedb.ConditionValue{
			{TargetColumn: "user_id", ConditionType: cachedb.Eq, Values: []any{userId}},
			{TargetColumn: "monster_id", ConditionType: cachedb.In, Values: monsterAnyIds},
		},
	); err != nil {
		return nil, err
	}

	return r.toEntities(&cachedM)
}

func (r *MonsterEntityCacheRepository) toEntities(contents cachedb.SelectCacheContents) ([]*entity.MonsterEntity, error) {
	var err error
	es := make([]*entity.MonsterEntity, 0)
	contents.ForEach(func(content cachedb.CacheContent) {
		m, ok := content.(*cachemodel.UserMonsterCacheModel)
		if !ok {
			err = errors.New("想定していない型です")
		}
		dm := &datamodel.UserMonster{
			UserID:    m.UserID,
			MonsterID: m.MonsterID,
			Exp:       m.Exp,
		}
		var e *entity.MonsterEntity
		e, err = entity.NewMonsterEntity(dm)
		es = append(es, e)
	})
	return es, err
}
