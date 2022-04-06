package file

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"math"
	"os"
	"strconv"
)

// Get info of block
func GetBlockInfo(fileSize int64) (blockSize int, length int) {
	switch {
	case fileSize <= 1<<28: //256M
		blockSize = 1 << 17 //128kb
	case fileSize <= 1<<29: //512M
		blockSize = 1 << 18 //256kb
	case fileSize <= 1<<30: //1G
		blockSize = 1 << 19 //512kb
	case fileSize <= 1<<31: //2G
		blockSize = 1 << 20 //(mb)
	case fileSize <= 1<<32: //4G
		blockSize = 1 << 21 //2mb
	case fileSize <= 1<<33: //8G
		blockSize = 1 << 22 //4mb
	case fileSize <= 1<<34: //16g
		blockSize = 1 << 23 //8mb
	default:
		blockSize = 1 << 24 //16mb
	}
	temp := float64(fileSize) / float64(blockSize)
	length = int(math.Ceil(temp))
	return
}

//get the hash of the data
func GetHashByContent(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}

//get the hash of the data
func GetHashByPath(path string) string {
	pFile, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer pFile.Close()
	md5h := md5.New()
	io.Copy(md5h, pFile)
	return hex.EncodeToString(md5h.Sum(nil))
}

//Comparison data hash
func ComparisonHash(data []byte, hash string) bool {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:]) == hash
}

//get prefix byte length
func PrefixLength(byteLength int) []byte {
	lengthByte := []byte{'0', '0', '0', '0', '0', '0'}
	bSize := strconv.Itoa(byteLength)
	cha := 6 - len(bSize)
	for i := len(bSize); i > 0; i-- {
		lengthByte[cha+i-1] = bSize[i-1]
	}
	return lengthByte
}

//get data byte length
func DataLength(length int) []byte {
	lengthByte := []byte{'0', '0', '0', '0', '0', '0', '0', '0'}
	bSize := strconv.Itoa(length)
	cha := 8 - len(bSize)
	for i := len(bSize); i > 0; i-- {
		lengthByte[cha+i-1] = bSize[i-1]
	}
	return lengthByte
}
