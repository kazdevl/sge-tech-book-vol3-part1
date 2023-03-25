package cachedb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"sort"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

type Model struct {
	UserId      int64
	SampleId    int64
	Count       int64
	cacheStatus CacheStatus
}

func (m *Model) Bind(rows *sql.Rows) error {
	return rows.Scan(&m.UserId, &m.SampleId, &m.Count)
}

func (m *Model) UniqueKeyKCondition() string {
	return "user_id=? AND sample_id=?"
}

func (m *Model) UniqueKeyConditionValues() []any {
	return []any{m.UserId, m.SampleId}
}

func (m *Model) Table() string {
	return "sample"
}

func (m *Model) UniqueKeyColumnValueStr() string {
	return fmt.Sprintf("user_id=%d,sample_id=%d", m.UserId, m.SampleId)
}

func (m *Model) SetCacheStatus(s CacheStatus) {
	m.cacheStatus = s
}

func (m *Model) GetCacheStatus() CacheStatus {
	return m.cacheStatus
}

func (m *Model) Update(content CacheContent) error {
	modelContent, ok := content.(*Model)
	if !ok {
		return errors.New("想定していない型で更新しようとしています")
	}
	m.UserId = modelContent.UserId
	m.SampleId = modelContent.SampleId
	m.Count = modelContent.Count
	return nil
}

func (m *Model) CreateCopy() CacheContent {
	newModel := &Model{
		UserId:      m.UserId,
		SampleId:    m.SampleId,
		Count:       m.Count,
		cacheStatus: m.cacheStatus,
	}
	return newModel
}

func (m *Model) IsSame(column string, v any) bool {
	switch column {
	case "user_id":
		userId, ok := v.(int64)
		if !ok {
			return false
		}
		return m.UserId == userId
	case "sample_id":
		sampleId, ok := v.(int64)
		if !ok {
			return false
		}
		return m.SampleId == sampleId
	case "count":
		count, ok := v.(int64)
		if !ok {
			return false
		}
		return m.Count == count
	default:
		return false
	}
}

func (m *Model) IsInclude(column string, vs []any) bool {
	switch column {
	case "user_id":
		userIds := toInt64Slice(vs)
		if len(userIds) == 0 {
			return false
		}
		return lo.Contains(userIds, m.UserId)
	case "sample_id":
		sampleIds := toInt64Slice(vs)
		if len(sampleIds) == 0 {
			return false
		}
		return lo.Contains(sampleIds, m.SampleId)
	case "count":
		counts := toInt64Slice(vs)
		if len(counts) == 0 {
			return false
		}
		return lo.Contains(counts, m.Count)
	default:
		return false
	}
}

func toInt64Slice(vs []any) []int64 {
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

type Models []*Model

func (ms Models) Table() string {
	return "sample"
}

func (ms *Models) BindAndAddContent(rows *sql.Rows) error {
	m := &Model{}
	if err := rows.Scan(&m.UserId, &m.SampleId, &m.Count); err != nil {
		return err
	}
	*ms = append(*ms, m)

	return nil
}

func (ms Models) ForEach(f func(cacheContent CacheContent)) {
	for _, m := range ms {
		f(m)
	}
}

func (ms *Models) Add(content CacheContent) {
	contentModel, ok := content.(*Model)
	if ok {
		*ms = append(*ms, contentModel)
	}
}

func (ms *Models) Reset() {
	*ms = make(Models, 0)
}

func TestCacheDB_Insert(t *testing.T) {
	ctx := context.Background()
	ctx = WithCacheContext(ctx)

	type fields struct {
		db *sqlx.DB
		tx *sqlx.Tx
	}
	type args struct {
		getCtx  func() context.Context
		content CacheContent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Model
		wantErr bool
	}{
		{
			name:   "正常:cacheが存在しなかった場合",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 10},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 10, cacheStatus: Insert},
			wantErr: false,
		},
		{
			name:   "正常:cacheが存在している場合(Deleteの状態)",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Delete}
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 10},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 10, cacheStatus: Update},
			wantErr: false,
		},
		{
			name:   "正常:cacheが存在している場合(Noneの状態)",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: None}
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 10},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 10, cacheStatus: Insert},
			wantErr: false,
		},
		{
			name:   "異常:cacheが存在している場合(Insertの状態)",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Insert}
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 10},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Insert},
			wantErr: true,
		},
		{
			name:   "異常:cacheが存在している場合(Selectの状態)",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Select}
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 10},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Select},
			wantErr: true,
		},
		{
			name:   "異常:cacheが存在している場合(Updateの状態)",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Update}
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 10},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Update},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdb := &CacheDB{
				db: tt.fields.db,
				tx: tt.fields.tx,
			}
			err := cdb.Insert(tt.args.getCtx(), tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			cache, _ := extractDBOperationCacheManager(ctx)
			assert.Equal(t, tt.want, cache.dbOperationResult["sample"]["user_id=1,sample_id=1"])
			cache.reset()
		})
	}
}

func TestCacheDB_Update(t *testing.T) {
	ctx := context.Background()
	ctx = WithCacheContext(ctx)

	type fields struct {
		db *sqlx.DB
		tx *sqlx.Tx
	}
	type args struct {
		getCtx  func() context.Context
		content CacheContent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    CacheContent
		wantErr bool
	}{
		{
			name:   "正常:cacheが存在している場合(Selectの状態)",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Select}
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 10},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 10, cacheStatus: Update},
			wantErr: false,
		},
		{
			name:   "正常:cacheが存在している場合(Updateの状態)",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Update}
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 10},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 10, cacheStatus: Update},
			wantErr: false,
		},
		{
			name:   "正常:cacheが存在している場合(Insertの状態)",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Insert}
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 10},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 10, cacheStatus: Insert},
			wantErr: false,
		},
		{
			name:   "異常:cacheが存在している場合(Deleteの状態)",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Delete}
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 10},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Delete},
			wantErr: true,
		},
		{
			name:   "異常:cacheが存在している場合(Noneの状態)",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: None}
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 10},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: None},
			wantErr: true,
		},
		{
			name:   "異常:cacheが存在しなかった場合",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 10},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdb := &CacheDB{
				db: tt.fields.db,
				tx: tt.fields.tx,
			}
			err := cdb.Update(tt.args.getCtx(), tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			cache, _ := extractDBOperationCacheManager(ctx)
			assert.Equal(t, tt.want, cache.dbOperationResult["sample"]["user_id=1,sample_id=1"])

			cache.reset()
		})
	}
}

func TestCacheDB_Delete(t *testing.T) {
	ctx := context.Background()
	ctx = WithCacheContext(ctx)

	type fields struct {
		db *sqlx.DB
		tx *sqlx.Tx
	}
	type args struct {
		getCtx  func() context.Context
		content CacheContent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    CacheContent
		wantErr bool
	}{
		{
			name:   "正常:cacheが存在している場合(Selectの状態)",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Select}
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 5},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Delete},
			wantErr: false,
		},
		{
			name:   "正常:cacheが存在している場合(Updateの状態)",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Update}
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 5},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Delete},
			wantErr: false,
		},
		{
			name:   "正常:cacheが存在している場合(Insertの状態)",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Insert}
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 5},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: None},
			wantErr: false,
		},
		{
			name:   "異常:cacheが存在している場合(Deleteの状態)",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Delete}
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 5},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Delete},
			wantErr: false,
		},
		{
			name:   "異常:cacheが存在している場合(Noneの状態)",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: None}
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 5},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: None},
			wantErr: true,
		},
		{
			name:   "異常:cacheが存在しなかった場合",
			fields: fields{},
			args: args{
				getCtx: func() context.Context {
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1, Count: 5},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdb := &CacheDB{
				db: tt.fields.db,
				tx: tt.fields.tx,
			}
			err := cdb.Delete(tt.args.getCtx(), tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			cache, _ := extractDBOperationCacheManager(ctx)
			assert.Equal(t, tt.want, cache.dbOperationResult["sample"]["user_id=1,sample_id=1"])

			cache.reset()
		})
	}
}

func TestCacheDB_GetAndSet(t *testing.T) {
	ctx := context.Background()
	ctx = WithCacheContext(ctx)

	type fields struct {
		getDB func() *sqlx.DB
		tx    *sqlx.Tx
	}
	type args struct {
		getCtx  func() context.Context
		content SelectCacheContent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    SelectCacheContent
		wantErr bool
	}{
		{
			name: "正常:cacheから取得",
			fields: fields{
				getDB: func() *sqlx.DB {
					return nil
				},
			},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Select}
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Select},
			wantErr: false,
		},
		{
			name: "正常:dbから取得してcacheに設定",
			fields: fields{
				getDB: func() *sqlx.DB {
					db, mock, err := sqlmock.New()
					assert.NoError(t, err)

					rows := sqlmock.NewRows([]string{"user_id", "sample_id", "count"}).
						AddRow(1, 1, 5)
					mock.ExpectQuery(`SELECT \* FROM sample WHERE user_id=\? AND sample_id=\?`).
						WithArgs(1, 1).
						WillReturnRows(rows)
					return sqlx.NewDb(db, "mysql")
				},
			},
			args: args{
				getCtx: func() context.Context {
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1},
			},
			want:    &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Select},
			wantErr: false,
		},
		{
			name: "正常:Delete状態のcacheから取得しようとしている",
			fields: fields{
				getDB: func() *sqlx.DB {
					return nil
				},
			},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Delete}
					return ctx
				},
				content: &Model{UserId: 1, SampleId: 1},
			},
			want:    &Model{UserId: 1, SampleId: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdb := &CacheDB{
				db: tt.fields.getDB(),
				tx: tt.fields.tx,
			}
			err := cdb.GetAndSet(tt.args.getCtx(), tt.args.content)
			if err != nil {
				t.Log(err)
			}
			assert.Equal(t, tt.wantErr, err != nil)

			assert.Equal(t, tt.want, tt.args.content)

			cache, _ := extractDBOperationCacheManager(ctx)
			cache.reset()
		})
	}
}

func TestCacheDB_FindAndSetByConditions(t *testing.T) {
	ctx := context.Background()
	ctx = WithCacheContext(ctx)

	type fields struct {
		getDB func() *sqlx.DB
		tx    *sqlx.Tx
	}
	type args struct {
		getCtx     func() context.Context
		contents   SelectCacheContents
		conditions []*ConditionValue
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    SelectCacheContents
		wantErr bool
	}{
		{
			name: "正常:cacheから取得",
			fields: fields{
				getDB: func() *sqlx.DB {
					return nil
				},
			},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.selectQueryCondition["sample"] = make(map[string]ConditionMaps)
					cache.selectQueryCondition["sample"]["SELECT * FROM sample WHERE user_id = ? AND sample_id IN (?)"] = ConditionMaps{
						toConditionMap([]*ConditionValue{
							{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
							{order: 1, TargetColumn: "sample_id", ConditionType: In, Values: []any{1, 2, 3}},
						}),
					}
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Select}
					cache.dbOperationResult["sample"]["user_id=1,sample_id=2"] = &Model{UserId: 1, SampleId: 2, Count: 5, cacheStatus: Select}
					cache.dbOperationResult["sample"]["user_id=1,sample_id=3"] = &Model{UserId: 1, SampleId: 3, Count: 5, cacheStatus: Insert}
					return ctx
				},
				contents: &Models{},
				conditions: []*ConditionValue{
					{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
					{order: 1, TargetColumn: "sample_id", ConditionType: In, Values: []any{1, 2, 3}},
				},
			},
			want: &Models{
				{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Select},
				{UserId: 1, SampleId: 2, Count: 5, cacheStatus: Select},
				{UserId: 1, SampleId: 3, Count: 5, cacheStatus: Insert},
			},
			wantErr: false,
		},
		{
			name: "正常:dbから取得してcacheに設定(同じ条件で検索したことがない)",
			fields: fields{
				getDB: func() *sqlx.DB {
					db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					assert.NoError(t, err)

					rows := sqlmock.NewRows([]string{"user_id", "sample_id", "count"}).
						AddRow(1, 1, 5).
						AddRow(1, 2, 5).
						AddRow(1, 3, 5)
					mock.ExpectQuery("SELECT * FROM sample WHERE user_id = ? AND sample_id IN (?,?,?)").
						WithArgs(1, 1, 2, 3).
						WillReturnRows(rows)
					return sqlx.NewDb(db, "mysql")
				},
			},
			args: args{
				getCtx: func() context.Context {
					return ctx
				},
				contents: &Models{},
				conditions: []*ConditionValue{
					{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{int64(1)}},
					{order: 1, TargetColumn: "sample_id", ConditionType: In, Values: []any{int64(1), int64(2), int64(3)}},
				},
			},
			want: &Models{
				{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Select},
				{UserId: 1, SampleId: 2, Count: 5, cacheStatus: Select},
				{UserId: 1, SampleId: 3, Count: 5, cacheStatus: Select},
			},
			wantErr: false,
		},
		{
			name: "正常:dbから取得してcacheに設定(同じ条件で検索したことはあるが、検索値が異なる)",
			fields: fields{
				getDB: func() *sqlx.DB {
					db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					assert.NoError(t, err)

					rows := sqlmock.NewRows([]string{"user_id", "sample_id", "count"}).
						AddRow(1, 1, 5).
						AddRow(1, 2, 5).
						AddRow(1, 3, 5).
						AddRow(1, 4, 5)
					mock.ExpectQuery("SELECT * FROM sample WHERE user_id = ? AND sample_id IN (?,?,?,?)").
						WithArgs(1, 1, 2, 3, 4).
						WillReturnRows(rows)
					return sqlx.NewDb(db, "mysql")
				},
			},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					cache.dbOperationResult["sample"] = make(map[string]CacheContent)
					cache.selectQueryCondition["sample"] = make(map[string]ConditionMaps)
					cache.selectQueryCondition["sample"]["SELECT * FROM sample WHERE user_id = ? AND sample_id IN (?)"] = ConditionMaps{
						toConditionMap([]*ConditionValue{
							{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{1}},
							{order: 1, TargetColumn: "sample_id", ConditionType: In, Values: []any{1, 2, 3}},
						}),
					}
					cache.dbOperationResult["sample"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Select}
					cache.dbOperationResult["sample"]["user_id=1,sample_id=2"] = &Model{UserId: 1, SampleId: 2, Count: 5, cacheStatus: Select}
					cache.dbOperationResult["sample"]["user_id=1,sample_id=3"] = &Model{UserId: 1, SampleId: 3, Count: 5, cacheStatus: Insert}
					return ctx
				},
				contents: &Models{},
				conditions: []*ConditionValue{
					{order: 0, TargetColumn: "user_id", ConditionType: Eq, Values: []any{int64(1)}},
					{order: 1, TargetColumn: "sample_id", ConditionType: In, Values: []any{int64(1), int64(2), int64(3), int64(4)}},
				},
			},
			want: &Models{
				{UserId: 1, SampleId: 1, Count: 5, cacheStatus: Select},
				{UserId: 1, SampleId: 2, Count: 5, cacheStatus: Select},
				{UserId: 1, SampleId: 3, Count: 5, cacheStatus: Insert},
				{UserId: 1, SampleId: 4, Count: 5, cacheStatus: Select},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdb := &CacheDB{
				db: tt.fields.getDB(),
				tx: tt.fields.tx,
			}
			err := cdb.FindAndSetByConditions(tt.args.getCtx(), tt.args.contents, tt.args.conditions)
			if err != nil {
				t.Log(err)
			}
			assert.Equal(t, tt.wantErr, err != nil)

			wantModels := tt.want.(*Models)
			argContents := tt.args.contents.(*Models)
			sort.Slice(*wantModels, func(i, j int) bool {
				wantModelNoPointer := *wantModels
				return wantModelNoPointer[i].SampleId < wantModelNoPointer[j].SampleId
			})
			sort.Slice(*argContents, func(i, j int) bool {
				argContentNoPointer := *argContents
				return argContentNoPointer[i].SampleId < argContentNoPointer[j].SampleId
			})
			assert.Equal(t, *wantModels, *argContents)

			cache, _ := extractDBOperationCacheManager(ctx)
			cache.reset()
		})
	}
}
