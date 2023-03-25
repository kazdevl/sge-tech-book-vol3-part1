package masterdata

import (
	"errors"
	"github.com/samber/lo"
)

type Monster struct {
	Id           int64       `csv:"id"`
	Name         string      `csv:"name"`
	Element      ElementType `csv:"element"`
	Rarity       RarityType  `csv:"rarity"`
	MaxHp        int64       `csv:"max_hp"`
	MaxMp        int64       `csv:"max_mp"`
	MaxOffensive int64       `csv:"max_offensive"`
	MaxDefensive int64       `csv:"max_defensive"`
	MaxSpeed     int64       `csv:"max_speed"`
}

type Monsters []*Monster

type MonsterMap map[int64]*Monster

type MonsterManager struct {
	monsters   Monsters
	monsterMap MonsterMap
}

func (ms Monsters) ToMap() map[int64]*Monster {
	return lo.SliceToMap(ms, func(m *Monster) (int64, *Monster) {
		return m.Id, m
	})
}

func (mm MonsterManager) GetSafely(id int64) (*Monster, error) {
	m, ok := mm.monsterMap[id]
	if !ok {
		return nil, errors.New("存在しないidです")
	}
	return m, nil
}

func (mm MonsterManager) Get(id int64) *Monster {
	return mm.monsterMap[id]
}

func (mm MonsterManager) All() []*Monster {
	return mm.monsters
}
