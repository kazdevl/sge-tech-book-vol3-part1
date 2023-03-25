package api

import (
	"game-server-example/pkg/handler/schema"
	"game-server-example/pkg/usecase"
	"github.com/labstack/echo/v4"
)

type MonsterHandler struct {
	monsterUsecase *usecase.MonsterUsecase
}

func NewMonsterHandler(monsterUsecase *usecase.MonsterUsecase) *MonsterHandler {
	return &MonsterHandler{monsterUsecase: monsterUsecase}
}

func (h *MonsterHandler) Enhance(ctx echo.Context, params schema.MonsterEnhanceParams) error {
	var req schema.MonsterEnhanceJSONBody
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	return h.monsterUsecase.Enhance(ctx.Request().Context(), params.UserId, *req.MonsterId, h.toUsecaseItemInfos((*[]struct {
		Count  *int64
		ItemId *int64
	})(req.Items)))
}

func (h *MonsterHandler) toUsecaseItemInfos(
	itemInfos *[]struct {
		Count  *int64
		ItemId *int64
	},
) []*usecase.ItemInfo {
	result := make([]*usecase.ItemInfo, 0, len(*itemInfos))
	for _, itemInfo := range *itemInfos {
		result = append(result, &usecase.ItemInfo{
			ItemId: *itemInfo.ItemId,
			Count:  *itemInfo.Count,
		})
	}
	return result
}
