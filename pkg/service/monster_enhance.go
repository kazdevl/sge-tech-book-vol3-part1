package service

import (
	"context"
	"errors"
	"game-server-example/pkg/entity"
	"game-server-example/pkg/ifrepository"
	"game-server-example/pkg/masterdata"
	"github.com/samber/lo"
)

type MonsterEnhanceService struct {
	monsterEntityRepository ifrepository.IFMonsterEntityRepository
	itemEntityRepository    ifrepository.IFItemEntityRepository
	coinEntityRepository    ifrepository.IFCoinEntityRepository
}

func NewMonsterEnhanceService(monsterEntityRepository ifrepository.IFMonsterEntityRepository, itemEntityRepository ifrepository.IFItemEntityRepository, coinEntityRepository ifrepository.IFCoinEntityRepository) *MonsterEnhanceService {
	return &MonsterEnhanceService{monsterEntityRepository: monsterEntityRepository, itemEntityRepository: itemEntityRepository, coinEntityRepository: coinEntityRepository}
}

type ItemInfo struct {
	ItemId int64
	Count  int64
}

type ItemInfos []*ItemInfo

func (is ItemInfos) ToMap() map[int64]int64 {
	result := make(map[int64]int64, len(is))
	for _, i := range is {
		result[i.ItemId] = i.Count
	}
	return result
}

func (is ItemInfos) ToIds() []int64 {
	return lo.Map(is, func(item *ItemInfo, _ int) int64 {
		return item.ItemId
	})
}

func (s *MonsterEnhanceService) Enhance(ctx context.Context, userId, monsterId int64, itemInfos ItemInfos) error {
	totalExp, err := s.consumeEnhanceItems(ctx, userId, itemInfos)
	if err != nil {
		return err
	}

	monsterEntity, err := s.monsterEntityRepository.Get(ctx, userId, monsterId)
	if err != nil {
		return err
	}

	if !monsterEntity.CanAddExp() {
		return errors.New("レベル上限に達しています")
	}

	beforeLevel := monsterEntity.Level()

	monsterEntity.AddExp(totalExp)
	if err = s.monsterEntityRepository.Save(ctx, monsterEntity); err != nil {
		return err
	}

	afterLevel := monsterEntity.Level()
	var totalConsumeCoin int64
	for level := beforeLevel + 1; level <= afterLevel; level++ {
		enhanceInfo := masterdata.Master().MonsterEnhanceTable.Get(level)
		totalConsumeCoin += enhanceInfo.Coin
	}

	coinEntity, err := s.coinEntityRepository.Get(ctx, userId)
	if err != nil {
		return err
	}
	if err = coinEntity.Consume(totalConsumeCoin); err != nil {
		return err
	}
	return s.coinEntityRepository.Save(ctx, coinEntity)
}

func (s *MonsterEnhanceService) consumeEnhanceItems(ctx context.Context, userId int64, itemInfos ItemInfos) (totalExp int64, err error) {
	itemIds := itemInfos.ToIds()
	itemConsumeMap := itemInfos.ToMap()

	if err = s.validateItemType(itemIds); err != nil {
		return 0, err
	}

	es, err := s.itemEntityRepository.FindByItemIds(ctx, userId, itemIds)
	if err != nil {
		return 0, err
	}

	lo.ForEach(es, func(e *entity.ItemEntity, _ int) {
		if err = e.Consume(itemConsumeMap[e.ItemId()]); err != nil {
			return
		}

		if err = s.itemEntityRepository.Save(ctx, e); err != nil {
			return
		}

		totalExp += e.Exp()
	})
	if err != nil {
		return 0, err
	}

	return totalExp, nil
}

func (s *MonsterEnhanceService) validateItemType(itemIds []int64) error {
	var err error

	lo.ForEach(itemIds, func(id int64, _ int) {
		var m *masterdata.Item
		m, err = masterdata.Master().Item.GetSafely(id)
		if err != nil {
			return
		}

		if m.ItemType != masterdata.Enhance {
			err = errors.New("itemのタイプが強化ではありません")
			return
		}
	})

	return err
}
