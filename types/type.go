package types

import (
	"encoding/binary"
	"strconv"
)

type TypeId int

// TODO implement more type
const (
	BIGINT = iota
	VARCHAR
)

type Type interface {
	GetData(val *Value) []byte
	GetStorageSize(val *Value) uint32
}

func GetTypeInstance(typeId TypeId) Type {
	instance, ok := typeInstances[typeId]
	if !ok {
		panic("unsupported type id: " + strconv.Itoa(int(typeId)))
	}
	return instance
}

var typeInstances = map[TypeId]Type{
	BIGINT:  &IntType{typeId: BIGINT},
	VARCHAR: &VarcharType{},
}

type IntType struct {
	typeId TypeId
}

func (t *IntType) GetData(value *Value) []byte {
	switch t.typeId {
	// TODO implement more type e.g SMALLINT TINYINT ...
	case BIGINT:
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, uint64(value.val.bigint))
		return buf
	default:
		return nil
	}
}

func (t *IntType) GetStorageSize(value *Value) uint32 {
	return value.size.length
}

type VarcharType struct{}

func (t *VarcharType) GetData(value *Value) []byte {
	return value.val.varLen
}

func (t *VarcharType) GetStorageSize(value *Value) uint32 {
	return value.size.length
}
