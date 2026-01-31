package prime

import (
	"bustub-go/types"
	"bustub-go/utils"
	"fmt"
	"math"
	"sync"
)

const BITSET_CAPACITY = 64

const CONSTANT = 0.79402

type HyperLogLog[T any] struct {
	cardinality  uint64
	bucketsMutex sync.Mutex
	buckets      []uint64
	bucketsSize  uint64
	nBits        int16
}

func NewHyperLogLog[T any](nBits int16) *HyperLogLog[T] {
	if nBits < 0 || nBits >= BITSET_CAPACITY {
		return nil
	}
	bucketsSize := uint64(1) << nBits
	return &HyperLogLog[T]{
		nBits:       nBits,
		bucketsSize: bucketsSize,
		buckets:     make([]uint64, bucketsSize),
	}
}

func (h *HyperLogLog[T]) GetCardinality() uint64 {
	return h.cardinality
}

func (h *HyperLogLog[T]) CalculateHash(val T) (utils.HashT, error) {
	var valObj *types.Value
	switch v := any(val).(type) {
	// TODO implement more type
	case int64:
		valObj = types.NewBigintValue(v)
	case string:
		valObj = types.NewVarcharValue(v)
	}
	if valObj == nil {
		return 0, fmt.Errorf("failed to create types.Value for type %T, valObj is nil", val)
	}
	return utils.HashValue(valObj)
}

func (h *HyperLogLog[T]) ComputeBinary(hash utils.HashT) uint64 {
	return uint64(hash)
}

func (h *HyperLogLog[T]) GetBucketIndex(binary uint64) uint64 {
	if h.nBits <= 0 {
		return 0
	}
	bucketIndex := uint64(0)
	for i := 0; i < int(h.nBits); i++ {
		bitPos := BITSET_CAPACITY - 1 - i
		bit := (binary >> bitPos) & 1
		bucketIndex <<= 1
		bucketIndex |= bit
	}
	return bucketIndex
}

func (h *HyperLogLog[T]) PositionOfLeftmostOne(binary uint64) uint64 {
	start := BITSET_CAPACITY - int(h.nBits)
	for i := start; i > 0; i-- {
		bit := (binary >> (i - 1)) & 1
		if bit == 1 {
			return uint64(BITSET_CAPACITY - i - int(h.nBits))
		}
	}
	return uint64(BITSET_CAPACITY - int(h.nBits))
}

func (h *HyperLogLog[T]) UpdateBucket(bucketIndex uint64, pValue uint64) {
	if bucketIndex >= h.bucketsSize {
		return
	}
	h.bucketsMutex.Lock()
	defer h.bucketsMutex.Unlock()
	if pValue > h.buckets[bucketIndex] {
		h.buckets[bucketIndex] = pValue
	}
}

func (h *HyperLogLog[T]) AddElem(val T) {
	hash, err := h.CalculateHash(val)
	if err != nil {
		panic(err)
	}
	binary := h.ComputeBinary(hash)
	bucketIndex := h.GetBucketIndex(binary)
	pValue := h.PositionOfLeftmostOne(binary) + 1
	h.UpdateBucket(bucketIndex, pValue)
}

func (h *HyperLogLog[T]) ComputeCardinality() {
	if h.nBits < 0 {
		return
	}
	h.bucketsMutex.Lock()
	defer h.bucketsMutex.Unlock()
	sum := 0.0
	for _, r := range h.buckets {
		sum += 1.0 / math.Pow(2, float64(r))
	}
	cardinality := CONSTANT * math.Pow(float64(h.bucketsSize), 2) / sum
	h.cardinality = uint64(cardinality)
}
