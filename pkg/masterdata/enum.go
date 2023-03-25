package masterdata

type RarityType int64

const (
	N RarityType = iota
	R
	SR
	UR
)

type ItemType int64

const (
	Enhance ItemType = iota
	Collection
	Exchange
)

type ElementType int64

const (
	Fire ElementType = iota
	Water
	Nature
	Thunder
	Rock
	Holy
	Dark
)
