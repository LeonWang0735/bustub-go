package utils

import (
	"bustub-go/types"
	"unsafe"
)

type HashT uint64

type HashUtil struct{}

func HashBytes(bytes []byte, length int) HashT {
	hash := HashT(length)
	for i := 0; i < length; i++ {
		hash = ((hash << 5) ^ (hash >> 27)) ^ HashT(int8(bytes[i]))
	}
	return hash
}

func Hash[T any](val T) HashT {
	ptrAddr := unsafe.Pointer(&val)
	ptrSize := unsafe.Sizeof(val)
	return HashBytes(unsafe.Slice((*byte)(ptrAddr), ptrSize), int(ptrSize))
}

func HashValue(val *types.Value) (HashT, error) {
	switch val.GetTypeId() {
	// TODO implement more type
	case types.BIGINT:
		raw, err := types.GetAs[int64](val)
		if err != nil {
			return 0, err
		}
		return Hash(raw), nil
	case types.VARCHAR:
		raw := val.GetData()
		length := val.GetStorageSize()
		return HashBytes(raw, int(length)), nil
	default:
		panic("Unsupported types")
	}
}
