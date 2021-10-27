package file

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

var (
	buffSize = 1 << 20
)

// ReadLineFromEnd --
type ReadLineFromEnd struct {
	f *os.File

	fileSize int
	bwr      *bytes.Buffer
	lineBuff []byte
	swapBuff []byte

	isFirst bool
}

// NewReadLineFromEnd --
func NewReadLineFromEnd(name string) (rd *ReadLineFromEnd, err error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if info.IsDir() {
		return nil, fmt.Errorf("not file")
	}
	fileSize := int(info.Size())
	rd = &ReadLineFromEnd{
		f:        f,
		fileSize: fileSize,
		bwr:      bytes.NewBuffer([]byte{}),
		lineBuff: make([]byte, 0),
		swapBuff: make([]byte, buffSize),
		isFirst:  true,
	}
	return rd, nil
}

// ReadLine 结尾包含'\n'
func (c *ReadLineFromEnd) ReadLine() (line []byte, err error) {
	var ok bool
	for {
		ok, err = c.buff()
		if err != nil {
			return nil, err
		}
		if ok {
			break
		}
	}
	line, err = c.bwr.ReadBytes('\n')
	if err == io.EOF && c.fileSize > 0 {
		err = nil
	}
	return line, err
}

// Close --
func (c *ReadLineFromEnd) Close() (err error) {
	return c.f.Close()
}

func (c *ReadLineFromEnd) buff() (ok bool, err error) {
	if c.fileSize == 0 {
		return true, nil
	}

	if c.bwr.Len() >= buffSize {
		return true, nil
	}

	offset := 0
	if c.fileSize > buffSize {
		offset = c.fileSize - buffSize
	}
	_, err = c.f.Seek(int64(offset), 0)
	if err != nil {
		return false, err
	}

	n, err := c.f.Read(c.swapBuff)
	if err != nil && err != io.EOF {
		return false, err
	}
	if c.fileSize < n {
		n = c.fileSize
	}
	if n == 0 {
		return true, nil
	}

	for {
		m := bytes.LastIndex(c.swapBuff[:n], []byte{'\n'})
		if m == -1 {
			break
		}
		if m < n-1 {
			err = c.writeLine(c.swapBuff[m+1 : n])
			if err != nil {
				return false, err
			}
			ok = true
		} else if m == n-1 && !c.isFirst {
			err = c.writeLine(nil)
			if err != nil {
				return false, err
			}
			ok = true
		}
		n = m
		if n == 0 {
			break
		}
	}
	if n > 0 {
		reverseBytes(c.swapBuff[:n])
		c.lineBuff = append(c.lineBuff, c.swapBuff[:n]...)
	}
	if offset == 0 {
		err = c.writeLine(nil)
		if err != nil {
			return false, err
		}
		ok = true
	}
	c.fileSize = offset
	if c.isFirst {
		c.isFirst = false
	}
	return ok, nil
}

func (c *ReadLineFromEnd) writeLine(b []byte) (err error) {
	if len(b) > 0 {
		_, err = c.bwr.Write(b)
		if err != nil {
			return err
		}
	}
	if len(c.lineBuff) > 0 {
		reverseBytes(c.lineBuff)
		_, err = c.bwr.Write(c.lineBuff)
		if err != nil {
			return err
		}
		c.lineBuff = c.lineBuff[:0]
	}
	_, err = c.bwr.Write([]byte{'\n'})
	if err != nil {
		return err
	}
	return nil
}

func reverseBytes(b []byte) {
	n := len(b)
	if n <= 1 {
		return
	}
	for i := 0; i < n; i++ {
		k := n - 1
		if k != i {
			b[i], b[k] = b[k], b[i]
		}
		n--
	}
}
