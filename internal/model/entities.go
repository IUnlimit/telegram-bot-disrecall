package model

type Static struct {
	// Rows 记录数
	Rows int64
	// MB 占用的储存空间
	MB float64
}

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func NewPair[K comparable, V any](key K, value V) *Pair[K, V] {
	return &Pair[K, V]{Key: key, Value: value}
}
