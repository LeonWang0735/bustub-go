package types

import (
	"fmt"
)

type Value struct {
	typeId TypeId
	val    Val
	size   Size
}

func NewBigintValue(v int64) *Value {
	return &Value{
		typeId: BIGINT,
		val: Val{
			bigint: v,
		},
		size: Size{
			length: 0,
		},
	}
}

func NewVarcharValue(v string) *Value {
	return &Value{
		typeId: VARCHAR,
		val: Val{
			varLen: []byte(v),
		},
		size: Size{
			length: uint32(len(v)),
		},
	}
}

func (v *Value) GetTypeId() TypeId {
	return v.typeId
}

func (v *Value) GetStorageSize() uint32 {
	return GetTypeInstance(v.typeId).GetStorageSize(v)
}

func (v *Value) GetData() []byte {
	return GetTypeInstance(v.typeId).GetData(v)
}

func GetAs[T any](v *Value) (T, error) {
	if v == nil {
		var zero T
		return zero, fmt.Errorf("value is nil")
	}

	var rawValue interface{}
	switch v.typeId {
	// TODO implement more type
	case BIGINT:
		rawValue = v.val.bigint
	case VARCHAR:
		if len(v.val.varLen) > 0 {
			rawValue = v.val.varLen
		} else {
			rawValue = v.val.constVarLen_
		}
	default:
		var zero T
		return zero, fmt.Errorf("unsupported type id: %d", v.typeId)
	}
	val, ok := rawValue.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf(
			"type mismatch: typeId=%d (raw type=%T), expected %T",
			v.typeId, rawValue, zero,
		)
	}
	return val, nil
}

type Val struct {
	boolean      int8
	tinyint      int8
	smallint     int16
	integer      int32
	bigint       int64
	decimal      float64
	timestamp    uint64
	varLen       []byte
	constVarLen_ string
}

type Size struct {
	length     uint32
	elemTypeId TypeId
}
