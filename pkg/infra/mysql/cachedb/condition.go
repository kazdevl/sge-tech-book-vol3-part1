package cachedb

import (
	"github.com/samber/lo"
	"sort"
	"strings"
)

type Condition int64

const (
	Eq Condition = iota
	In
)

type ConditionValue struct {
	TargetColumn  string
	ConditionType Condition
	Values        []any
	order         int
}

func (c *ConditionValue) setOrder(i int) {
	c.order = i
}

func (c *ConditionValue) isSame(condition *ConditionValue) bool {
	if c.TargetColumn != condition.TargetColumn {
		return false
	}
	if c.ConditionType != condition.ConditionType {
		return false
	}

	switch c.ConditionType {
	case Eq:
		return c.Values[0] == condition.Values[0]
	case In:
		if len(c.Values) < len(condition.Values) {
			return false
		}

		count := 0
		for _, value := range condition.Values {
			for _, cachedValue := range c.Values {
				if value == cachedValue {
					count++
					break
				}
			}
		}
		return count == len(condition.Values)
	}
	return false
}

func (c *ConditionValue) union(vs []any) {
	// TODO 型の保証 仮でint64
	int64Vs := make([]int64, 0, len(vs))
	for _, v := range vs {
		int64Vs = append(int64Vs, v.(int64))
	}
	int64Values := make([]int64, 0, len(c.Values))
	for _, value := range c.Values {
		int64Values = append(int64Values, value.(int64))
	}

	unionValues := lo.Union(int64Values, int64Vs)
	unionAnyValues := make([]any, 0, len(unionValues))
	for _, uv := range unionValues {
		unionAnyValues = append(unionAnyValues, uv)
	}
	c.Values = unionAnyValues
}

func (c *ConditionValue) getQueryInfo() (string, string, []any) {
	switch c.ConditionType {
	case Eq:
		return c.TargetColumn + " = ?", c.TargetColumn + " = ?", c.Values
	case In:
		return c.TargetColumn + " IN (?)", c.TargetColumn + " IN (" + c.getPlaceHolderStr() + ")", c.Values
	}
	return "", "", nil
}

func (c *ConditionValue) getPlaceHolderStr() string {
	strs := make([]string, 0, len(c.Values))

	for i := 0; i < len(c.Values); i++ {
		strs = append(strs, "?")
	}
	valuesStr := strings.Join(strs, ",")
	return valuesStr
}

type ConditionMap map[int]*ConditionValue

func toConditionMap(vs []*ConditionValue) ConditionMap {
	return lo.SliceToMap(vs, func(v *ConditionValue) (int, *ConditionValue) {
		return v.order, v
	})
}

func (cm ConditionMap) CreateQuery(table string) (string, string, []any) {
	cs := lo.MapToSlice(cm, func(_ int, value *ConditionValue) *ConditionValue {
		return value
	})

	sort.Slice(cs, func(i, j int) bool {
		return cs[i].order < cs[j].order
	})

	key := "SELECT * FROM " + table + " WHERE "
	query := "SELECT * FROM " + table + " WHERE "
	allArgs := make([]any, 0)
	for i, c := range cs {
		queryKey, conditionFormat, args := c.getQueryInfo()
		if len(cs)-1 != i {
			key += queryKey + " AND "
			query += conditionFormat + " AND "
		} else {
			key += queryKey
			query += conditionFormat
		}

		allArgs = append(allArgs, args...)
	}
	return key, query, allArgs
}

func (cm ConditionMap) needUnionConditionValue(addableConditionMap ConditionMap) ([]int, bool) {
	targetOrders := make([]int, 0)

	for order, c := range cm {
		condition := addableConditionMap[order]
		switch condition.ConditionType {
		case Eq:
			if !condition.isSame(c) {
				return nil, false
			}
		case In:
			if !condition.isSame(c) {
				targetOrders = append(targetOrders, order)
			}
		}
	}
	return targetOrders, len(targetOrders) > 0
}

type ConditionMaps []ConditionMap

func (cms ConditionMaps) IsOnceCalled(conditions []*ConditionValue) bool {
	for _, cm := range cms {
		if len(cm) != len(conditions) {
			continue
		}

		onceCalled := true
		for _, condition := range conditions {
			if cm[condition.order] != nil && !cm[condition.order].isSame(condition) {
				onceCalled = false
				break
			}
		}
		if onceCalled {
			return true
		}
	}
	return false
}

func (cms ConditionMaps) UnionConditionValues(addableConditions []*ConditionValue) bool {
	addableConditionMap := toConditionMap(addableConditions)

	for _, cm := range cms {
		targetOrders, needUnion := cm.needUnionConditionValue(addableConditionMap)
		if needUnion {
			for _, targetOrder := range targetOrders {
				cm[targetOrder].union(addableConditions[targetOrder].Values)
			}
			return true
		}
	}
	return false
}
