package driver

import (
	"io"

	"github.com/freeconf/yang/meta"
	"github.com/ugorji/go/codec"
)

type Encoder struct {
	pack *codec.Encoder
}

func (e *Encoder) Encode(m *meta.Module, out io.Writer) error {
	var hnd codec.CborHandle
	e.pack = codec.NewEncoder(out, &hnd)
	return e.pack.Encode(ref{m})
}

type ref struct {
	obj interface{}
}

func (r *ref) CodecEncodeSelf(e *codec.Encoder) {
	switch x := r.obj.(type) {
	case *meta.Module:
		chkerr(e.Encode(x.Ident()))
		chkerr(e.Encode(x.Description()))
	}
}

func (r *ref) CodecDecodeSelf(e *codec.Decoder) {
	panic("encoder only")
}

func chkerr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type Decoder struct {
	pack *codec.Decoder
}

func (d *Decoder) Decode(in io.Reader) (*meta.Module, error) {
	var hnd codec.CborHandle
	d.pack = codec.NewDecoder(in, &hnd)
	ref := &ref2{obj: &meta.Module{}, b: new(meta.Builder)}
	err := d.pack.Decode(ref)
	return ref.obj.(*meta.Module), err
}

type ref2 struct {
	b   *meta.Builder
	obj interface{}
}

func (r *ref2) CodecEncodeSelf(e *codec.Encoder) {
	panic("decoder only")
}

func (r *ref2) CodecDecodeSelf(d *codec.Decoder) {
	switch x := r.obj.(type) {
	case *meta.Module:
		var s string
		chkerr(d.Decode(&s))
		x = r.b.Module(s, nil)
		r.obj = x
		chkerr(d.Decode(&s))
		r.b.Description(x, s)
	}
}
