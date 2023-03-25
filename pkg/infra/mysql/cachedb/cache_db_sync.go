package cachedb

import (
	"context"
	"errors"
)

func (cdb *CacheDB) Begin() error {
	tx, err := cdb.db.Beginx()
	if err != nil {
		return err
	}
	cdb.tx = tx
	return nil
}

func (cdb *CacheDB) Commit(ctx context.Context) error {
	cacheManager, err := extractDBOperationCacheManager(ctx)
	if err != nil {
		return err
	}
	if cacheManager.haveUnSyncedCache() {
		return errors.New("同期漏れしている変更差分が存在します")
	}
	if cdb.tx == nil {
		return errors.New("transactionが設定されていません")
	}

	cacheManager.reset()
	return cdb.tx.Commit()
}

func (cdb *CacheDB) Rollback(ctx context.Context) error {
	cacheManager, err := extractDBOperationCacheManager(ctx)
	if err != nil {
		return err
	}
	if cdb.tx == nil {
		return errors.New("transactionが設定されていません")
	}

	cacheManager.reset()
	return cdb.tx.Rollback()
}

func (cdb *CacheDB) SyncedToDB(ctx context.Context) (map[string]map[CacheStatus][]CacheContent, error) {
	if cdb.tx == nil {
		if err := cdb.Begin(); err != nil {
			return nil, err
		}
	}

	cacheManager, err := extractDBOperationCacheManager(ctx)
	if err != nil {
		return nil, err
	}

	// 1. キャッシュ一覧をテーブルとキャッシュ状態をキーにしたマップ型に変換する
	cachesMapByStatus := make(map[string]map[CacheStatus][]CacheContent)
	for table, resultMap := range cacheManager.dbOperationResult {
		cachesMapByStatus[table] = make(map[CacheStatus][]CacheContent)
		for _, result := range resultMap {
			cacheStatus := result.GetCacheStatus()
			if cacheStatus == Select || cacheStatus == None {
				continue
			}

			if cachesMapByStatus[table][cacheStatus] == nil {
				cachesMapByStatus[table][cacheStatus] = make([]CacheContent, 0)
			}
			cachesMapByStatus[table][cacheStatus] = append(cachesMapByStatus[table][cacheStatus], result.CreateCopy())

			// 処理済みとしてcontextにcacheされているcacheの内容のstatusをNoneに設定する
			result.SetCacheStatus(None)
		}
	}

	// 2. modelBulkExecuterMapを利用して、各テーブルごとにDB操作を実行する
	for table, cachesMap := range cachesMapByStatus {
		for cacheStatus, caches := range cachesMap {
			switch cacheStatus {
			case Insert:
				if err = cdb.modelBulkExecutorMap[table].BulkInsert(ctx, cdb.tx, caches); err != nil {
					return nil, err
				}
			case Update:
				if err = cdb.modelBulkExecutorMap[table].BulkUpdate(ctx, cdb.tx, caches); err != nil {
					return nil, err
				}
			case Delete:
				if err = cdb.modelBulkExecutorMap[table].BulkDelete(ctx, cdb.tx, caches); err != nil {
					return nil, err
				}
			}
		}
	}

	return cachesMapByStatus, nil
}
