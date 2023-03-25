package entity

import (
	"game-server-example/pkg/infra/mysql/datamodel"
	"game-server-example/pkg/masterdata"
)

const maxLevel = 100

type MonsterEntity struct {
	userId  int64
	monster *masterdata.Monster
	exp     int64
}

func NewMonsterEntity(m *datamodel.UserMonster) (*MonsterEntity, error) {
	monster, err := masterdata.Master().Monster.GetSafely(m.MonsterID)
	if err != nil {
		return nil, err
	}

	return &MonsterEntity{
		userId:  m.UserID,
		monster: monster,
		exp:     m.Exp,
	}, nil
}

func (e *MonsterEntity) UserId() int64 {
	return e.userId
}

func (e *MonsterEntity) MonsterId() int64 {
	return e.monster.Id
}

func (e *MonsterEntity) Exp() int64 {
	return e.exp
}

func (e *MonsterEntity) CanAddExp() bool {
	return e.Level() != maxLevel
}

func (e *MonsterEntity) AddExp(exp int64) {
	if e.CanAddExp() {
		e.exp += exp
	}
}

func (e *MonsterEntity) Level() int64 {
	for _, enhanceInfo := range masterdata.Master().MonsterEnhanceTable.All() {
		if enhanceInfo.TotalExp >= e.exp {
			return enhanceInfo.Level
		}
	}
	return 0
}
