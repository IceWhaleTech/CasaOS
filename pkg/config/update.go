package config

import "github.com/IceWhaleTech/CasaOS/pkg/utils/file"

//检查目录是否存在
func mkdirDATAAll() {
	dirArray := [7]string{"/DATA/AppData", "/DATA/Documents", "/DATA/Downloads", "/DATA/Gallery", "/DATA/Media/Movies", "/DATA/Media/TV Shows", "/DATA/Media/Music"}
	for _, v := range dirArray {
		file.IsNotExistMkDir(v)
	}
}

func UpdateSetup() {
	mkdirDATAAll()
}
