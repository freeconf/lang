package lang

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"

	"github.com/freeconf/lang/pb"
)

type FileSystemService struct {
	pb.UnimplementedFileSystemServer
	d *Driver
}

func (fs *FileSystemService) ReaderInit(ctx context.Context, req *pb.ReaderInitRequest) (*pb.ReaderInitResponse, error) {
	var rdr io.Reader
	var err error
	var stream *streamReader

	if req.GetStream() {
		stream = newStreamReader()
		rdr = stream
	} else if req.GetFname() != "" {
		rdr, err = os.Open(req.GetFname())
		if err != nil {
			return nil, err
		}
	} else if req.GetInlineContent() != nil {
		rdr = bytes.NewBuffer(req.GetInlineContent())
	} else {
		return nil, errors.New("invalid file reader init request")
	}
	streamHnd := fs.d.handles.Put(rdr)
	return &pb.ReaderInitResponse{StreamHnd: streamHnd}, nil
}

func (fs *FileSystemService) ReaderStream(srv pb.FileSystem_ReaderStreamServer) error {
	for {
		data, err := srv.Recv()
		if err != nil {
			return err
		}
		if data == nil {
			break
		}
		rdr := fs.d.handles.Require(data.StreamHnd).(*streamReader)
		rdr.chunks <- data.Chunk
	}
	return srv.SendAndClose(&pb.ReaderStreamResponse{})
}

func (fs *FileSystemService) WriterInit(ctx context.Context, req *pb.WriterInitRequest) (*pb.WriterInitResponse, error) {
	var wtr io.Writer
	var err error
	var stream *streamWriter
	if req.GetStream() {
		stream = newRemoteWriter()
		wtr = stream
	} else if req.GetFname() != "" {
		wtr, err = os.Create(req.GetFname())
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("invalid file writer init request")
	}
	streamHnd := fs.d.handles.Put(wtr)
	return &pb.WriterInitResponse{StreamHnd: streamHnd}, nil
}

func (fs *FileSystemService) WriterStream(req *pb.WriterStreamRequest, srv pb.FileSystem_WriterStreamServer) error {
	stream := fs.d.handles.Require(req.StreamHnd).(*streamWriter)
	for chunk := range stream.chunks {
		if err := srv.Send(&pb.WriterStreamData{Chunk: chunk}); err != nil {
			return err
		}
	}
	return nil
}

func (fs *FileSystemService) CloseStream(ctx context.Context, req *pb.CloseStreamRequest) (*pb.CloseStreamResponse, error) {
	var err error
	if obj := fs.d.handles.Get(req.StreamHnd); obj != nil {
		err = obj.(io.Closer).Close()
	}
	return &pb.CloseStreamResponse{}, err
}

type streamWriter struct {
	chunks chan []byte
}

func newRemoteWriter() *streamWriter {
	return &streamWriter{
		chunks: make(chan []byte),
	}
}

func (r *streamWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	r.chunks <- p
	return len(p), nil
}

func (r *streamWriter) Close() error {
	close(r.chunks)
	return nil
}

type streamReader struct {
	chunks   chan []byte
	remainer []byte
}

func newStreamReader() *streamReader {
	return &streamReader{
		chunks: make(chan []byte),
	}
}

func (r *streamReader) Read(p []byte) (int, error) {
	if len(r.remainer) == 0 {
		r.remainer = <-r.chunks
		if len(r.remainer) == 0 {
			return 0, io.EOF
		}
	}
	len := copy(p, r.remainer)
	r.remainer = r.remainer[len:]
	return len, nil
}

func (r *streamReader) Close() error {
	r.chunks <- nil
	return nil
}
