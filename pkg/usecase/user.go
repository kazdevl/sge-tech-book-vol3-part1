package usecase

import (
	"context"
	"game-server-example/pkg/entity"
	"game-server-example/pkg/service"
)

type UserUsecase struct {
	userRegisterService *service.UserRegisterService
	userGetDataService  *service.UserGetDataService
}

func NewUserUsecase(userRegisterService *service.UserRegisterService, userGetDataService *service.UserGetDataService) *UserUsecase {
	return &UserUsecase{userRegisterService: userRegisterService, userGetDataService: userGetDataService}
}

func (u *UserUsecase) RegisterData(ctx context.Context) (int64, error) {
	userId, err := u.userRegisterService.Register(ctx)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (u *UserUsecase) GetData(ctx context.Context, userId int64) (
	[]*entity.MonsterEntity,
	[]*entity.ItemEntity,
	*entity.CoinEntity,
	error,
) {
	return u.userGetDataService.GetAllEntities(ctx, userId)
}
