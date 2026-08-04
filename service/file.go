/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2021-12-20 14:15:46
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-07-04 16:18:23
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
	"sync"
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/moby/sys/mountinfo"
	"go.uber.org/zap"
)

var FileQueue sync.Map

var OpStrArr []string
var OpStrArrMu sync.Mutex

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

func FileOperate(k string) {
	list, ok := FileQueue.Load(k)
	if !ok {
		return
	}

	temp := list.(model.FileOperate)
	if temp.ProcessedSize > 0 {
		return
	}
	for i := 0; i < len(temp.Item); i++ {
		v := temp.Item[i]
		if temp.Type == "move" {
			lastPath := v.From[strings.LastIndex(v.From, "/")+1:]
			if !file.CheckNotExist(temp.To + "/" + lastPath) {
				if temp.Style == "skip" {
					temp.Item[i].Finished = true
					continue
				} else {
					os.RemoveAll(temp.To + "/" + lastPath)
				}
			}
			err := file.CopyDir(v.From, temp.To, temp.Style)
			if err == nil {
				err = os.RemoveAll(v.From)
				if err != nil {
					logger.Error("file move error", zap.Any("err", err))
					err = file.MoveFile(v.From, temp.To+"/"+lastPath)
					if err != nil {
						logger.Error("MoveFile error", zap.Any("err", err))
						continue
					}

				}
			}

		} else if temp.Type == "copy" {
			err := file.CopyDir(v.From, temp.To, temp.Style)
			if err != nil {
				continue
			}
		} else {
			continue
		}

	}
	temp.Finished = true
	FileQueue.Store(k, temp)
}

func ExecOpFile() {
	OpStrArrMu.Lock()
	length := len(OpStrArr)
	if length == 0 {
		OpStrArrMu.Unlock()
		return
	}
	if length > 1 {
		length = 1
	}
	items := make([]string, length)
	copy(items, OpStrArr[:length])
	OpStrArrMu.Unlock()
	for i := 0; i < length; i++ {
		go FileOperate(items[i])
	}
}

// file move or copy and send notify
func CheckFileStatus() {
	for {
		OpStrArrMu.Lock()
		length := len(OpStrArr)
		if length == 0 {
			OpStrArrMu.Unlock()
			return
		}
		currentItems := make([]string, length)
		copy(currentItems, OpStrArr)
		OpStrArrMu.Unlock()
		for _, v := range currentItems {
			var total int64 = 0
			item, ok := FileQueue.Load(v)
			if !ok {
				continue
			}
			temp := item.(model.FileOperate)
			for i := 0; i < len(temp.Item); i++ {
				if !temp.Item[i].Finished {
					size, err := file.GetFileOrDirSize(temp.To + "/" + filepath.Base(temp.Item[i].From))
					if err != nil {
						continue
					}
					temp.Item[i].ProcessedSize = size
					if size == temp.Item[i].Size {
						temp.Item[i].Finished = true
					}
					total += size
				} else {
					total += temp.Item[i].ProcessedSize
				}
			}
			temp.ProcessedSize = total
			FileQueue.Store(v, temp)
		}
		time.Sleep(time.Second * 3)
	}
}
func IsMounted(path string) bool {
	mounted, _ := mountinfo.Mounted(path)
	if mounted {
		return true
	}
	connections := MyService.Connections().GetConnectionsList()
	for _, v := range connections {
		if v.MountPoint == path {
			return true
		}
	}
	return false
}
