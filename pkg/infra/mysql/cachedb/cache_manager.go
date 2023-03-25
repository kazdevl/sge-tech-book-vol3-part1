package cachedb

import "sync"

type DBOperationCacheManager struct {
	dbOperationResult    map[string]map[string]CacheContent  // key: table, uk_info
	selectQueryCondition map[string]map[string]ConditionMaps // key: table, query
	mutex                sync.RWMutex
}

func newDBOperationCacheManager() *DBOperationCacheManager {
	return &DBOperationCacheManager{
		dbOperationResult:    make(map[string]map[string]CacheContent, 0),
		selectQueryCondition: make(map[string]map[string]ConditionMaps, 0),
	}
}

func (t *DBOperationCacheManager) reset() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.dbOperationResult = make(map[string]map[string]CacheContent)
	t.selectQueryCondition = make(map[string]map[string]ConditionMaps)
}

func (t *DBOperationCacheManager) haveUnSyncedCache() bool {
	for _, tableQueryResultMap := range t.dbOperationResult {
		for _, cachedContent := range tableQueryResultMap {
			cacheStatus := cachedContent.GetCacheStatus()
			if cacheStatus == Update || cacheStatus == Insert || cacheStatus == Delete {
				return true
			}
		}
	}
	return false
}
