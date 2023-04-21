package lang

import (
	"fmt"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/val"
)

func encodeVal(v val.Value) *pb.XVal {
	if v == nil {
		return nil
	}
	switch v.Format() {
	case val.FmtString:
		return &pb.XVal{Value: &pb.XVal_Str{Str: v.String()}}
	case val.FmtInt32:
		return &pb.XVal{Value: &pb.XVal_I32{I32: int32(v.Value().(int))}}
	case val.FmtInt64:
		return &pb.XVal{Value: &pb.XVal_I64{I64: v.Value().(int64)}}
	}
	panic(fmt.Sprintf("not implemented type %T", v))
}

func decodeVal(v *pb.XVal) val.Value {
	if v == nil {
		return nil
	}
	switch x := v.Value.(type) {
	case *pb.XVal_Str:
		return val.String(x.Str)
	case *pb.XVal_I32:
		return val.Int32(x.I32)
	case *pb.XVal_I64:
		return val.Int64(x.I64)
	}
	panic(fmt.Sprintf("not implemented type %s", v))
}
