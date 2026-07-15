package file

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mholt/archiver/v3"
)

// GetSize get the file size
func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)
	return len(content), err
}

// GetExt get the file ext
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

// CheckNotExist check if the file exists
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

// CheckPermission check if the file has permission
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

// IsNotExistMkDir create a directory if it does not exist
func IsNotExistMkDir(src string) error {
	if CheckNotExist(src) {
		if err := MkDir(src); err != nil {
			return err
		}
	}
	return nil
}

// MkDir create a directory with safe permissions 0750
func MkDir(src string) error {
	err := os.MkdirAll(src, 0o750)
	if err != nil {
		return err
	}
	return nil
}

// RMDir remove a directory
func RMDir(src string) error {
	err := os.RemoveAll(src)
	if err != nil {
		return err
	}
	os.Remove(src)
	return nil
}

func RemoveAll(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return os.Remove(path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return os.Remove(dir)
}

// Open a file according to a specific mode
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// MustOpen maximize trying to open the file
func MustOpen(fileName, filePath string) (*os.File, error) {
	src := filePath
	if CheckPermission(src) {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err := IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}

// Exists checks if a file or directory exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// IsDir checks if the given path is a directory
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile checks if the given path is a file
func IsFile(path string) bool {
	return !IsDir(path)
}

func CreateFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func CreateFileAndWriteContent(path, content string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(content)
	write.Flush()
	return nil
}

// IsNotExistCreateFile create a file if it does not exist
func IsNotExistCreateFile(src string) error {
	if CheckNotExist(src) {
		if err := CreateFile(src); err != nil {
			return err
		}
	}
	return nil
}

func ReadFullFile(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		return []byte("")
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte("")
	}
	return content
}

// copyFileContents holds the logic shared by CopyFile and CopySingleFile
func copyFileContents(src, dst, style string) error {
	if Exists(dst) {
		if style == "skip" {
			return nil
		}
		os.Remove(dst)
	}

	srcfd, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcfd.Close()

	dstfd, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}

	srcinfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcinfo.Mode()&0o750)
}

// CopyFile copies a single file from src into the dst directory
func CopyFile(src, dst, style string) error {
	lastPath := src[strings.LastIndex(src, "/")+1:]

	if !strings.HasSuffix(dst, "/") {
		dst += "/"
	}
	dst += lastPath

	return copyFileContents(src, dst, style)
}

// CopySingleFile copies src to the exact dst path given
func CopySingleFile(src, dst, style string) error {
	return copyFileContents(src, dst, style)
}

// GetNoDuplicateFileName checks for duplicate file names and returns a unique name
func GetNoDuplicateFileName(fullPath string) string {
	dirPath, fileName := filepath.Split(fullPath)
	fileSuffix := path.Ext(fileName)
	filenameOnly := strings.TrimSuffix(fileName, fileSuffix)
	for i := 0; Exists(fullPath); i++ {
		fullPath = path.Join(dirPath, filenameOnly+"("+strconv.Itoa(i+1)+")"+fileSuffix)
	}
	return fullPath
}

// CopyDir copies a whole directory recursively
func CopyDir(src, dst, style string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	if !srcinfo.IsDir() {
		if err = CopyFile(src, dst, style); err != nil {
			fmt.Println(err)
		}
		return nil
	}
	lastPath := src[strings.LastIndex(src, "/")+1:]
	dst += "/" + lastPath

	if Exists(dst) {
		if style == "skip" {
			return nil
		}
		os.Remove(dst)
	}
	if err = os.MkdirAll(dst, srcinfo.Mode()&0o750); err != nil {
		return err
	}
	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := dst

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp, style); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp, style); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func WriteToPath(data []byte, path, name string) error {
	fullPath := path
	if strings.HasSuffix(path, "/") {
		fullPath += name
	} else {
		fullPath += "/" + name
	}
	return WriteToFullPath(data, fullPath, 0o600)
}

func WriteToFullPath(data []byte, fullPath string, perm fs.FileMode) error {
	if err := IsNotExistCreateFile(fullPath); err != nil {
		return err
	}

	safePerm := perm & 0o750

	file, err := os.OpenFile(fullPath,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		safePerm,
	)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)

	return err
}

// SpliceFiles concatenates multiple file chunks into a single file
func SpliceFiles(dir, path string, length, startPoint int) error {
	fullPath := path

	if err := IsNotExistCreateFile(fullPath); err != nil {
		return err
	}

	file, _ := os.OpenFile(fullPath,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0o600,
	)

	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)

	for i := 0; i < length+startPoint-1; i++ {
		data, err := ioutil.ReadFile(dir + "/" + strconv.Itoa(i+startPoint))
		if err != nil {
			return err
		}
		if _, err := bufferedWriter.Write(data); err != nil {
			return err
		}
	}

	bufferedWriter.Flush()

	return nil
}

func GetCompressionAlgorithm(t string) (string, archiver.Writer, error) {
	switch t {
	case "zip", "":
		return ".zip", archiver.NewZip(), nil
	case "tar":
		return ".tar", archiver.NewTar(), nil
	case "targz":
		return ".tar.gz", archiver.NewTarGz(), nil
	case "tarbz2":
		return ".tar.bz2", archiver.NewTarBz2(), nil
	case "tarxz":
		return ".tar.xz", archiver.NewTarXz(), nil
	case "tarlz4":
		return ".tar.lz4", archiver.NewTarLz4(), nil
	case "tarsz":
		return ".tar.sz", archiver.NewTarSz(), nil
	default:
		return "", nil, errors.New("format not implemented")
	}
}

func AddFile(ar archiver.Writer, path, commonPath string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !info.IsDir() && !info.Mode().IsRegular() {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if path != commonPath {
		filename := strings.TrimPrefix(path, commonPath)
		filename = strings.TrimPrefix(filename, string(filepath.Separator))
		err = ar.Write(archiver.File{
			FileInfo: archiver.FileInfo{
				FileInfo:   info,
				CustomName: filename,
			},
			ReadCloser: file,
		})
		if err != nil {
			return err
		}
	}

	if info.IsDir() {
		names, err := file.Readdirnames(0)
		if err != nil {
			return err
		}

		for _, name := range names {
			err = AddFile(ar, filepath.Join(path, name), commonPath)
			if err != nil {
				log.Printf("Failed to archive %v", err)
			}
		}
	}

	return nil
}

func CommonPrefix(sep byte, paths ...string) string {
	switch len(paths) {
	case 0:
		return ""
	case 1:
		return path.Clean(paths[0])
	}

	c := []byte(path.Clean(paths[0]))
	c = append(c, sep)

	for _, v := range paths[1:] {
		v = path.Clean(v) + string(sep)

		if len(v) < len(c) {
			c = c[:len(v)]
		}
		for i := 0; i < len(c); i++ {
			if v[i] != c[i] {
				c = c[:i]
				break
			}
		}
	}

	for i := len(c) - 1; i >= 0; i-- {
		if c[i] == sep {
			c = c[:i]
			break
		}
	}

	return string(c)
}

func GetFileOrDirSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	if fileInfo.IsDir() {
		return DirSizeB(path + "/")
	}
	return fileInfo.Size(), nil
}

// DirSizeB calculates total size of a directory in bytes
func DirSizeB(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

func ReadLine(lineNumber int, path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	fileScanner := bufio.NewScanner(file)
	lineCount := 1
	for fileScanner.Scan() {
		if lineCount == lineNumber {
			return fileScanner.Text()
		}
		lineCount++
	}
	defer file.Close()
	return ""
}

func NameAccumulation(name, dir string) string {
	path := filepath.Join(dir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return name
	}
	base := name
	strings.Split(base, "_")
	index := strings.LastIndex(base, "_")
	if index < 0 {
		index = len(base)
	}
	for i := 1; ; i++ {
		newPath := filepath.Join(dir, fmt.Sprintf("%s_%d", base[:index], i))
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return fmt.Sprintf("%s_%d", base[:index], i)
		}
	}
}

// parseHeaderFields extracts key-value pairs from header data
func parseHeaderFields(item []byte) (string, string, bool) {
	tarr := bytes.Split(item, []byte(";"))
	if len(tarr) != 2 {
		return "", "", false
	}

	tbyte := tarr[1]
	tbyte = bytes.ReplaceAll(tbyte, []byte("\r\n--"), []byte(""))
	tbyte = bytes.ReplaceAll(tbyte, []byte("name=\""), []byte(""))
	tempArr := bytes.Split(tbyte, []byte("\"\r\n\r\n"))
	if len(tempArr) != 2 {
		return "", "", false
	}
	return strings.TrimSpace(string(tempArr[0])), strings.TrimSpace(string(tempArr[1])), true
}

// ParseFileHeader parses file header from raw data
func ParseFileHeader(h, boundary []byte) (map[string]string, bool) {
	arr := bytes.Split(h, boundary)
	result := make(map[string]string)
	for _, item := range arr {
		key, value, ok := parseHeaderFields(item)
		if ok {
			result[key] = value
		}
	}
	return result, true
}

// readStreamChunk reads a chunk from the stream
func readStreamChunk(stream io.ReadCloser, buf []byte) (int, error) {
	readLen, err := stream.Read(buf)
	if err != nil && err != io.EOF {
		return 0, err
	}
	return readLen, nil
}

// locateBoundary finds the boundary in the data buffer
func locateBoundary(data, boundary []byte, readTotal int) int {
	return bytes.LastIndex(data[:readTotal], boundary)
}

// extractHeaderFromData extracts header map and remaining data
func extractHeaderFromData(data, boundary []byte, readTotal, boundaryLoc int) (map[string]string, []byte, bool, error) {
	startLoc := boundaryLoc + len(boundary)
	fileHeadLoc := bytes.Index(data[startLoc:readTotal], []byte("\r\n\r\n"))
	if fileHeadLoc == -1 {
		return nil, nil, false, nil
	}
	fileHeadLoc += startLoc

	headMap, ok := ParseFileHeader(data, boundary)
	if !ok {
		return headMap, nil, false, fmt.Errorf("ParseFileHeader fail: %s", string(data[startLoc:fileHeadLoc]))
	}
	return headMap, data[fileHeadLoc+4 : readTotal], true, nil
}

// ReadToBoundary reads data from stream until a boundary is found
func ReadToBoundary(boundary []byte, stream io.ReadCloser, target io.WriteCloser) ([]byte, bool, error) {
	readData := make([]byte, 1024*8)
	readDataLen := 0
	buf := make([]byte, 1024*4)
	bLen := len(boundary)
	reachEnd := false

	for !reachEnd {
		readLen, err := readStreamChunk(stream, buf)
		if err != nil {
			return nil, true, err
		}
		if readLen <= 0 {
			reachEnd = true
			continue
		}

		copy(readData[readDataLen:], buf[:readLen])
		readDataLen += readLen

		if readDataLen < bLen+4 {
			continue
		}

		loc := bytes.Index(readData[:readDataLen], boundary)
		if loc >= 0 {
			target.Write(readData[:loc-4])
			return readData[loc:readDataLen], reachEnd, nil
		}

		target.Write(readData[:readDataLen-bLen-4])
		copy(readData[0:], readData[readDataLen-bLen-4:])
		readDataLen = bLen + 4
	}

	target.Write(readData[:readDataLen])
	return nil, reachEnd, nil
}

// ParseFromHead parses the file header from the beginning of the stream
func ParseFromHead(readData []byte, readTotal int, boundary []byte, stream io.ReadCloser) (map[string]string, []byte, error) {
	buf := make([]byte, 1024*8)
	foundBoundary := false
	boundaryLoc := -1

	for {
		readLen, err := readStreamChunk(stream, buf)
		if err != nil {
			return nil, nil, err
		}
		if readLen <= 0 {
			break
		}

		if readTotal+readLen > cap(readData) {
			return nil, nil, fmt.Errorf("not found boundary")
		}

		copy(readData[readTotal:], buf[:readLen])
		readTotal += readLen

		if !foundBoundary {
			boundaryLoc = locateBoundary(readData, boundary, readTotal)
			if boundaryLoc == -1 {
				continue
			}
			foundBoundary = true
		}

		headMap, remainingData, ok, err := extractHeaderFromData(readData, boundary, readTotal, boundaryLoc)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			return headMap, remainingData, nil
		}
	}

	return nil, nil, fmt.Errorf("reach to stream EOF")
}
