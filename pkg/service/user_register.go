package service

import (
	"context"
	"game-server-example/pkg/ifrepository"
	"game-server-example/pkg/masterdata"
)

type UserRegisterService struct {
	uuidGenerator           ifrepository.IFUUIDGenerator
	monsterEntityRepository ifrepository.IFMonsterEntityRepository
	itemEntityRepository    ifrepository.IFItemEntityRepository
	coinEntityRepository    ifrepository.IFCoinEntityRepository
}

func NewUserRegisterService(uuidGenerator ifrepository.IFUUIDGenerator, monsterEntityRepository ifrepository.IFMonsterEntityRepository, itemEntityRepository ifrepository.IFItemEntityRepository, coinEntityRepository ifrepository.IFCoinEntityRepository) *UserRegisterService {
	return &UserRegisterService{uuidGenerator: uuidGenerator, monsterEntityRepository: monsterEntityRepository, itemEntityRepository: itemEntityRepository, coinEntityRepository: coinEntityRepository}
}

func (s *UserRegisterService) Register(ctx context.Context) (int64, error) {
	userId := s.uuidGenerator.GetUUID()

	if err := s.addCoin(ctx, userId); err != nil {
		return 0, err
	}
	if err := s.addMonsters(ctx, userId); err != nil {
		return 0, err
	}
	if err := s.addItems(ctx, userId); err != nil {
		return 0, err
	}
	return userId, nil
}

func (s *UserRegisterService) addCoin(ctx context.Context, userId int64) error {
	e, err := s.coinEntityRepository.Create(ctx, userId)
	if err != nil {
		return err
	}
	e.Add(100000)
	return s.coinEntityRepository.Save(ctx, e)
}

func (s *UserRegisterService) addMonsters(ctx context.Context, userId int64) error {
	for _, monster := range masterdata.Master().Monster.All() {
		_, err := s.monsterEntityRepository.Create(ctx, userId, monster.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *UserRegisterService) addItems(ctx context.Context, userId int64) error {
	for _, item := range masterdata.Master().Item.All() {
		e, err := s.itemEntityRepository.Create(ctx, userId, item.Id)
		if err != nil {
			return err
		}

		e.Add(1000)
		if err = s.itemEntityRepository.Save(ctx, e); err != nil {
			return err
		}
	}
	return nil
}
