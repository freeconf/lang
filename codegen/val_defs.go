package codegen

type ValType int

const (
	ValTypeBinary ValType = iota + 1
	ValTypeBits
	ValTypeBool
	ValTypeDecimal64
	ValTypeEmpty
	ValTypeEnum
	ValTypeIdentityRef
	ValTypeInstanceRef
	ValTypeInt8
	ValTypeInt16
	ValTypeInt32
	ValTypeInt64
	ValTypeLeafRef
	ValTypeString
	ValTypeUInt8
	ValTypeUInt16
	ValTypeUInt32
	ValTypeUInt64
	ValTypeUnion
	ValTypeAny

	ValTypeBinaryList
	ValTypeBitsList
	ValTypeBoolList
	ValTypeDecimal64List
	ValTypeEmptyList
	ValTypeEnumList
	ValTypeIdentityRefList
	ValTypeInstanceRefList
	ValTypeInt8List
	ValTypeInt16List
	ValTypeInt32List
	ValTypeInt64List
	ValTypeLeafRefList
	ValTypeStringList
	ValTypeUInt8List
	ValTypeUInt16List
	ValTypeUInt32List
	ValTypeUInt64List
	ValTypeUnionList
	ValTypeAnyList
)

type Val struct {
	ValTypeId ValType
	Data      any
}
