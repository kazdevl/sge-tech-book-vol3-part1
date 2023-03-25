package cachedb

import "database/sql"

type CacheStatus int64

const (
	None CacheStatus = iota
	Select
	Insert
	Update
	Delete
)

type CacheContent interface {
	Table() string
	UniqueKeyColumnValueStr() string
	SetCacheStatus(s CacheStatus)
	GetCacheStatus() CacheStatus
	Update(content CacheContent) error
	CreateCopy() CacheContent
	IsSame(column string, v any) bool
	IsInclude(column string, vs []any) bool
}

type SelectCacheContent interface {
	CacheContent
	Bind(*sql.Rows) error
	UniqueKeyKCondition() string
	UniqueKeyConditionValues() []any
}

type SelectCacheContents interface {
	Table() string
	BindAndAddContent(*sql.Rows) error
	ForEach(func(cacheContent CacheContent))
	Add(content CacheContent)
	Reset()
}
