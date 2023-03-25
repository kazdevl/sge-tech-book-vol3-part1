package service

import (
	"context"
	"game-server-example/pkg/entity"
	"game-server-example/pkg/ifrepository"
)

type UserGetDataService struct {
	monsterEntityRepository ifrepository.IFMonsterEntityRepository
	itemEntityRepository    ifrepository.IFItemEntityRepository
	coinEntityRepository    ifrepository.IFCoinEntityRepository
}

func NewUserGetDataService(monsterEntityRepository ifrepository.IFMonsterEntityRepository, itemEntityRepository ifrepository.IFItemEntityRepository, coinEntityRepository ifrepository.IFCoinEntityRepository) *UserGetDataService {
	return &UserGetDataService{monsterEntityRepository: monsterEntityRepository, itemEntityRepository: itemEntityRepository, coinEntityRepository: coinEntityRepository}
}

func (s *UserGetDataService) GetAllEntities(ctx context.Context, userId int64) (
	[]*entity.MonsterEntity,
	[]*entity.ItemEntity,
	*entity.CoinEntity,
	error,
) {
	monsterEntities, err := s.monsterEntityRepository.FindByUserId(ctx, userId)
	if err != nil {
		return nil, nil, nil, err
	}
	itemEntities, err := s.itemEntityRepository.FindByUserId(ctx, userId)
	if err != nil {
		return nil, nil, nil, err
	}
	coinEntity, err := s.coinEntityRepository.Get(ctx, userId)
	if err != nil {
		return nil, nil, nil, err
	}
	return monsterEntities, itemEntities, coinEntity, nil
}
