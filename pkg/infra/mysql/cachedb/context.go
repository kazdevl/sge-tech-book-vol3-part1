package cachedb

import (
	"context"
	"errors"
)

type cacheKey struct{}

func WithCacheContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, cacheKey{}, newDBOperationCacheManager())
}

func extractDBOperationCacheManager(ctx context.Context) (*DBOperationCacheManager, error) {
	v := ctx.Value(cacheKey{})
	if v == nil {
		return nil, errors.New("contextに値が設定されていません")
	}
	cache, ok := v.(*DBOperationCacheManager)
	if !ok {
		return nil, errors.New("想定している型ではありません")
	}
	return cache, nil
}
