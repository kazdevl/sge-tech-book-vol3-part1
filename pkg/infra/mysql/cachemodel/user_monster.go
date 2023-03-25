package cachemodel

import (
	"database/sql"
	"errors"
	"fmt"
	"game-server-example/pkg/infra/mysql/cachedb"
	"game-server-example/pkg/infra/mysql/datamodel"
	"github.com/samber/lo"
)

type UserMonsterCacheModel struct {
	datamodel.UserMonster
	cacheStatus cachedb.CacheStatus
}

func NewUserMonsterCacheModel(userMonster datamodel.UserMonster, cacheStatus cachedb.CacheStatus) *UserMonsterCacheModel {
	return &UserMonsterCacheModel{UserMonster: userMonster, cacheStatus: cacheStatus}
}

func (m *UserMonsterCacheModel) Bind(rows *sql.Rows) error {
	return rows.Scan(&m.UserID, &m.MonsterID, &m.Exp)
}

func (m *UserMonsterCacheModel) UniqueKeyKCondition() string {
	return "user_id=? AND monster_id=?"
}

func (m *UserMonsterCacheModel) UniqueKeyConditionValues() []any {
	return []any{m.UserID, m.MonsterID}
}

func (m *UserMonsterCacheModel) Table() string {
	return "user_monster"
}

func (m *UserMonsterCacheModel) UniqueKeyColumnValueStr() string {
	return fmt.Sprintf("user_id=%d,monster_id=%d", m.UserID, m.MonsterID)
}

func (m *UserMonsterCacheModel) SetCacheStatus(s cachedb.CacheStatus) {
	m.cacheStatus = s
}

func (m *UserMonsterCacheModel) GetCacheStatus() cachedb.CacheStatus {
	return m.cacheStatus
}

func (m *UserMonsterCacheModel) Update(content cachedb.CacheContent) error {
	model, ok := content.(*UserMonsterCacheModel)
	if !ok {
		return errors.New("想定していない型で更新しようとしています")
	}
	m.UserID = model.UserID
	m.MonsterID = model.MonsterID
	m.Exp = model.Exp
	return nil
}

func (m *UserMonsterCacheModel) CreateCopy() cachedb.CacheContent {
	newUserCoin := &UserMonsterCacheModel{
		UserMonster: datamodel.UserMonster{
			UserID:    m.UserID,
			MonsterID: m.MonsterID,
			Exp:       m.Exp,
		},
		cacheStatus: m.cacheStatus,
	}
	return newUserCoin
}

func (m *UserMonsterCacheModel) IsSame(column string, v any) bool {
	switch column {
	case "user_id":
		userId, ok := v.(int64)
		if !ok {
			return false
		}
		return m.UserID == userId
	case "monster_id":
		monsterId, ok := v.(int64)
		if !ok {
			return false
		}
		return m.MonsterID == monsterId
	case "Exp":
		count, ok := v.(int64)
		if !ok {
			return false
		}
		return m.Exp == count
	default:
		return false
	}
}

func (m *UserMonsterCacheModel) IsInclude(column string, vs []any) bool {
	switch column {
	case "user_id":
		userIds := m.toInt64Slice(vs)
		if len(userIds) == 0 {
			return false
		}
		return lo.Contains(userIds, m.UserID)
	case "monster_id":
		monsterIds := m.toInt64Slice(vs)
		if len(monsterIds) == 0 {
			return false
		}
		return lo.Contains(monsterIds, m.UserID)
	case "Exp":
		counts := m.toInt64Slice(vs)
		if len(counts) == 0 {
			return false
		}
		return lo.Contains(counts, m.Exp)
	default:
		return false
	}
}

func (m *UserMonsterCacheModel) toInt64Slice(vs []any) []int64 {
	result := make([]int64, 0, len(vs))
	for _, v := range vs {
		int64V, ok := v.(int64)
		if !ok {
			return nil
		}
		result = append(result, int64V)
	}
	return result
}

type UserMonsterCacheModels []*UserMonsterCacheModel

func (ms UserMonsterCacheModels) Table() string {
	return "user_monster"
}

func (ms *UserMonsterCacheModels) BindAndAddContent(rows *sql.Rows) error {
	m := &UserMonsterCacheModel{}
	if err := rows.Scan(&m.UserID, &m.MonsterID, &m.Exp); err != nil {
		return err
	}
	*ms = append(*ms, m)

	return nil
}

func (ms UserMonsterCacheModels) ForEach(f func(cacheContent cachedb.CacheContent)) {
	for _, m := range ms {
		f(m)
	}
}

func (ms *UserMonsterCacheModels) Add(content cachedb.CacheContent) {
	m, ok := content.(*UserMonsterCacheModel)
	if ok {
		*ms = append(*ms, m)
	}
}

func (ms *UserMonsterCacheModels) Reset() {
	*ms = make(UserMonsterCacheModels, 0)
}
