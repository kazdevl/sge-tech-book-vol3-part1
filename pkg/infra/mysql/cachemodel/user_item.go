package cachemodel

import (
	"database/sql"
	"errors"
	"fmt"
	"game-server-example/pkg/infra/mysql/cachedb"
	"game-server-example/pkg/infra/mysql/datamodel"
	"github.com/samber/lo"
)

type UserItemCacheModel struct {
	datamodel.UserItem
	cacheStatus cachedb.CacheStatus
}

func NewUserItemCacheModel(userItem datamodel.UserItem, cacheStatus cachedb.CacheStatus) *UserItemCacheModel {
	return &UserItemCacheModel{UserItem: userItem, cacheStatus: cacheStatus}
}

func (m *UserItemCacheModel) Bind(rows *sql.Rows) error {
	return rows.Scan(&m.UserID, &m.ItemID, &m.Count)
}

func (m *UserItemCacheModel) UniqueKeyKCondition() string {
	return "user_id=? AND item_id=?"
}

func (m *UserItemCacheModel) UniqueKeyConditionValues() []any {
	return []any{m.UserID, m.ItemID}
}

func (m *UserItemCacheModel) Table() string {
	return "user_item"
}

func (m *UserItemCacheModel) UniqueKeyColumnValueStr() string {
	return fmt.Sprintf("user_id=%d,item_id=%d", m.UserID, m.ItemID)
}

func (m *UserItemCacheModel) SetCacheStatus(s cachedb.CacheStatus) {
	m.cacheStatus = s
}

func (m *UserItemCacheModel) GetCacheStatus() cachedb.CacheStatus {
	return m.cacheStatus
}

func (m *UserItemCacheModel) Update(content cachedb.CacheContent) error {
	model, ok := content.(*UserItemCacheModel)
	if !ok {
		return errors.New("想定していない型で更新しようとしています")
	}
	m.UserID = model.UserID
	m.ItemID = model.ItemID
	m.Count = model.Count
	return nil
}

func (m *UserItemCacheModel) CreateCopy() cachedb.CacheContent {
	newUserCoin := &UserItemCacheModel{
		UserItem: datamodel.UserItem{
			UserID: m.UserID,
			ItemID: m.ItemID,
			Count:  m.Count,
		},
		cacheStatus: m.cacheStatus,
	}
	return newUserCoin
}

func (m *UserItemCacheModel) IsSame(column string, v any) bool {
	switch column {
	case "user_id":
		userId, ok := v.(int64)
		if !ok {
			return false
		}
		return m.UserID == userId
	case "item_id":
		itemId, ok := v.(int64)
		if !ok {
			return false
		}
		return m.ItemID == itemId
	case "Count":
		count, ok := v.(int64)
		if !ok {
			return false
		}
		return m.Count == count
	default:
		return false
	}
}

func (m *UserItemCacheModel) IsInclude(column string, vs []any) bool {
	switch column {
	case "user_id":
		userIds := m.toInt64Slice(vs)
		if len(userIds) == 0 {
			return false
		}
		return lo.Contains(userIds, m.UserID)
	case "item_id":
		itemIds := m.toInt64Slice(vs)
		if len(itemIds) == 0 {
			return false
		}
		return lo.Contains(itemIds, m.UserID)
	case "Count":
		counts := m.toInt64Slice(vs)
		if len(counts) == 0 {
			return false
		}
		return lo.Contains(counts, m.Count)
	default:
		return false
	}
}

func (m *UserItemCacheModel) toInt64Slice(vs []any) []int64 {
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

type UserItemCacheModels []*UserItemCacheModel

func (ms UserItemCacheModels) Table() string {
	return "user_item"
}

func (ms *UserItemCacheModels) BindAndAddContent(rows *sql.Rows) error {
	m := &UserItemCacheModel{}
	if err := rows.Scan(&m.UserID, &m.ItemID, &m.Count); err != nil {
		return err
	}
	*ms = append(*ms, m)

	return nil
}

func (ms UserItemCacheModels) ForEach(f func(cacheContent cachedb.CacheContent)) {
	for _, m := range ms {
		f(m)
	}
}

func (ms *UserItemCacheModels) Add(content cachedb.CacheContent) {
	m, ok := content.(*UserItemCacheModel)
	if ok {
		*ms = append(*ms, m)
	}
}

func (ms *UserItemCacheModels) Reset() {
	*ms = make(UserItemCacheModels, 0)
}
