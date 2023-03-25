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

type CoinEntityRepository struct {
	db *mysql.ApplicationDB
}

func NewCoinEntityRepository(db *mysql.ApplicationDB) *CoinEntityRepository {
	return &CoinEntityRepository{db: db}
}

func (r *CoinEntityRepository) Create(ctx context.Context, userId int64) (*entity.CoinEntity, error) {
	m := &datamodel.UserCoin{UserID: userId}
	if err := m.Insert(ctx, r.db.Tx, boil.Infer()); err != nil {
		return nil, err
	}
	return entity.NewCoinEntity(m), nil
}

func (r *CoinEntityRepository) Save(ctx context.Context, e *entity.CoinEntity) error {
	m := &datamodel.UserCoin{UserID: e.UserId(), Num: e.Num()}
	_, err := m.Update(ctx, r.db.Tx, boil.Infer())
	return err
}

func (r *CoinEntityRepository) Get(ctx context.Context, userId int64) (*entity.CoinEntity, error) {
	m, err := datamodel.UserCoins(
		datamodel.UserCoinWhere.UserID.EQ(userId),
	).One(ctx, r.db.Db)
	if err != nil {
		return nil, err
	}
	return entity.NewCoinEntity(m), nil
}

func (r *CoinEntityRepository) BulkInsert(ctx context.Context, tx *sqlx.Tx, contents []cachedb.CacheContent) error {
	dms, err := r.toDataModels(contents)
	if err != nil {
		return err
	}
	return dms.InsertAll(ctx, tx, boil.Infer())
}

func (r *CoinEntityRepository) BulkUpdate(ctx context.Context, tx *sqlx.Tx, contents []cachedb.CacheContent) error {
	dms, err := r.toDataModels(contents)
	if err != nil {
		return err
	}

	queryWithoutArgs := "UPDATE user_coin SET " + r.createUpdateCaseQuery(dms) + " WHERE " + r.createUpdateWhereQuery(dms)
	if _, err = tx.ExecContext(ctx, queryWithoutArgs); err != nil {
		return err
	}

	return nil
}

func (r *CoinEntityRepository) BulkDelete(ctx context.Context, tx *sqlx.Tx, contents []cachedb.CacheContent) error {
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

func (r *CoinEntityRepository) toDataModels(contents []cachedb.CacheContent) (datamodel.UserCoinSlice, error) {
	var err error
	dms := lo.Map(contents, func(content cachedb.CacheContent, _ int) *datamodel.UserCoin {
		m, ok := content.(*cachemodel.UserCoinCacheModel)
		if !ok {
			err = errors.New("想定していない型です")
		}
		return &datamodel.UserCoin{
			UserID: m.UserID,
			Num:    m.Num,
		}
	})
	return dms, err
}

func (r *CoinEntityRepository) createUpdateCaseQuery(dms []*datamodel.UserCoin) string {
	query := "num = case "
	for _, dm := range dms {
		condition := fmt.Sprintf("WHEN user_id=%d THEN %d ", dm.UserID, dm.Num)
		query += condition
	}
	query += "Else 0 End"
	return query
}

func (r *CoinEntityRepository) createUpdateWhereQuery(dms []*datamodel.UserCoin) string {
	query := ""
	for index, dm := range dms {
		condition := fmt.Sprintf("(user_id=%d)", dm.UserID)
		if index != len(dms)-1 {
			condition += " OR "
		}
		query += condition
	}
	return query
}
