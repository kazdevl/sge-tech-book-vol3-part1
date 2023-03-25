package masterdata

import (
	"errors"
	"github.com/samber/lo"
)

type MonsterEnhanceTable struct {
	Level      int64 `csv:"level"`
	GrowthRate int64 `csv:"growth_rate"` // 千分率
	TotalExp   int64 `csv:"total_exp"`
	Coin       int64 `csv:"coin"` // このレベルになるのに必要なコイン
}

type MonsterEnhanceTables []*MonsterEnhanceTable

type MonsterEnhanceTableMap map[int64]*MonsterEnhanceTable

type MonsterEnhanceTableManager struct {
	monsterEnhanceTables   MonsterEnhanceTables
	monsterEnhanceTableMap MonsterEnhanceTableMap
}

func (ms MonsterEnhanceTables) ToMap() map[int64]*MonsterEnhanceTable {
	return lo.SliceToMap(ms, func(m *MonsterEnhanceTable) (int64, *MonsterEnhanceTable) {
		return m.Level, m
	})
}

func (mm MonsterEnhanceTableManager) GetSafely(level int64) (*MonsterEnhanceTable, error) {
	m, ok := mm.monsterEnhanceTableMap[level]
	if !ok {
		return nil, errors.New("存在しないlevelです")
	}
	return m, nil
}

func (mm MonsterEnhanceTableManager) Get(level int64) *MonsterEnhanceTable {
	return mm.monsterEnhanceTableMap[level]
}

func (mm MonsterEnhanceTableManager) All() MonsterEnhanceTables {
	return mm.monsterEnhanceTables
}
