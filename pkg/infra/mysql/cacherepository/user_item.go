package cacherepository

import (
	"context"
	"errors"
	"game-server-example/pkg/entity"
	"game-server-example/pkg/infra/mysql/cachedb"
	"game-server-example/pkg/infra/mysql/cachemodel"
	"game-server-example/pkg/infra/mysql/datamodel"
)

type ItemEntityCacheRepository struct {
	cacheDB *cachedb.CacheDB
}

func NewItemEntityCacheRepository(cacheDB *cachedb.CacheDB) *ItemEntityCacheRepository {
	return &ItemEntityCacheRepository{cacheDB: cacheDB}
}

func (r *ItemEntityCacheRepository) Create(ctx context.Context, userId, itemId int64) (*entity.ItemEntity, error) {
	m := &datamodel.UserItem{UserID: userId, ItemID: itemId}
	if err := r.cacheDB.Insert(ctx, cachemodel.NewUserItemCacheModel(*m, cachedb.Insert)); err != nil {
		return nil, err
	}
	return entity.NewItemEntity(m)
}

func (r *ItemEntityCacheRepository) Save(ctx context.Context, e *entity.ItemEntity) error {
	m := datamodel.UserItem{UserID: e.UserId(), ItemID: e.ItemId(), Count: e.Count()}
	if err := r.cacheDB.Update(ctx, cachemodel.NewUserItemCacheModel(m, cachedb.Update)); err != nil {
		return err
	}
	return nil
}

func (r *ItemEntityCacheRepository) Get(ctx context.Context, userId, itemId int64) (*entity.ItemEntity, error) {
	m := datamodel.UserItem{UserID: userId, ItemID: itemId}
	cachedM := cachemodel.NewUserItemCacheModel(m, cachedb.None)
	if err := r.cacheDB.GetAndSet(ctx, cachedM); err != nil {
		return nil, err
	}

	return entity.NewItemEntity(&datamodel.UserItem{
		UserID: cachedM.UserID,
		ItemID: cachedM.ItemID,
		Count:  cachedM.Count,
	})
}

func (r *ItemEntityCacheRepository) FindByUserId(ctx context.Context, userId int64) ([]*entity.ItemEntity, error) {
	cachedM := make(cachemodel.UserItemCacheModels, 0)
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

func (r *ItemEntityCacheRepository) FindByItemIds(ctx context.Context, userId int64, itemIds []int64) ([]*entity.ItemEntity, error) {
	cachedM := make(cachemodel.UserItemCacheModels, 0)
	itemAnyIds := make([]any, 0, len(itemIds))
	for _, itemId := range itemIds {
		itemAnyIds = append(itemAnyIds, itemId)
	}

	if err := r.cacheDB.FindAndSetByConditions(
		ctx,
		&cachedM,
		[]*cachedb.ConditionValue{
			{TargetColumn: "user_id", ConditionType: cachedb.Eq, Values: []any{userId}},
			{TargetColumn: "item_id", ConditionType: cachedb.In, Values: itemAnyIds},
		},
	); err != nil {
		return nil, err
	}

	return r.toEntities(&cachedM)
}

func (r *ItemEntityCacheRepository) toEntities(contents cachedb.SelectCacheContents) ([]*entity.ItemEntity, error) {
	var err error
	es := make([]*entity.ItemEntity, 0)
	contents.ForEach(func(content cachedb.CacheContent) {
		m, ok := content.(*cachemodel.UserItemCacheModel)
		if !ok {
			err = errors.New("想定していない型です")
		}
		dm := &datamodel.UserItem{
			UserID: m.UserID,
			ItemID: m.ItemID,
			Count:  m.Count,
		}
		var e *entity.ItemEntity
		e, err = entity.NewItemEntity(dm)
		es = append(es, e)
	})
	return es, err
}
