package file

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ArchiveWriter interface {
	Create(w io.Writer) error
	Write(filename string, file io.Reader, fileSize int64, modTime int64) error
	Close() error
}

type zipArchiveWriter struct{ w *zip.Writer }

func (z *zipArchiveWriter) Create(w io.Writer) error { z.w = zip.NewWriter(w); return nil }
func (z *zipArchiveWriter) Write(name string, r io.Reader, sz int64, _ int64) error {
	f, err := z.w.Create(name)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}
func (z *zipArchiveWriter) Close() error { return z.w.Close() }

type tarArchiveWriter struct{ w *tar.Writer }

func (t *tarArchiveWriter) Create(w io.Writer) error { t.w = tar.NewWriter(w); return nil }
func (t *tarArchiveWriter) Write(name string, r io.Reader, sz int64, _ int64) error {
	if err := t.w.WriteHeader(&tar.Header{Name: name, Size: sz, Mode: 0644}); err != nil {
		return err
	}
	_, err := io.Copy(t.w, r)
	return err
}
func (t *tarArchiveWriter) Close() error { return t.w.Close() }

type tarGzArchiveWriter struct {
	tw tarArchiveWriter
	gw *gzip.Writer
}

func (tg *tarGzArchiveWriter) Create(w io.Writer) error {
	tg.gw = gzip.NewWriter(w)
	return tg.tw.Create(tg.gw)
}
func (tg *tarGzArchiveWriter) Write(n string, r io.Reader, s int64, m int64) error {
	return tg.tw.Write(n, r, s, m)
}
func (tg *tarGzArchiveWriter) Close() error {
	if err := tg.tw.Close(); err != nil {
		return err
	}
	return tg.gw.Close()
}

func getCompressionWriter(t string) (string, ArchiveWriter, error) {
	switch t {
	case "zip", "":
		return ".zip", &zipArchiveWriter{}, nil
	case "tar":
		return ".tar", &tarArchiveWriter{}, nil
	case "targz":
		return ".tar.gz", &tarGzArchiveWriter{}, nil
	default:
		return "", nil, errors.New("format not implemented")
	}
}

func AddFileToArchive(ar ArchiveWriter, path, commonPath string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() && !info.Mode().IsRegular() {
		return nil
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if path != commonPath {
		filename := strings.TrimPrefix(path, commonPath)
		filename = strings.TrimPrefix(filename, string(filepath.Separator))
		return ar.Write(filename, f, info.Size(), info.ModTime().Unix())
	}
	return nil
}
