package cachedb

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"reflect"
	"sort"
	"testing"
)

func TestCacheDB_SyncedToDB(t *testing.T) {
	ctx := context.Background()
	ctx = WithCacheContext(ctx)

	type fields struct {
		db                      *sqlx.DB
		tx                      *sqlx.Tx
		getModelBulkExecutorMap func(ctrl *gomock.Controller) map[string]BulkExecutor
	}
	type args struct {
		getCtx func() context.Context
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		want            map[string]map[CacheStatus][]CacheContent
		wantCacheStatus map[string]map[string]CacheContent
		wantErr         bool
	}{
		{
			name: "データの同期成功(Insert・Update・Delete)",
			fields: fields{
				getModelBulkExecutorMap: func(ctrl *gomock.Controller) map[string]BulkExecutor {
					mockForSample1 := NewMockBulkExecutor(ctrl)
					mockForSample1.EXPECT().BulkInsert(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).DoAndReturn(func(ctx context.Context, tx *sqlx.Tx, contents []CacheContent) error {
						expect := []CacheContent{
							&Model{UserId: 1, SampleId: 3, Count: 5, cacheStatus: Insert},
							&Model{UserId: 1, SampleId: 4, Count: 5, cacheStatus: Insert},
						}
						toSortContents(contents)

						assert.True(t, reflect.DeepEqual(expect, contents))
						return nil
					})
					mockForSample1.EXPECT().BulkUpdate(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).DoAndReturn(func(ctx context.Context, tx *sqlx.Tx, contents []CacheContent) error {
						expect := []CacheContent{
							&Model{UserId: 1, SampleId: 5, Count: 5, cacheStatus: Update},
							&Model{UserId: 1, SampleId: 6, Count: 5, cacheStatus: Update},
						}
						toSortContents(contents)

						assert.True(t, reflect.DeepEqual(expect, contents))
						return nil
					})
					mockForSample1.EXPECT().BulkDelete(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).DoAndReturn(func(ctx context.Context, tx *sqlx.Tx, contents []CacheContent) error {
						expect := []CacheContent{
							&Model{UserId: 1, SampleId: 7, Count: 5, cacheStatus: Delete},
							&Model{UserId: 1, SampleId: 8, Count: 5, cacheStatus: Delete},
						}
						toSortContents(contents)

						assert.True(t, reflect.DeepEqual(expect, contents))
						return nil
					})

					mockForSample2 := NewMockBulkExecutor(ctrl)
					mockForSample2.EXPECT().BulkInsert(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).DoAndReturn(func(ctx context.Context, tx *sqlx.Tx, contents []CacheContent) error {
						expect := []CacheContent{
							&Model{UserId: 1, SampleId: 3, Count: 10, cacheStatus: Insert},
							&Model{UserId: 1, SampleId: 4, Count: 10, cacheStatus: Insert},
						}
						toSortContents(contents)

						assert.True(t, reflect.DeepEqual(expect, contents))
						return nil
					})
					mockForSample2.EXPECT().BulkUpdate(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).DoAndReturn(func(ctx context.Context, tx *sqlx.Tx, contents []CacheContent) error {
						expect := []CacheContent{
							&Model{UserId: 1, SampleId: 5, Count: 10, cacheStatus: Update},
							&Model{UserId: 1, SampleId: 6, Count: 10, cacheStatus: Update},
						}
						toSortContents(contents)

						assert.True(t, reflect.DeepEqual(expect, contents))
						return nil
					})
					mockForSample2.EXPECT().BulkDelete(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).DoAndReturn(func(ctx context.Context, tx *sqlx.Tx, contents []CacheContent) error {
						expect := []CacheContent{
							&Model{UserId: 1, SampleId: 7, Count: 10, cacheStatus: Delete},
							&Model{UserId: 1, SampleId: 8, Count: 10, cacheStatus: Delete},
						}
						toSortContents(contents)

						assert.True(t, reflect.DeepEqual(expect, contents))
						return nil
					})

					return map[string]BulkExecutor{"sample1": mockForSample1, "sample2": mockForSample2}
				},
				tx: &sqlx.Tx{},
			},
			args: args{
				getCtx: func() context.Context {
					cache, _ := extractDBOperationCacheManager(ctx)
					// sample1テーブルのcache
					cache.dbOperationResult["sample1"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample1"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: None}
					cache.dbOperationResult["sample1"]["user_id=1,sample_id=2"] = &Model{UserId: 1, SampleId: 2, Count: 5, cacheStatus: Select}
					cache.dbOperationResult["sample1"]["user_id=1,sample_id=3"] = &Model{UserId: 1, SampleId: 3, Count: 5, cacheStatus: Insert}
					cache.dbOperationResult["sample1"]["user_id=1,sample_id=4"] = &Model{UserId: 1, SampleId: 4, Count: 5, cacheStatus: Insert}
					cache.dbOperationResult["sample1"]["user_id=1,sample_id=5"] = &Model{UserId: 1, SampleId: 5, Count: 5, cacheStatus: Update}
					cache.dbOperationResult["sample1"]["user_id=1,sample_id=6"] = &Model{UserId: 1, SampleId: 6, Count: 5, cacheStatus: Update}
					cache.dbOperationResult["sample1"]["user_id=1,sample_id=7"] = &Model{UserId: 1, SampleId: 7, Count: 5, cacheStatus: Delete}
					cache.dbOperationResult["sample1"]["user_id=1,sample_id=8"] = &Model{UserId: 1, SampleId: 8, Count: 5, cacheStatus: Delete}

					// sample2テーブルのcache
					cache.dbOperationResult["sample2"] = make(map[string]CacheContent)
					cache.dbOperationResult["sample2"]["user_id=1,sample_id=1"] = &Model{UserId: 1, SampleId: 1, Count: 10, cacheStatus: None}
					cache.dbOperationResult["sample2"]["user_id=1,sample_id=2"] = &Model{UserId: 1, SampleId: 2, Count: 10, cacheStatus: Select}
					cache.dbOperationResult["sample2"]["user_id=1,sample_id=3"] = &Model{UserId: 1, SampleId: 3, Count: 10, cacheStatus: Insert}
					cache.dbOperationResult["sample2"]["user_id=1,sample_id=4"] = &Model{UserId: 1, SampleId: 4, Count: 10, cacheStatus: Insert}
					cache.dbOperationResult["sample2"]["user_id=1,sample_id=5"] = &Model{UserId: 1, SampleId: 5, Count: 10, cacheStatus: Update}
					cache.dbOperationResult["sample2"]["user_id=1,sample_id=6"] = &Model{UserId: 1, SampleId: 6, Count: 10, cacheStatus: Update}
					cache.dbOperationResult["sample2"]["user_id=1,sample_id=7"] = &Model{UserId: 1, SampleId: 7, Count: 10, cacheStatus: Delete}
					cache.dbOperationResult["sample2"]["user_id=1,sample_id=8"] = &Model{UserId: 1, SampleId: 8, Count: 10, cacheStatus: Delete}
					return ctx
				},
			},
			want: map[string]map[CacheStatus][]CacheContent{
				"sample1": {
					Insert: []CacheContent{
						&Model{UserId: 1, SampleId: 3, Count: 5, cacheStatus: Insert},
						&Model{UserId: 1, SampleId: 4, Count: 5, cacheStatus: Insert},
					},
					Update: []CacheContent{
						&Model{UserId: 1, SampleId: 5, Count: 5, cacheStatus: Update},
						&Model{UserId: 1, SampleId: 6, Count: 5, cacheStatus: Update},
					},
					Delete: []CacheContent{
						&Model{UserId: 1, SampleId: 7, Count: 5, cacheStatus: Delete},
						&Model{UserId: 1, SampleId: 8, Count: 5, cacheStatus: Delete},
					},
				},
				"sample2": {
					Insert: []CacheContent{
						&Model{UserId: 1, SampleId: 3, Count: 10, cacheStatus: Insert},
						&Model{UserId: 1, SampleId: 4, Count: 10, cacheStatus: Insert},
					},
					Update: []CacheContent{
						&Model{UserId: 1, SampleId: 5, Count: 10, cacheStatus: Update},
						&Model{UserId: 1, SampleId: 6, Count: 10, cacheStatus: Update},
					},
					Delete: []CacheContent{
						&Model{UserId: 1, SampleId: 7, Count: 10, cacheStatus: Delete},
						&Model{UserId: 1, SampleId: 8, Count: 10, cacheStatus: Delete},
					},
				},
			},
			wantCacheStatus: map[string]map[string]CacheContent{
				"sample1": {
					"user_id=1,sample_id=1": &Model{UserId: 1, SampleId: 1, Count: 5, cacheStatus: None},
					"user_id=1,sample_id=2": &Model{UserId: 1, SampleId: 2, Count: 5, cacheStatus: Select},
					"user_id=1,sample_id=3": &Model{UserId: 1, SampleId: 3, Count: 5, cacheStatus: None},
					"user_id=1,sample_id=4": &Model{UserId: 1, SampleId: 4, Count: 5, cacheStatus: None},
					"user_id=1,sample_id=5": &Model{UserId: 1, SampleId: 5, Count: 5, cacheStatus: None},
					"user_id=1,sample_id=6": &Model{UserId: 1, SampleId: 6, Count: 5, cacheStatus: None},
					"user_id=1,sample_id=7": &Model{UserId: 1, SampleId: 7, Count: 5, cacheStatus: None},
					"user_id=1,sample_id=8": &Model{UserId: 1, SampleId: 8, Count: 5, cacheStatus: None},
				},
				"sample2": {
					"user_id=1,sample_id=1": &Model{UserId: 1, SampleId: 1, Count: 10, cacheStatus: None},
					"user_id=1,sample_id=2": &Model{UserId: 1, SampleId: 2, Count: 10, cacheStatus: Select},
					"user_id=1,sample_id=3": &Model{UserId: 1, SampleId: 3, Count: 10, cacheStatus: None},
					"user_id=1,sample_id=4": &Model{UserId: 1, SampleId: 4, Count: 10, cacheStatus: None},
					"user_id=1,sample_id=5": &Model{UserId: 1, SampleId: 5, Count: 10, cacheStatus: None},
					"user_id=1,sample_id=6": &Model{UserId: 1, SampleId: 6, Count: 10, cacheStatus: None},
					"user_id=1,sample_id=7": &Model{UserId: 1, SampleId: 7, Count: 10, cacheStatus: None},
					"user_id=1,sample_id=8": &Model{UserId: 1, SampleId: 8, Count: 10, cacheStatus: None},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cdb := &CacheDB{
				db:                   tt.fields.db,
				tx:                   tt.fields.tx,
				modelBulkExecutorMap: tt.fields.getModelBulkExecutorMap(ctrl),
			}
			got, err := cdb.SyncedToDB(tt.args.getCtx())
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equalf(t, tt.want, got, "SyncedToDB(ctx)")

			cache, _ := extractDBOperationCacheManager(ctx)
			assert.Equal(t, tt.wantCacheStatus, cache.dbOperationResult)

			cache.reset()
		})
	}
}

func toSortContents(contents []CacheContent) {
	contentsModel := make(Models, 0, len(contents))
	for _, content := range contents {
		contentsModel = append(contentsModel, content.(*Model))
	}
	sort.Slice(contentsModel, func(i, j int) bool {
		return contentsModel[i].SampleId < contentsModel[j].SampleId
	})

	for index, contentModel := range contentsModel {
		contents[index] = contentModel
	}
}
