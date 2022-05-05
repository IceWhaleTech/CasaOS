package file

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/disintegration/imaging"
	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
)

func GetImage(path string, width, height int) ([]byte, error) {
	if thumbnail, err := GetThumbnailByOwnerPhotos(path); err == nil {
		return thumbnail, nil
	} else {
		return GetThumbnailByWebPhoto(path, width, height)
	}
}
func GetThumbnailByOwnerPhotos(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	buff := &bytes.Buffer{}

	defer file.Close()
	offset := 0
	offsets := []int{12, 30}

	head := make([]byte, 0xffff)

	r := io.TeeReader(file, buff)
	_, err = r.Read(head)
	if err != nil {
		return nil, err
	}

	for _, offset = range offsets {
		if _, err = exif.ParseExifHeader(head[offset:]); err == nil {
			break
		}
	}

	im, err := exifcommon.NewIfdMappingWithStandard()
	if err != nil {
		return nil, err
	}

	_, index, err := exif.Collect(im, exif.NewTagIndex(), head[offset:])
	if err != nil {
		return nil, err
	}

	ifd := index.RootIfd.NextIfd()
	if err != nil {
		return nil, err
	}
	thumbnail, err := ifd.Thumbnail()
	if err != nil {
		return nil, err
	}
	return thumbnail, nil
}
func GetThumbnailByWebPhoto(path string, width, height int) ([]byte, error) {
	src, err := imaging.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	img := imaging.Resize(src, width, 0, imaging.Lanczos)

	f, err := imaging.FormatFromFilename(path)
	if err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}
	imaging.Encode(&buf, img, f)
	return buf.Bytes(), nil
}
