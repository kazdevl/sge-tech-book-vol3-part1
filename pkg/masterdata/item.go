package masterdata

import (
	"errors"
	"github.com/samber/lo"
)

type Item struct {
	Id       int64      `csv:"id"`
	Name     string     `csv:"name"`
	Rarity   RarityType `csv:"rarity"`
	ItemType ItemType   `csv:"item_type"`
	AddExp   int64      `csv:"add_exp"`
}

type Items []*Item

type ItemMap map[int64]*Item

type ItemManager struct {
	items   Items
	itemMap ItemMap
}

func (ms Items) ToMap() map[int64]*Item {
	return lo.SliceToMap(ms, func(m *Item) (int64, *Item) {
		return m.Id, m
	})
}

func (mm ItemManager) GetSafely(id int64) (*Item, error) {
	m, ok := mm.itemMap[id]
	if !ok {
		return nil, errors.New("存在しないidです")
	}
	return m, nil
}

func (mm ItemManager) Get(id int64) *Item {
	return mm.itemMap[id]
}

func (mm ItemManager) All() []*Item {
	return mm.items
}
