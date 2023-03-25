package api

import (
	"game-server-example/pkg/entity"
	"game-server-example/pkg/handler/acontext"
	"game-server-example/pkg/handler/schema"
	"game-server-example/pkg/usecase"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) UserGetData(c echo.Context, params schema.UserGetDataParams) error {
	ctx := c.Request().Context()
	monsters, items, coin, err := h.userUsecase.GetData(ctx, params.UserId)
	if err != nil {
		return err
	}

	apiResponse, err := acontext.ExtractAPIResponse(ctx)
	if err != nil {
		return err
	}
	apiResponse.Common.Update = &struct {
		UserCoin    *[]schema.UserCoin    `json:"user_coin,omitempty"`
		UserItem    *[]schema.UserItem    `json:"user_item,omitempty"`
		UserMonster *[]schema.UserMonster `json:"user_monster,omitempty"`
	}{
		UserCoin:    h.toCoinResponses([]*entity.CoinEntity{coin}),
		UserItem:    h.toItemResponses(items),
		UserMonster: h.toMonsterResponses(monsters),
	}
	return nil
}

func (h *UserHandler) UserRegister(c echo.Context) error {
	ctx := c.Request().Context()
	userId, err := h.userUsecase.RegisterData(ctx)
	if err != nil {
		return err
	}

	apiResponse, err := acontext.ExtractAPIResponse(ctx)
	if err != nil {
		return err
	}
	apiResponse.Original.UserRegister = &schema.UserRegisterResponseContent{UserId: &userId}

	return nil
}

func (h *UserHandler) toCoinResponses(coins []*entity.CoinEntity) *[]schema.UserCoin {
	res := make([]schema.UserCoin, 0, len(coins))
	for _, coin := range coins {
		num := coin.Num()
		userId := coin.UserId()
		res = append(res, schema.UserCoin{
			Currency: &num,
			UserId:   &userId,
		})
	}
	return &res
}

func (h *UserHandler) toItemResponses(items []*entity.ItemEntity) *[]schema.UserItem {
	res := make([]schema.UserItem, 0, len(items))
	for _, item := range items {
		count := item.Count()
		userId := item.UserId()
		itemId := item.ItemId()
		res = append(res, schema.UserItem{
			Count:  &count,
			UserId: &userId,
			ItemId: &itemId,
		})
	}
	return &res
}

func (h *UserHandler) toMonsterResponses(monsters []*entity.MonsterEntity) *[]schema.UserMonster {
	res := make([]schema.UserMonster, 0, len(monsters))
	for _, monster := range monsters {
		exp := monster.Exp()
		userId := monster.UserId()
		monsterId := monster.MonsterId()
		res = append(res, schema.UserMonster{
			Exp:       &exp,
			UserId:    &userId,
			MonsterId: &monsterId,
		})
	}
	return &res
}
