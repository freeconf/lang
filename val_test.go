package lang

import (
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/val"
)

func TestVal(t *testing.T) {
	v := val.Int32(99)
	pv := encodeVal(v)
	rt := decodeVal(pv)
	fc.AssertEqual(t, v.Format(), rt.Format())
	fc.AssertEqual(t, v.Value(), rt.Value())
}

func TestValList(t *testing.T) {
	v := val.Int32List([]int{99, 100})
	pv := encodeVal(v)
	rt := decodeVal(pv)
	fc.AssertEqual(t, v.Format(), rt.Format())
	fc.AssertEqual(t, v.Value(), rt.Value())
}
