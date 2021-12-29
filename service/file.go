package service

import (
	"context"
	"io"
)

// type InteruptReader struct {
// 	r        io.Reader
// 	interupt chan int
// }

// func NewInteruptReader(r io.Reader) InteruptReader {
// 	return InteruptReader{
// 		r,
// 		make(chan int),
// 	}
// }

// func (r InteruptReader) Read(p []byte) (n int, err error) {
// 	if r.r == nil {
// 		return 0, io.EOF
// 	}
// 	select {
// 	case <-r.interupt:
// 		return r.r.Read(p)
// 	default:
// 		r.r = nil
// 		return 0, io.EOF
// 	}
// }

// func (r InteruptReader) Cancel() {
// 	r.interupt <- 0
// }

type reader struct {
	ctx context.Context
	r   io.Reader
}

// NewReader wraps an io.Reader to handle context cancellation.
//
// Context state is checked BEFORE every Read.
func NewReader(ctx context.Context, r io.Reader) io.Reader {
	if r, ok := r.(*reader); ok && ctx == r.ctx {
		return r
	}
	return &reader{ctx: ctx, r: r}
}

func (r *reader) Read(p []byte) (n int, err error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
		return r.r.Read(p)
	}
}

type writer struct {
	ctx context.Context
	w   io.Writer
}

type copier struct {
	writer
}

func NewWriter(ctx context.Context, w io.Writer) io.Writer {
	if w, ok := w.(*copier); ok && ctx == w.ctx {
		return w
	}
	return &copier{writer{ctx: ctx, w: w}}
}

// Write implements io.Writer, but with context awareness.
func (w *writer) Write(p []byte) (n int, err error) {
	select {
	case <-w.ctx.Done():
		return 0, w.ctx.Err()
	default:
		return w.w.Write(p)
	}
}
