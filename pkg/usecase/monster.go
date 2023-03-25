package usecase

import (
	"context"
	"game-server-example/pkg/service"
	"github.com/samber/lo"
)

type MonsterUsecase struct {
	monsterEnhanceService   *service.MonsterEnhanceService
	checkPerformanceService *service.CheckPerformanceService
}

type ItemInfo struct {
	ItemId int64
	Count  int64
}

func NewMonsterUsecase(monsterEnhanceService *service.MonsterEnhanceService, checkPerformanceService *service.CheckPerformanceService) *MonsterUsecase {
	return &MonsterUsecase{monsterEnhanceService: monsterEnhanceService, checkPerformanceService: checkPerformanceService}
}

func (u *MonsterUsecase) Enhance(ctx context.Context, userId, monsterId int64, itemInfos []*ItemInfo) error {
	err := u.monsterEnhanceService.Enhance(ctx, userId, monsterId, lo.Map(itemInfos, func(item *ItemInfo, _ int) *service.ItemInfo {
		return &service.ItemInfo{
			ItemId: item.ItemId,
			Count:  item.Count,
		}
	}))
	if err != nil {
		return err
	}

	return u.checkPerformanceService.CheckPerformance(ctx, userId, monsterId, lo.Map(itemInfos, func(itemInfo *ItemInfo, _ int) int64 {
		return itemInfo.ItemId
	}))
}
