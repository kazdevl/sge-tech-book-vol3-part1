package entity

import (
	"errors"
	"game-server-example/pkg/infra/mysql/datamodel"
)

type CoinEntity struct {
	userId int64
	num    int64
}

func NewCoinEntity(m *datamodel.UserCoin) *CoinEntity {
	return &CoinEntity{userId: m.UserID, num: m.Num}
}

func (e *CoinEntity) Consume(count int64) error {
	if e.num < count {
		return errors.New("消費数が所持数を上まっています")
	}
	e.num -= count
	return nil
}

func (e *CoinEntity) Add(num int64) {
	e.num += num
}

func (e *CoinEntity) UserId() int64 {
	return e.userId
}

func (e *CoinEntity) Num() int64 {
	return e.num
}
