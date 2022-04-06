package file

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	path2 "path"
	"path/filepath"
	"strconv"
	"strings"
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
	if notExist := CheckNotExist(src); notExist {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

// MkDir create a directory
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}
	os.Chmod(src, 0777)

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
	//dir, err := os.Getwd()
	//if err != nil {
	//	return nil, fmt.Errorf("os.Getwd err: %v", err)
	//}

	src := filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err := IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

func CreateFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func CreateFileAndWriteContent(path string, content string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()
	write := bufio.NewWriter(file)

	write.WriteString(content)

	write.Flush()
	return nil
}

// IsNotExistMkDir create a directory if it does not exist
func IsNotExistCreateFile(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
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

// File copies a single file from src to dst
func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	lastPath := src[strings.LastIndex(src, "/")+1:]

	if !strings.HasSuffix(dst, "/") {
		dst += "/"
	}
	dstPath := dst
	dst += lastPath
	for i := 0; Exists(dst); i++ {
		name := strings.Split(lastPath, ".")
		nameIndex := 0
		if len(name) > 2 {
			nameIndex = len(name) - 2
		}
		name[nameIndex] = name[nameIndex] + "(Copy)"
		dst = dstPath
		for _, v := range name {
			dst += v + "."
		}
		dst = strings.TrimSuffix(dst, ".")
	}

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

//Check for duplicate file names
func GetNoDuplicateFileName(fullPath string) string {
	path, fileName := filepath.Split(fullPath)
	fileSuffix := path2.Ext(fileName)
	filenameOnly := strings.TrimSuffix(fileName, fileSuffix)
	for i := 0; Exists(fullPath); i++ {
		fullPath = path2.Join(path, filenameOnly+"("+strconv.Itoa(i+1)+")"+fileSuffix)
	}
	return fullPath
}

// Dir copies a whole directory recursively
func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	if !srcinfo.IsDir() {
		if err = CopyFile(src, dst); err != nil {
			fmt.Println(err)
		}
		return nil
	}
	dstPath := dst
	lastPath := src[strings.LastIndex(src, "/")+1:]
	dst += "/" + lastPath
	for i := 0; Exists(dst); i++ {
		dst = dstPath + "/" + lastPath + strconv.Itoa(i+1)
	}
	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}
	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := dst //path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

//文件写入临时目录
func WriteToPath(data []byte, path, name string) error {
	fullPath := path
	if strings.HasSuffix(path, "/") {
		fullPath += name
	} else {
		fullPath += "/" + name
	}
	IsNotExistCreateFile(fullPath)
	file, err := os.OpenFile(fullPath,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)

	return err
}

//最终拼接
func SpliceFiles(dir, path string, length int, startPoint int) error {

	fullPath := path

	IsNotExistCreateFile(fullPath)

	file, _ := os.OpenFile(fullPath,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	defer file.Close()
	bufferedWriter := bufio.NewWriter(file)
	for i := 0; i < length+startPoint; i++ {
		data, err := ioutil.ReadFile(dir + "/" + strconv.Itoa(i+startPoint))
		if err != nil {
			return err
		}
		_, err = bufferedWriter.Write(data)
		if err != nil {
			return err
		}
	}

	bufferedWriter.Flush()

	return nil
}
