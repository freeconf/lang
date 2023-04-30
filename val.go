package lang

// This file is generated from val.go.in

import (
	"fmt"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/val"
)

func encodeVal(v val.Value) *pb.Val {
	if v == nil {
		return nil
	}
	f := pb.Format(v.Format())
	switch v.Format() {
	case val.FmtBinary:
		xVal := v.Value().([]byte)
		x := xVal
		return &pb.Val{Format: f, Value: &pb.ValUnion{Value: &pb.ValUnion_BinaryVal{BinaryVal: x}}}
	case val.FmtBits:
		xVal := v.Value().([]byte)
		x := xVal
		return &pb.Val{Format: f, Value: &pb.ValUnion{Value: &pb.ValUnion_BitsVal{BitsVal: x}}}
	case val.FmtBool:
		xVal := v.Value().(bool)
		x := xVal
		return &pb.Val{Format: f, Value: &pb.ValUnion{Value: &pb.ValUnion_BoolVal{BoolVal: x}}}
	case val.FmtDecimal64:
		xVal := v.Value().(float64)
		x := xVal
		return &pb.Val{Format: f, Value: &pb.ValUnion{Value: &pb.ValUnion_Decimal64Val{Decimal64Val: x}}}
	case val.FmtEmpty:
		xVal := v.Value().(val.Value)
		x := newPbEmpty(xVal)
		return &pb.Val{Format: f, Value: &pb.ValUnion{Value: &pb.ValUnion_EmptyVal{EmptyVal: x}}}
	case val.FmtEnum:
		xVal := v.Value().(val.Enum)
		x := newPbEnum(xVal)
		return &pb.Val{Format: f, Value: &pb.ValUnion{Value: &pb.ValUnion_EnumVal{EnumVal: x}}}
	case val.FmtIdentityRef:
		xVal := v.Value().(val.IdentRef)
		x := newPbIdentRef(xVal)
		return &pb.Val{Format: f, Value: &pb.ValUnion{Value: &pb.ValUnion_IdentRefVal{IdentRefVal: x}}}
	case val.FmtInt8:
		xVal := v.Value().(int8)
		x := int32(xVal)
		return &pb.Val{Format: f, Value: &pb.ValUnion{Value: &pb.ValUnion_Int8Val{Int8Val: x}}}
	case val.FmtInt16:
		xVal := v.Value().(int16)
		x := int32(xVal)
		return &pb.Val{Format: f, Value: &pb.ValUnion{Value: &pb.ValUnion_Int16Val{Int16Val: x}}}
	case val.FmtInt32:
		xVal := v.Value().(int)
		x := int32(xVal)
		return &pb.Val{Format: f, Value: &pb.ValUnion{Value: &pb.ValUnion_Int32Val{Int32Val: x}}}
	case val.FmtInt64:
		xVal := v.Value().(int64)
		x := xVal
		return &pb.Val{Format: f, Value: &pb.ValUnion{Value: &pb.ValUnion_Int64Val{Int64Val: x}}}
	case val.FmtString:
		xVal := v.Value().(string)
		x := xVal
		return &pb.Val{Format: f, Value: &pb.ValUnion{Value: &pb.ValUnion_StringVal{StringVal: x}}}
	case val.FmtUInt8:
		xVal := v.Value().(uint8)
		x := uint32(xVal)
		return &pb.Val{Format: f, Value: &pb.ValUnion{Value: &pb.ValUnion_Uint8Val{Uint8Val: x}}}
	case val.FmtUInt16:
		xVal := v.Value().(uint16)
		x := uint32(xVal)
		return &pb.Val{Format: f, Value: &pb.ValUnion{Value: &pb.ValUnion_Uint16Val{Uint16Val: x}}}
	case val.FmtUInt32:
		xVal := v.Value().(uint)
		x := uint32(xVal)
		return &pb.Val{Format: f, Value: &pb.ValUnion{Value: &pb.ValUnion_Uint32Val{Uint32Val: x}}}
	case val.FmtUInt64:
		xVal := v.Value().(uint64)
		x := xVal
		return &pb.Val{Format: f, Value: &pb.ValUnion{Value: &pb.ValUnion_Uint64Val{Uint64Val: x}}}
	case val.FmtBinaryList:
		xVals := v.Value().([][]byte)
		vals := make([]*pb.ValUnion, len(xVals))
		for i, xVal := range xVals {
			x := xVal
			vals[i] = &pb.ValUnion{Value: &pb.ValUnion_BinaryVal{BinaryVal: x}}
		}
		return &pb.Val{Format: f, ListValue: vals}
	case val.FmtBitsList:
		xVals := v.Value().([][]byte)
		vals := make([]*pb.ValUnion, len(xVals))
		for i, xVal := range xVals {
			x := xVal
			vals[i] = &pb.ValUnion{Value: &pb.ValUnion_BitsVal{BitsVal: x}}
		}
		return &pb.Val{Format: f, ListValue: vals}
	case val.FmtBoolList:
		xVals := v.Value().([]bool)
		vals := make([]*pb.ValUnion, len(xVals))
		for i, xVal := range xVals {
			x := xVal
			vals[i] = &pb.ValUnion{Value: &pb.ValUnion_BoolVal{BoolVal: x}}
		}
		return &pb.Val{Format: f, ListValue: vals}
	case val.FmtDecimal64List:
		xVals := v.Value().([]float64)
		vals := make([]*pb.ValUnion, len(xVals))
		for i, xVal := range xVals {
			x := xVal
			vals[i] = &pb.ValUnion{Value: &pb.ValUnion_Decimal64Val{Decimal64Val: x}}
		}
		return &pb.Val{Format: f, ListValue: vals}
	case val.FmtEmptyList:
		xVals := v.Value().([]val.Value)
		vals := make([]*pb.ValUnion, len(xVals))
		for i, xVal := range xVals {
			x := newPbEmpty(xVal)
			vals[i] = &pb.ValUnion{Value: &pb.ValUnion_EmptyVal{EmptyVal: x}}
		}
		return &pb.Val{Format: f, ListValue: vals}
	case val.FmtEnumList:
		xVals := v.Value().([]val.Enum)
		vals := make([]*pb.ValUnion, len(xVals))
		for i, xVal := range xVals {
			x := newPbEnum(xVal)
			vals[i] = &pb.ValUnion{Value: &pb.ValUnion_EnumVal{EnumVal: x}}
		}
		return &pb.Val{Format: f, ListValue: vals}
	case val.FmtIdentityRefList:
		xVals := v.Value().([]val.IdentRef)
		vals := make([]*pb.ValUnion, len(xVals))
		for i, xVal := range xVals {
			x := newPbIdentRef(xVal)
			vals[i] = &pb.ValUnion{Value: &pb.ValUnion_IdentRefVal{IdentRefVal: x}}
		}
		return &pb.Val{Format: f, ListValue: vals}
	case val.FmtInt8List:
		xVals := v.Value().([]int8)
		vals := make([]*pb.ValUnion, len(xVals))
		for i, xVal := range xVals {
			x := int32(xVal)
			vals[i] = &pb.ValUnion{Value: &pb.ValUnion_Int8Val{Int8Val: x}}
		}
		return &pb.Val{Format: f, ListValue: vals}
	case val.FmtInt16List:
		xVals := v.Value().([]int16)
		vals := make([]*pb.ValUnion, len(xVals))
		for i, xVal := range xVals {
			x := int32(xVal)
			vals[i] = &pb.ValUnion{Value: &pb.ValUnion_Int16Val{Int16Val: x}}
		}
		return &pb.Val{Format: f, ListValue: vals}
	case val.FmtInt32List:
		xVals := v.Value().([]int)
		vals := make([]*pb.ValUnion, len(xVals))
		for i, xVal := range xVals {
			x := int32(xVal)
			vals[i] = &pb.ValUnion{Value: &pb.ValUnion_Int32Val{Int32Val: x}}
		}
		return &pb.Val{Format: f, ListValue: vals}
	case val.FmtInt64List:
		xVals := v.Value().([]int64)
		vals := make([]*pb.ValUnion, len(xVals))
		for i, xVal := range xVals {
			x := xVal
			vals[i] = &pb.ValUnion{Value: &pb.ValUnion_Int64Val{Int64Val: x}}
		}
		return &pb.Val{Format: f, ListValue: vals}
	case val.FmtStringList:
		xVals := v.Value().([]string)
		vals := make([]*pb.ValUnion, len(xVals))
		for i, xVal := range xVals {
			x := xVal
			vals[i] = &pb.ValUnion{Value: &pb.ValUnion_StringVal{StringVal: x}}
		}
		return &pb.Val{Format: f, ListValue: vals}
	case val.FmtUInt8List:
		xVals := v.Value().([]uint8)
		vals := make([]*pb.ValUnion, len(xVals))
		for i, xVal := range xVals {
			x := uint32(xVal)
			vals[i] = &pb.ValUnion{Value: &pb.ValUnion_Uint8Val{Uint8Val: x}}
		}
		return &pb.Val{Format: f, ListValue: vals}
	case val.FmtUInt16List:
		xVals := v.Value().([]uint16)
		vals := make([]*pb.ValUnion, len(xVals))
		for i, xVal := range xVals {
			x := uint32(xVal)
			vals[i] = &pb.ValUnion{Value: &pb.ValUnion_Uint16Val{Uint16Val: x}}
		}
		return &pb.Val{Format: f, ListValue: vals}
	case val.FmtUInt32List:
		xVals := v.Value().([]uint)
		vals := make([]*pb.ValUnion, len(xVals))
		for i, xVal := range xVals {
			x := uint32(xVal)
			vals[i] = &pb.ValUnion{Value: &pb.ValUnion_Uint32Val{Uint32Val: x}}
		}
		return &pb.Val{Format: f, ListValue: vals}
	case val.FmtUInt64List:
		xVals := v.Value().([]uint64)
		vals := make([]*pb.ValUnion, len(xVals))
		for i, xVal := range xVals {
			x := xVal
			vals[i] = &pb.ValUnion{Value: &pb.ValUnion_Uint64Val{Uint64Val: x}}
		}
		return &pb.Val{Format: f, ListValue: vals}
	}
	panic(fmt.Sprintf("not implemented type %T", v))
}

func decodeVal(v *pb.Val) val.Value {
	if v == nil {
		return nil
	}
	f := val.Format(v.Format)
	switch f {
	case val.FmtBinary:
		pval := v.Value.Value.(*pb.ValUnion_BinaryVal).BinaryVal
		return val.Binary(pval)
	case val.FmtBits:
		pval := v.Value.Value.(*pb.ValUnion_BitsVal).BitsVal
		return val.Bits(pval)
	case val.FmtBool:
		pval := v.Value.Value.(*pb.ValUnion_BoolVal).BoolVal
		return val.Bool(pval)
	case val.FmtDecimal64:
		pval := v.Value.Value.(*pb.ValUnion_Decimal64Val).Decimal64Val
		return val.Decimal64(pval)
	case val.FmtEmpty:
		pval := v.Value.Value.(*pb.ValUnion_EmptyVal).EmptyVal
		return newEmpty(pval)
	case val.FmtEnum:
		pval := v.Value.Value.(*pb.ValUnion_EnumVal).EnumVal
		return newEnum(pval)
	case val.FmtIdentityRef:
		pval := v.Value.Value.(*pb.ValUnion_IdentRefVal).IdentRefVal
		return newIdentRef(pval)
	case val.FmtInt8:
		pval := v.Value.Value.(*pb.ValUnion_Int8Val).Int8Val
		return val.Int8(int8(pval))
	case val.FmtInt16:
		pval := v.Value.Value.(*pb.ValUnion_Int16Val).Int16Val
		return val.Int16(int16(pval))
	case val.FmtInt32:
		pval := v.Value.Value.(*pb.ValUnion_Int32Val).Int32Val
		return val.Int32(int(pval))
	case val.FmtInt64:
		pval := v.Value.Value.(*pb.ValUnion_Int64Val).Int64Val
		return val.Int64(pval)
	case val.FmtString:
		pval := v.Value.Value.(*pb.ValUnion_StringVal).StringVal
		return val.String(pval)
	case val.FmtUInt8:
		pval := v.Value.Value.(*pb.ValUnion_Uint8Val).Uint8Val
		return val.UInt8(uint8(pval))
	case val.FmtUInt16:
		pval := v.Value.Value.(*pb.ValUnion_Uint16Val).Uint16Val
		return val.UInt16(uint16(pval))
	case val.FmtUInt32:
		pval := v.Value.Value.(*pb.ValUnion_Uint32Val).Uint32Val
		return val.UInt32(uint(pval))
	case val.FmtUInt64:
		pval := v.Value.Value.(*pb.ValUnion_Uint64Val).Uint64Val
		return val.UInt64(pval)
	case val.FmtBinaryList:
		xVals := make([][]byte, len(v.ListValue))
		for i, next := range v.ListValue {
			pval := next.Value.(*pb.ValUnion_BinaryVal).BinaryVal
			xVals[i] = pval
		}
		return val.BinaryList(xVals)
	case val.FmtBitsList:
		xVals := make([][]byte, len(v.ListValue))
		for i, next := range v.ListValue {
			pval := next.Value.(*pb.ValUnion_BitsVal).BitsVal
			xVals[i] = pval
		}
		return val.BitsList(xVals)
	case val.FmtBoolList:
		xVals := make([]bool, len(v.ListValue))
		for i, next := range v.ListValue {
			pval := next.Value.(*pb.ValUnion_BoolVal).BoolVal
			xVals[i] = pval
		}
		return val.BoolList(xVals)
	case val.FmtDecimal64List:
		xVals := make([]float64, len(v.ListValue))
		for i, next := range v.ListValue {
			pval := next.Value.(*pb.ValUnion_Decimal64Val).Decimal64Val
			xVals[i] = pval
		}
		return val.Decimal64List(xVals)
	case val.FmtEmptyList:
		xVals := make([]val.Value, len(v.ListValue))
		for i, next := range v.ListValue {
			pval := next.Value.(*pb.ValUnion_EmptyVal).EmptyVal
			xVals[i] = newEmpty(pval)
		}
		return newEmptyList(xVals)
	case val.FmtEnumList:
		xVals := make([]val.Enum, len(v.ListValue))
		for i, next := range v.ListValue {
			pval := next.Value.(*pb.ValUnion_EnumVal).EnumVal
			xVals[i] = newEnum(pval)
		}
		return val.EnumList(xVals)
	case val.FmtIdentityRefList:
		xVals := make([]val.IdentRef, len(v.ListValue))
		for i, next := range v.ListValue {
			pval := next.Value.(*pb.ValUnion_IdentRefVal).IdentRefVal
			xVals[i] = newIdentRef(pval)
		}
		return val.IdentRefList(xVals)
	case val.FmtInt8List:
		xVals := make([]int8, len(v.ListValue))
		for i, next := range v.ListValue {
			pval := next.Value.(*pb.ValUnion_Int8Val).Int8Val
			xVals[i] = int8(pval)
		}
		return val.Int8List(xVals)
	case val.FmtInt16List:
		xVals := make([]int16, len(v.ListValue))
		for i, next := range v.ListValue {
			pval := next.Value.(*pb.ValUnion_Int16Val).Int16Val
			xVals[i] = int16(pval)
		}
		return val.Int16List(xVals)
	case val.FmtInt32List:
		xVals := make([]int, len(v.ListValue))
		for i, next := range v.ListValue {
			pval := next.Value.(*pb.ValUnion_Int32Val).Int32Val
			xVals[i] = int(pval)
		}
		return val.Int32List(xVals)
	case val.FmtInt64List:
		xVals := make([]int64, len(v.ListValue))
		for i, next := range v.ListValue {
			pval := next.Value.(*pb.ValUnion_Int64Val).Int64Val
			xVals[i] = pval
		}
		return val.Int64List(xVals)
	case val.FmtStringList:
		xVals := make([]string, len(v.ListValue))
		for i, next := range v.ListValue {
			pval := next.Value.(*pb.ValUnion_StringVal).StringVal
			xVals[i] = pval
		}
		return val.StringList(xVals)
	case val.FmtUInt8List:
		xVals := make([]uint8, len(v.ListValue))
		for i, next := range v.ListValue {
			pval := next.Value.(*pb.ValUnion_Uint8Val).Uint8Val
			xVals[i] = uint8(pval)
		}
		return val.UInt8List(xVals)
	case val.FmtUInt16List:
		xVals := make([]uint16, len(v.ListValue))
		for i, next := range v.ListValue {
			pval := next.Value.(*pb.ValUnion_Uint16Val).Uint16Val
			xVals[i] = uint16(pval)
		}
		return val.UInt16List(xVals)
	case val.FmtUInt32List:
		xVals := make([]uint, len(v.ListValue))
		for i, next := range v.ListValue {
			pval := next.Value.(*pb.ValUnion_Uint32Val).Uint32Val
			xVals[i] = uint(pval)
		}
		return val.UInt32List(xVals)
	case val.FmtUInt64List:
		xVals := make([]uint64, len(v.ListValue))
		for i, next := range v.ListValue {
			pval := next.Value.(*pb.ValUnion_Uint64Val).Uint64Val
			xVals[i] = pval
		}
		return val.UInt64List(xVals)
	default:
		panic(fmt.Sprintf("proto decoder not implemented type %s", v))
	}
}

func newEmpty(_ any) val.Value {
	return val.NotEmpty
}

func newEmptyList(_ any) val.Value {
	return val.NotEmpty
}

func newIdentRef(ref *pb.IdentRef) val.IdentRef {
	return val.IdentRef{
		Base:  ref.Base,
		Label: ref.Label,
	}
}

func newEnum(v *pb.EnumVal) val.Enum {
	return val.Enum{
		Id:    int(v.Id),
		Label: v.Label,
	}
}

func newInt8(v int32) val.Int8 {
	return val.Int8(int8(v))
}

func newInt16(v int32) val.Int16 {
	return val.Int16(int16(v))
}

func newUInt8(v uint32) val.UInt8 {
	return val.UInt8(uint8(v))
}

func newUInt16(v uint32) val.UInt16 {
	return val.UInt16(uint16(v))
}

func newPbEnum(e val.Enum) *pb.EnumVal {
	return &pb.EnumVal{
		Id:    int32(e.Id),
		Label: e.Label,
	}
}

func newPbIdentRef(ref val.IdentRef) *pb.IdentRef {
	return &pb.IdentRef{
		Base:  ref.Base,
		Label: ref.Label,
	}
}

func newPbEmpty(v val.Value) string {
	return v.String()
}

func encodeVals(vals []val.Value) []*pb.Val {
	resp := make([]*pb.Val, len(vals))
	for i, v := range vals {
		resp[i] = encodeVal(v)
	}
	return resp
}

func decodeVals(vals []*pb.Val) []val.Value {
	resp := make([]val.Value, len(vals))
	for i, v := range vals {
		resp[i] = decodeVal(v)
	}
	return resp
}
