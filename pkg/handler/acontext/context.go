package acontext

import (
	"context"
	"errors"
	"game-server-example/pkg/handler/schema"
)

type key struct{}

func WithAPIResponse(ctx context.Context) context.Context {
	return context.WithValue(ctx, key{}, &schema.APIResponse{
		Common:   &schema.CommonResponse{},
		Original: &schema.OriginalResponse{},
	})
}

func ExtractAPIResponse(ctx context.Context) (*schema.APIResponse, error) {
	v, ok := ctx.Value(key{}).(*schema.APIResponse)
	if !ok {
		return nil, errors.New("適切な値がcontextに設定されていません")
	}
	return v, nil
}
