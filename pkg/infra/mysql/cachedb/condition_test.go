package cachedb

import (
	"github.com/magiconair/properties/assert"
	"reflect"
	"sort"
	"testing"
)

func TestConditionValue_isSame(t *testing.T) {
	type fields struct {
		TargetColumn  string
		ConditionType Condition
		Values        []any
	}
	type args struct {
		condition *ConditionValue
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "正常：conditionがEqでカラムと値が等しい",
			fields: fields{
				TargetColumn:  "Sample",
				ConditionType: Eq,
				Values:        []any{1},
			},
			args: args{
				condition: &ConditionValue{
					TargetColumn:  "Sample",
					ConditionType: Eq,
					Values:        []any{1},
				},
			},
			want: true,
		},
		{
			name: "正常：conditionがInでカラムと値が等しい",
			fields: fields{
				TargetColumn:  "Sample",
				ConditionType: In,
				Values:        []any{1, 2, 3},
			},
			args: args{
				condition: &ConditionValue{
					TargetColumn:  "Sample",
					ConditionType: In,
					Values:        []any{1, 2, 3},
				},
			},
			want: true,
		},
		{
			name: "正常：conditionがInでカラムが同じで引数のconditionの値が全てcacheのconditionの値に含まれている",
			fields: fields{
				TargetColumn:  "Sample",
				ConditionType: In,
				Values:        []any{1, 2, 3},
			},
			args: args{
				condition: &ConditionValue{
					TargetColumn:  "Sample",
					ConditionType: In,
					Values:        []any{1, 2},
				},
			},
			want: true,
		},
		{
			name: "正常：カラムが異なる",
			fields: fields{
				TargetColumn:  "Sample1",
				ConditionType: In,
				Values:        []any{1, 2, 3},
			},
			args: args{
				condition: &ConditionValue{
					TargetColumn:  "Sample2",
					ConditionType: In,
					Values:        []any{1, 2, 3},
				},
			},
			want: false,
		},
		{
			name: "正常：conditionが異なる",
			fields: fields{
				TargetColumn:  "Sample",
				ConditionType: In,
				Values:        []any{1, 2, 3},
			},
			args: args{
				condition: &ConditionValue{
					TargetColumn:  "Sample",
					ConditionType: Eq,
					Values:        []any{1},
				},
			},
			want: false,
		},
		{
			name: "正常：conditionがEqで値が異なる",
			fields: fields{
				TargetColumn:  "Sample",
				ConditionType: Eq,
				Values:        []any{2},
			},
			args: args{
				condition: &ConditionValue{
					TargetColumn:  "Sample",
					ConditionType: Eq,
					Values:        []any{1},
				},
			},
			want: false,
		},
		{
			name: "正常：conditionがInで値が異なる",
			fields: fields{
				TargetColumn:  "Sample",
				ConditionType: In,
				Values:        []any{1, 2, 3},
			},
			args: args{
				condition: &ConditionValue{
					TargetColumn:  "Sample",
					ConditionType: In,
					Values:        []any{1, 2, 5},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConditionValue{
				TargetColumn:  tt.fields.TargetColumn,
				ConditionType: tt.fields.ConditionType,
				Values:        tt.fields.Values,
			}
			if got := c.isSame(tt.args.condition); got != tt.want {
				t.Errorf("isSame() = %v, wantKey %v", got, tt.want)
			}
		})
	}
}

func TestConditionValue_union(t *testing.T) {
	type fields struct {
		TargetColumn  string
		ConditionType Condition
		Values        []any
	}
	type args struct {
		vs []any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []any
	}{
		{
			name: "正常",
			fields: fields{
				Values: []any{int64(1), int64(2), int64(3)},
			},
			args: args{vs: []any{int64(1), int64(2), int64(4), int64(5)}},
			want: []any{int64(1), int64(2), int64(3), int64(4), int64(5)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConditionValue{
				TargetColumn:  tt.fields.TargetColumn,
				ConditionType: tt.fields.ConditionType,
				Values:        tt.fields.Values,
			}
			c.union(tt.args.vs)
			if !reflect.DeepEqual(c.Values, tt.want) {
				t.Errorf("union() = %v, wantKey %v", c.Values, tt.want)
			}
		})
	}
}

func TestConditionValue_getQueryInfo(t *testing.T) {
	type fields struct {
		TargetColumn  string
		ConditionType Condition
		Values        []any
		order         int
	}
	tests := []struct {
		name          string
		fields        fields
		wantKey       string
		wantCondition string
		wantArgs      []any
	}{
		{
			name: "正常：Eqのconditionのformatとargsの作成",
			fields: fields{
				TargetColumn:  "sample",
				ConditionType: Eq,
				Values:        []any{1},
			},
			wantKey:       "sample = ?",
			wantCondition: "sample = ?",
			wantArgs:      []any{1},
		},
		{
			name: "正常：Inのconditionのformatとargsの作成",
			fields: fields{
				TargetColumn:  "sample",
				ConditionType: In,
				Values:        []any{1, 2, 3},
			},
			wantKey:       "sample IN (?)",
			wantCondition: "sample IN (?,?,?)",
			wantArgs:      []any{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConditionValue{
				TargetColumn:  tt.fields.TargetColumn,
				ConditionType: tt.fields.ConditionType,
				Values:        tt.fields.Values,
				order:         tt.fields.order,
			}
			got, got1, got2 := c.getQueryInfo()
			if got != tt.wantKey {
				t.Errorf("getQueryInfo() got = %v, wantKey %v", got, tt.wantKey)
			}
			if got1 != tt.wantCondition {
				t.Errorf("getQueryInfo() got = %v, wantCondition %v", got, tt.wantCondition)
			}
			if !reflect.DeepEqual(got2, tt.wantArgs) {
				t.Errorf("getQueryInfo() got2 = %v, wantArgs %v", got1, tt.wantArgs)
			}
		})
	}
}

func TestConditionMap_CreateWhereConditionQuery(t *testing.T) {
	type args struct {
		table string
	}
	tests := []struct {
		name      string
		cm        ConditionMap
		args      args
		wantKey   string
		wantQuery string
		wantArgs  []any
	}{
		{
			name: "正常にクエリを作成できている",
			cm: toConditionMap([]*ConditionValue{
				{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
				{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{1, 2, 3}},
			}),
			args:      args{table: "user_item"},
			wantKey:   "SELECT * FROM user_item WHERE user_id = ? AND item_id IN (?)",
			wantQuery: "SELECT * FROM user_item WHERE user_id = ? AND item_id IN (?,?,?)",
			wantArgs:  []any{1, 1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := tt.cm.CreateQuery(tt.args.table)
			if got != tt.wantKey {
				t.Errorf("CreateQuery() got = %v, wantKey %v", got, tt.wantKey)
			}
			if got1 != tt.wantQuery {
				t.Errorf("CreateQuery() got1 = %v, wantQuery %v", got, tt.wantQuery)
			}
			if !reflect.DeepEqual(got2, tt.wantArgs) {
				t.Errorf("CreateQuery() got2 = %v, wantArgs %v", got1, tt.wantArgs)
			}
		})
	}
}

func TestConditionMap_needUpdateConditionValue(t *testing.T) {
	type args struct {
		addableConditionMap ConditionMap
	}
	tests := []struct {
		name             string
		cm               ConditionMap
		args             args
		wantTargetOrders []int
		wantNeedUnion    bool
	}{
		{
			name: "正常：更新不要",
			cm: toConditionMap([]*ConditionValue{
				{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
				{order: 1, TargetColumn: "item_id", ConditionType: Eq, Values: []any{2}},
			}),
			args: args{addableConditionMap: toConditionMap([]*ConditionValue{
				{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
				{order: 1, TargetColumn: "item_id", ConditionType: Eq, Values: []any{3}},
			})},
			wantTargetOrders: nil,
			wantNeedUnion:    false,
		},
		{
			name: "正常：IN句で更新必要_1",
			cm: toConditionMap([]*ConditionValue{
				{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
				{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{1, 2, 3}},
			}),
			args: args{addableConditionMap: toConditionMap([]*ConditionValue{
				{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
				{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{2, 3, 4}},
			})},
			wantTargetOrders: []int{1},
			wantNeedUnion:    true,
		},
		{
			name: "正常：IN句で更新必要_2",
			cm: toConditionMap([]*ConditionValue{
				{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
				{order: 1, TargetColumn: "count", ConditionType: Eq, Values: []any{10}},
				{order: 2, TargetColumn: "item_id", ConditionType: In, Values: []any{1, 2, 3}},
			}),
			args: args{addableConditionMap: toConditionMap([]*ConditionValue{
				{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
				{order: 1, TargetColumn: "count", ConditionType: Eq, Values: []any{10}},
				{order: 2, TargetColumn: "item_id", ConditionType: In, Values: []any{2, 3, 4}},
			})},
			wantTargetOrders: []int{2},
			wantNeedUnion:    true,
		},
		{
			name: "正常：複数のIN句で更新必要_2",
			cm: toConditionMap([]*ConditionValue{
				{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
				{order: 1, TargetColumn: "count", ConditionType: Eq, Values: []any{10}},
				{order: 2, TargetColumn: "item_id", ConditionType: In, Values: []any{1, 2, 3}},
				{order: 3, TargetColumn: "sample_id", ConditionType: In, Values: []any{1, 2, 3}},
			}),
			args: args{addableConditionMap: toConditionMap([]*ConditionValue{
				{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
				{order: 1, TargetColumn: "count", ConditionType: Eq, Values: []any{10}},
				{order: 2, TargetColumn: "item_id", ConditionType: In, Values: []any{2, 3, 4}},
				{order: 3, TargetColumn: "sample_id", ConditionType: In, Values: []any{1, 2, 5}},
			})},
			wantTargetOrders: []int{2, 3},
			wantNeedUnion:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.cm.needUnionConditionValue(tt.args.addableConditionMap)

			sort.Slice(got, func(i, j int) bool {
				return got[i] < got[j]
			})

			assert.Equal(t, got, tt.wantTargetOrders)
			assert.Equal(t, got1, tt.wantNeedUnion)
		})
	}
}

func TestConditionMaps_IsOneCalled(t *testing.T) {
	type args struct {
		conditions []*ConditionValue
	}
	tests := []struct {
		name string
		cms  ConditionMaps
		args args
		want bool
	}{
		{
			name: "一度実行済み_1つ目のCacheにHit",
			cms: ConditionMaps{
				toConditionMap([]*ConditionValue{
					{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
					{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{1, 2, 3}},
				}),
				toConditionMap([]*ConditionValue{
					{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{2}},
					{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{1, 2, 3}},
				}),
			},
			args: args{conditions: []*ConditionValue{
				{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
				{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{1, 2}},
			}},
			want: true,
		},
		{
			name: "一度実行済み_2つ目のCacheにHit",
			cms: ConditionMaps{
				toConditionMap([]*ConditionValue{
					{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
					{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{1, 2, 3}},
				}),
				toConditionMap([]*ConditionValue{
					{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{2}},
					{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{1, 2, 3}},
				}),
			},
			args: args{conditions: []*ConditionValue{
				{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{2}},
				{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{1, 3}},
			}},
			want: true,
		},
		{
			name: "実行したことがない_IN句",
			cms: ConditionMaps{
				toConditionMap([]*ConditionValue{
					{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
					{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{1, 2, 3}},
				}),
				toConditionMap([]*ConditionValue{
					{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{2}},
					{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{1, 2, 3}},
				}),
			},
			args: args{conditions: []*ConditionValue{
				{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
				{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{1, 3, 5}},
			}},
			want: false,
		},
		{
			name: "実行したことがない_Eq句",
			cms: ConditionMaps{
				toConditionMap([]*ConditionValue{
					{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
					{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{1, 2, 3}},
				}),
				toConditionMap([]*ConditionValue{
					{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{2}},
					{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{1, 2, 3}},
				}),
			},
			args: args{conditions: []*ConditionValue{
				{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{3}},
				{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{1, 2, 3}},
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cms.IsOnceCalled(tt.args.conditions); got != tt.want {
				t.Errorf("IsOnceCalled() = %v, wantKey %v", got, tt.want)
			}
		})
	}
}

func TestConditionMaps_UpdateConditionValues(t *testing.T) {
	type args struct {
		addableConditions []*ConditionValue
	}
	tests := []struct {
		name  string
		cms   ConditionMaps
		args  args
		want  bool
		want1 map[int][]any
	}{
		{
			name: "IN句の更新",
			cms: ConditionMaps{
				toConditionMap([]*ConditionValue{
					{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{int64(1)}},
					{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{int64(1), int64(2), int64(3)}},
				}),
				toConditionMap([]*ConditionValue{
					{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{int64(2)}},
					{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{int64(1), int64(2), int64(3)}},
				}),
			},
			args: args{addableConditions: []*ConditionValue{
				{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{int64(1)}},
				{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{int64(1), int64(3), int64(4), int64(5)}},
			}},
			want: true,
			want1: map[int][]any{
				0: {int64(1), int64(2), int64(3), int64(4), int64(5)},
				1: {int64(1), int64(2), int64(3)},
			},
		},
		{
			name: "Eq句が異なる場合",
			cms: ConditionMaps{
				toConditionMap([]*ConditionValue{
					{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{int64(1)}},
					{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{int64(1), int64(2), int64(3)}},
				}),
				toConditionMap([]*ConditionValue{
					{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{int64(2)}},
					{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{int64(1), int64(2), int64(3)}},
				}),
			},
			args: args{addableConditions: []*ConditionValue{
				{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{int64(3)}},
				{order: 1, TargetColumn: "item_id", ConditionType: In, Values: []any{int64(1), int64(3)}},
			}},
			want: false,
			want1: map[int][]any{
				0: {int64(1), int64(2), int64(3)},
				1: {int64(1), int64(2), int64(3)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cms.UnionConditionValues(tt.args.addableConditions)
			if got != tt.want {
				t.Errorf("UnionConditionValues() = %v, wantKey %v", got, tt.want)
			}

			gotValuesMap := make(map[int][]any)
			for i, conditionMap := range tt.cms {
				gotValuesMap[i] = conditionMap[1].Values
			}
			assert.Equal(t, tt.want1[0], gotValuesMap[0])
			assert.Equal(t, tt.want1[1], gotValuesMap[1])
		})
	}
}
