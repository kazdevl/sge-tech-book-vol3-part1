package entity

import (
	"errors"
	"game-server-example/pkg/infra/mysql/datamodel"
	"game-server-example/pkg/masterdata"
)

type ItemEntity struct {
	userId int64
	item   *masterdata.Item
	count  int64
}

func NewItemEntity(m *datamodel.UserItem) (*ItemEntity, error) {
	item, err := masterdata.Master().Item.GetSafely(m.ItemID)
	if err != nil {
		return nil, err
	}
	return &ItemEntity{userId: m.UserID, item: item, count: m.Count}, nil
}

func (e *ItemEntity) ItemId() int64 {
	return e.item.Id
}

func (e *ItemEntity) UserId() int64 {
	return e.userId
}

func (e *ItemEntity) Count() int64 {
	return e.count
}

func (e *ItemEntity) Consume(count int64) error {
	if e.count < count {
		return errors.New("消費数が所持数を上まっています")
	}
	e.count -= count
	return nil
}

func (e *ItemEntity) Exp() int64 {
	return e.item.AddExp
}

func (e *ItemEntity) Add(count int64) {
	e.count = count
}
