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

type MonsterEntityRepository struct {
	db *mysql.ApplicationDB
}

func NewMonsterEntityRepository(db *mysql.ApplicationDB) *MonsterEntityRepository {
	return &MonsterEntityRepository{db: db}
}

func (r *MonsterEntityRepository) Create(ctx context.Context, userId, monsterId int64) (*entity.MonsterEntity, error) {
	m := &datamodel.UserMonster{UserID: userId, MonsterID: monsterId}
	if err := m.Insert(ctx, r.db.Tx, boil.Infer()); err != nil {
		return nil, err
	}
	return entity.NewMonsterEntity(m)
}

func (r *MonsterEntityRepository) Get(ctx context.Context, userId, monsterId int64) (*entity.MonsterEntity, error) {
	m, err := datamodel.UserMonsters(
		datamodel.UserMonsterWhere.UserID.EQ(userId),
		datamodel.UserMonsterWhere.MonsterID.EQ(monsterId),
	).One(ctx, r.db.Db)
	if err != nil {
		return nil, err
	}
	return entity.NewMonsterEntity(m)
}

func (r *MonsterEntityRepository) Save(ctx context.Context, e *entity.MonsterEntity) error {
	m := &datamodel.UserMonster{UserID: e.UserId(), MonsterID: e.MonsterId(), Exp: e.Exp()}
	_, err := m.Update(ctx, r.db.Tx, boil.Infer())
	return err
}

func (r *MonsterEntityRepository) FindByUserId(ctx context.Context, userId int64) ([]*entity.MonsterEntity, error) {
	ms, err := datamodel.UserMonsters(
		datamodel.UserMonsterWhere.UserID.EQ(userId),
	).All(ctx, r.db.Db)
	if err != nil {
		return nil, err
	}

	es := make([]*entity.MonsterEntity, 0, len(ms))
	for _, m := range ms {
		e, err := entity.NewMonsterEntity(m)
		if err != nil {
			return nil, err
		}
		es = append(es, e)
	}
	return es, nil
}

func (r *MonsterEntityRepository) BulkInsert(ctx context.Context, tx *sqlx.Tx, contents []cachedb.CacheContent) error {
	dms, err := r.toDataModels(contents)
	if err != nil {
		return err
	}
	return dms.InsertAll(ctx, tx, boil.Infer())
}

func (r *MonsterEntityRepository) BulkUpdate(ctx context.Context, tx *sqlx.Tx, contents []cachedb.CacheContent) error {
	dms, err := r.toDataModels(contents)
	if err != nil {
		return err
	}

	queryWithoutArgs := "UPDATE user_monster SET " + r.createUpdateCaseQuery(dms) + " WHERE " + r.createUpdateWhereQuery(dms)
	if _, err = tx.ExecContext(ctx, queryWithoutArgs); err != nil {
		return err
	}

	return nil
}

func (r *MonsterEntityRepository) BulkDelete(ctx context.Context, tx *sqlx.Tx, contents []cachedb.CacheContent) error {
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

func (r *MonsterEntityRepository) toDataModels(contents []cachedb.CacheContent) (datamodel.UserMonsterSlice, error) {
	var err error
	dms := lo.Map(contents, func(content cachedb.CacheContent, _ int) *datamodel.UserMonster {
		m, ok := content.(*cachemodel.UserMonsterCacheModel)
		if !ok {
			err = errors.New("想定していない型です")
		}
		return &datamodel.UserMonster{
			UserID:    m.UserID,
			MonsterID: m.MonsterID,
			Exp:       m.Exp,
		}
	})
	return dms, err
}

func (r *MonsterEntityRepository) createUpdateCaseQuery(dms []*datamodel.UserMonster) string {
	query := "exp = case "
	for _, dm := range dms {
		condition := fmt.Sprintf("WHEN user_id=%d AND monster_id=%d THEN %d ", dm.UserID, dm.MonsterID, dm.Exp)
		query += condition
	}
	query += "Else 0 End"
	return query
}

func (r *MonsterEntityRepository) createUpdateWhereQuery(dms []*datamodel.UserMonster) string {
	query := ""
	for index, dm := range dms {
		condition := fmt.Sprintf("(user_id=%d AND monster_id=%d)", dm.UserID, dm.MonsterID)
		if index != len(dms)-1 {
			condition += " OR "
		}
		query += condition
	}
	return query
}
