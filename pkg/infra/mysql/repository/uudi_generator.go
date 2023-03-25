package repository

import "github.com/google/uuid"

type UUIDGenerator struct {
}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (u *UUIDGenerator) GetUUID() int64 {
	generator := uuid.Must(uuid.NewRandom())
	return int64(generator.ID())
}
