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

type zipArchiveWriter struct {
	zw *zip.Writer
}

func (z *zipArchiveWriter) Create(w io.Writer) error {
	z.zw = zip.NewWriter(w)
	return nil
}

func (z *zipArchiveWriter) Write(filename string, file io.Reader, fileSize int64, modTime int64) error {
	f, err := z.zw.Create(filename)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, file)
	return err
}

func (z *zipArchiveWriter) Close() error {
	return z.zw.Close()
}

type tarArchiveWriter struct {
	tw *tar.Writer
}

func (t *tarArchiveWriter) Create(w io.Writer) error {
	t.tw = tar.NewWriter(w)
	return nil
}

func (t *tarArchiveWriter) Write(filename string, file io.Reader, fileSize int64, modTime int64) error {
	header := &tar.Header{
		Name: filename,
		Size: fileSize,
		Mode: 0644,
	}
	if err := t.tw.WriteHeader(header); err != nil {
		return err
	}
	_, err := io.Copy(t.tw, file)
	return err
}

func (t *tarArchiveWriter) Close() error {
	return t.tw.Close()
}

type tarGzArchiveWriter struct {
	tw *tarArchiveWriter
	gw *gzip.Writer
}

func (tg *tarGzArchiveWriter) Create(w io.Writer) error {
	tg.gw = gzip.NewWriter(w)
	tg.tw = &tarArchiveWriter{}
	return tg.tw.Create(tg.gw)
}

func (tg *tarGzArchiveWriter) Write(filename string, file io.Reader, fileSize int64, modTime int64) error {
	return tg.tw.Write(filename, file, fileSize, modTime)
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
