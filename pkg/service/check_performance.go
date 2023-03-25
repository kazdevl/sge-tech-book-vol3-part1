package service

import (
	"context"
	"game-server-example/pkg/ifrepository"
)

type CheckPerformanceService struct {
	monsterEntityRepository ifrepository.IFMonsterEntityRepository
	itemEntityRepository    ifrepository.IFItemEntityRepository
}

func NewCheckPerformanceService(monsterEntityRepository ifrepository.IFMonsterEntityRepository, itemEntityRepository ifrepository.IFItemEntityRepository) *CheckPerformanceService {
	return &CheckPerformanceService{monsterEntityRepository: monsterEntityRepository, itemEntityRepository: itemEntityRepository}
}

func (s *CheckPerformanceService) CheckPerformance(ctx context.Context, userId, monsterId int64, itemIds []int64) error {
	monster, err := s.monsterEntityRepository.Get(ctx, userId, monsterId)
	if err != nil {
		return err
	}
	for i := 0; i < 10; i++ {
		monster.AddExp(1)
		if err = s.monsterEntityRepository.Save(ctx, monster); err != nil {
			return err
		}
	}

	items, err := s.itemEntityRepository.FindByItemIds(ctx, userId, itemIds)
	if err != nil {
		return err
	}

	for _, item := range items {
		item.Consume(1)
		if err = s.itemEntityRepository.Save(ctx, item); err != nil {
			return err
		}
	}

	return nil
}
