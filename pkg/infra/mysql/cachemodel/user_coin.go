package cachemodel

import (
	"database/sql"
	"errors"
	"fmt"
	"game-server-example/pkg/infra/mysql/cachedb"
	"game-server-example/pkg/infra/mysql/datamodel"
	"github.com/samber/lo"
)

type UserCoinCacheModel struct {
	datamodel.UserCoin
	cacheStatus cachedb.CacheStatus
}

func NewUserCoinCacheModel(userCoin datamodel.UserCoin, cacheStatus cachedb.CacheStatus) *UserCoinCacheModel {
	return &UserCoinCacheModel{UserCoin: userCoin, cacheStatus: cacheStatus}
}

func (m *UserCoinCacheModel) Bind(rows *sql.Rows) error {
	return rows.Scan(&m.UserID, &m.Num)
}

func (m *UserCoinCacheModel) UniqueKeyKCondition() string {
	return "user_id=?"
}

func (m *UserCoinCacheModel) UniqueKeyConditionValues() []any {
	return []any{m.UserID}
}

func (m *UserCoinCacheModel) Table() string {
	return "user_coin"
}

func (m *UserCoinCacheModel) UniqueKeyColumnValueStr() string {
	return fmt.Sprintf("user_id=%d", m.UserID)
}

func (m *UserCoinCacheModel) SetCacheStatus(s cachedb.CacheStatus) {
	m.cacheStatus = s
}

func (m *UserCoinCacheModel) GetCacheStatus() cachedb.CacheStatus {
	return m.cacheStatus
}

func (m *UserCoinCacheModel) Update(content cachedb.CacheContent) error {
	model, ok := content.(*UserCoinCacheModel)
	if !ok {
		return errors.New("想定していない型で更新しようとしています")
	}
	m.UserID = model.UserID
	m.Num = model.Num
	return nil
}

func (m *UserCoinCacheModel) CreateCopy() cachedb.CacheContent {
	newUserCoin := &UserCoinCacheModel{
		UserCoin: datamodel.UserCoin{
			UserID: m.UserID,
			Num:    m.Num,
		},
		cacheStatus: m.cacheStatus,
	}
	return newUserCoin
}

func (m *UserCoinCacheModel) IsSame(column string, v any) bool {
	switch column {
	case "user_id":
		userId, ok := v.(int64)
		if !ok {
			return false
		}
		return m.UserID == userId
	case "num":
		num, ok := v.(int64)
		if !ok {
			return false
		}
		return m.Num == num
	default:
		return false
	}
}

func (m *UserCoinCacheModel) IsInclude(column string, vs []any) bool {
	switch column {
	case "user_id":
		userIDs := m.toInt64Slice(vs)
		if len(userIDs) == 0 {
			return false
		}
		return lo.Contains(userIDs, m.UserID)
	case "Num":
		nums := m.toInt64Slice(vs)
		if len(nums) == 0 {
			return false
		}
		return lo.Contains(nums, m.Num)
	default:
		return false
	}
}

func (m *UserCoinCacheModel) toInt64Slice(vs []any) []int64 {
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

type UserCoinCacheModels []*UserCoinCacheModel

func (ms UserCoinCacheModels) Table() string {
	return "user_coin"
}

func (ms *UserCoinCacheModels) BindAndAddContent(rows *sql.Rows) error {
	m := &UserCoinCacheModel{}
	if err := rows.Scan(&m.UserID, &m.Num); err != nil {
		return err
	}
	*ms = append(*ms, m)

	return nil
}

func (ms UserCoinCacheModels) ForEach(f func(cacheContent cachedb.CacheContent)) {
	for _, m := range ms {
		f(m)
	}
}

func (ms *UserCoinCacheModels) Add(content cachedb.CacheContent) {
	m, ok := content.(*UserCoinCacheModel)
	if ok {
		*ms = append(*ms, m)
	}
}

func (ms *UserCoinCacheModels) Reset() {
	*ms = make(UserCoinCacheModels, 0)
}
