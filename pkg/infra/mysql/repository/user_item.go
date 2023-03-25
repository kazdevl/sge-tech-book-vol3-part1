package repository

import (
	"context"
	"errors"
	"fmt"
	"game-server-example/pkg/entity"
	"game-server-example/pkg/infra/mysql"
	"game-server-example/pkg/infra/mysql/cachedb"
	"game-server-example/pkg/infra/mysql/cachemodel"
	"game-server-example/pkg/infra/mysql/datamodel"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type ItemEntityRepository struct {
	db *mysql.ApplicationDB
}

func NewItemEntityRepository(db *mysql.ApplicationDB) *ItemEntityRepository {
	return &ItemEntityRepository{db: db}
}

func (r *ItemEntityRepository) Create(ctx context.Context, userId, itemId int64) (*entity.ItemEntity, error) {
	m := &datamodel.UserItem{UserID: userId, ItemID: itemId}
	if err := m.Insert(ctx, r.db.Tx, boil.Infer()); err != nil {
		return nil, err
	}
	return entity.NewItemEntity(m)
}

func (r *ItemEntityRepository) Save(ctx context.Context, e *entity.ItemEntity) error {
	m := &datamodel.UserItem{UserID: e.UserId(), ItemID: e.ItemId(), Count: e.Count()}
	_, err := m.Update(ctx, r.db.Tx, boil.Infer())
	return err
}

func (r *ItemEntityRepository) Get(ctx context.Context, userId, itemId int64) (*entity.ItemEntity, error) {
	m, err := datamodel.UserItems(
		datamodel.UserItemWhere.UserID.EQ(userId),
		datamodel.UserItemWhere.ItemID.EQ(itemId),
	).One(ctx, r.db.Db)
	if err != nil {
		return nil, err
	}
	return entity.NewItemEntity(m)
}

func (r *ItemEntityRepository) FindByItemIds(ctx context.Context, userId int64, itemIds []int64) ([]*entity.ItemEntity, error) {
	ms, err := datamodel.UserItems(
		datamodel.UserItemWhere.UserID.EQ(userId),
		datamodel.UserItemWhere.ItemID.IN(itemIds),
	).All(ctx, r.db.Db)
	if err != nil {
		return nil, err
	}

	es := make([]*entity.ItemEntity, 0, len(ms))
	for _, m := range ms {
		e, err := entity.NewItemEntity(m)
		if err != nil {
			return nil, err
		}
		es = append(es, e)
	}
	return es, nil
}

func (r *ItemEntityRepository) FindByUserId(ctx context.Context, userId int64) ([]*entity.ItemEntity, error) {
	ms, err := datamodel.UserItems(
		datamodel.UserItemWhere.UserID.EQ(userId),
	).All(ctx, r.db.Db)
	if err != nil {
		return nil, err
	}

	es := make([]*entity.ItemEntity, 0, len(ms))
	for _, m := range ms {
		e, err := entity.NewItemEntity(m)
		if err != nil {
			return nil, err
		}
		es = append(es, e)
	}
	return es, nil
}

func (r *ItemEntityRepository) BulkInsert(ctx context.Context, tx *sqlx.Tx, contents []cachedb.CacheContent) error {
	dms, err := r.toDataModels(contents)
	if err != nil {
		return err
	}
	return dms.InsertAll(ctx, tx, boil.Infer())
}

func (r *ItemEntityRepository) BulkUpdate(ctx context.Context, tx *sqlx.Tx, contents []cachedb.CacheContent) error {
	dms, err := r.toDataModels(contents)
	if err != nil {
		return err
	}

	queryWithoutArgs := "UPDATE user_item SET " + r.createUpdateCaseQuery(dms) + " WHERE " + r.createUpdateWhereQuery(dms)
	if _, err = tx.ExecContext(ctx, queryWithoutArgs); err != nil {
		return err
	}

	return nil
}

func (r *ItemEntityRepository) BulkDelete(ctx context.Context, tx *sqlx.Tx, contents []cachedb.CacheContent) error {
	dms, err := r.toDataModels(contents)
	if err != nil {
		return err
	}
	effectedRows, err := dms.DeleteAll(ctx, tx)
	if err != nil {
		return err
	}
	if effectedRows < int64(len(contents)) {
		return errors.New("変更漏れしています")
	}
	return nil
}

func (r *ItemEntityRepository) toDataModels(contents []cachedb.CacheContent) (datamodel.UserItemSlice, error) {
	var err error
	dms := lo.Map(contents, func(content cachedb.CacheContent, _ int) *datamodel.UserItem {
		m, ok := content.(*cachemodel.UserItemCacheModel)
		if !ok {
			err = errors.New("想定していない型です")
		}
		return &datamodel.UserItem{
			UserID: m.UserID,
			ItemID: m.ItemID,
			Count:  m.Count,
		}
	})
	return dms, err
}

func (r *ItemEntityRepository) createUpdateCaseQuery(dms []*datamodel.UserItem) string {
	query := "count = case "
	for _, dm := range dms {
		condition := fmt.Sprintf("WHEN user_id=%d AND item_id=%d THEN %d ", dm.UserID, dm.ItemID, dm.Count)
		query += condition
	}
	query += "Else 0 End"
	return query
}

func (r *ItemEntityRepository) createUpdateWhereQuery(dms []*datamodel.UserItem) string {
	query := ""
	for index, dm := range dms {
		condition := fmt.Sprintf("(user_id=%d AND item_id=%d)", dm.UserID, dm.ItemID)
		if index != len(dms)-1 {
			condition += " OR "
		}
		query += condition
	}
	return query
}
