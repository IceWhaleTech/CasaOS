package config

import (
	"runtime"

	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
)

//检查目录是否存在
func mkdirDATAAll() {
	sysType := runtime.GOOS
	var dirArray []string
	if sysType == "linux" {
		dirArray = []string{"/DATA/AppData", "/DATA/Documents", "/DATA/Downloads", "/DATA/Gallery", "/DATA/Media/Movies", "/DATA/Media/TV Shows", "/DATA/Media/Music"}
	}

	if sysType == "windows" {
		dirArray = []string{"C:\\CasaOS\\DATA\\AppData", "C:\\CasaOS\\DATA\\Documents", "C:\\CasaOS\\DATA\\Downloads", "C:\\CasaOS\\DATA\\Gallery", "C:\\CasaOS\\DATA\\Media/Movies", "C:\\CasaOS\\DATA\\Media\\TV Shows", "C:\\CasaOS\\DATA\\Media\\Music"}
	}
	if sysType == "darwin" {
		dirArray = []string{"./CasaOS/DATA/AppData", "./CasaOS/DATA/Documents", "./CasaOS/DATA/Downloads", "./CasaOS/DATA/Gallery", "./CasaOS/DATA/Media/Movies", "./CasaOS/DATA/Media/TV Shows", "./CasaOS/DATA/Media/Music"}
	}

	for _, v := range dirArray {
		file.IsNotExistMkDir(v)
	}

}

func UpdateSetup() {
	mkdirDATAAll()
}
