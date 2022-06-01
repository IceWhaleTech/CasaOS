/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2021-12-20 14:15:46
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-05-30 18:49:46
 * @FilePath: /CasaOS/service/file.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package service

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
)

var FileQueue map[string]model.FileOperate

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
func FileOperate(list model.FileOperate) {
	for _, v := range list.Item {
		if list.Type == "move" {
			lastPath := v.From[strings.LastIndex(v.From, "/")+1:]
			if !file.CheckNotExist(list.To + "/" + lastPath) {
				continue
			}
			err := os.Rename(v.From, list.To+"/"+lastPath)
			if err != nil {
				continue
			}
		} else if list.Type == "copy" {
			err := file.CopyDir(v.From, list.To)
			if err != nil {
				continue
			}
		} else {
			continue
		}
	}
}

// file move or copy and send notify
func CheckFileStatus() {
	for {
		if len(FileQueue) == 0 {
			return
		}
		for k, v := range FileQueue {
			var total int64 = 0
			for i := 0; i < len(v.Item); i++ {
				if !v.Item[i].Finished {
					size, err := file.GetFileOrDirSize(v.To + "/" + filepath.Base(v.Item[i].From))
					if err != nil {
						continue
					}
					v.Item[i].ProcessedSize = size
					if size == v.Item[i].Size {
						v.Item[i].Finished = true
					}
					total += size
				} else {
					total += v.Item[i].ProcessedSize
				}

			}
			v.ProcessedSize = total
			FileQueue[k] = v
		}
		time.Sleep(time.Second * 3)
	}
}
