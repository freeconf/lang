package lang

import (
	"fmt"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/val"
)

func encodeVal(v val.Value) *pb.Val {
	if v == nil {
		return nil
	}
	switch v.Format() {
	case val.FmtString:
		return &pb.Val{Value: &pb.Val_Str{Str: v.String()}}
	case val.FmtInt32:
		return &pb.Val{Value: &pb.Val_I32{I32: int32(v.Value().(int))}}
	case val.FmtInt64:
		return &pb.Val{Value: &pb.Val_I64{I64: v.Value().(int64)}}
	}
	panic(fmt.Sprintf("not implemented type %T", v))
}

func decodeVal(v *pb.Val) val.Value {
	if v == nil {
		return nil
	}
	switch x := v.Value.(type) {
	case *pb.Val_Str:
		return val.String(x.Str)
	case *pb.Val_I32:
		return val.Int32(x.I32)
	case *pb.Val_I64:
		return val.Int64(x.I64)
	}
	panic(fmt.Sprintf("not implemented type %s", v))
}
