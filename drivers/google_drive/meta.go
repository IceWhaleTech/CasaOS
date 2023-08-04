package google_drive

import (
	"github.com/IceWhaleTech/CasaOS/internal/driver"
	"github.com/IceWhaleTech/CasaOS/internal/op"
)

const ICONURL = "./img/driver/GoogleDrive.svg"

type Addition struct {
	driver.RootID
	RefreshToken   string `json:"refresh_token" required:"true" omit:"true"`
	OrderBy        string `json:"order_by" type:"string" help:"such as: folder,name,modifiedTime" omit:"true"`
	OrderDirection string `json:"order_direction" type:"select" options:"asc,desc" omit:"true"`
	ClientID       string `json:"client_id" required:"true" default:"" omit:"true"`
	ClientSecret   string `json:"client_secret" required:"true" default:"" omit:"true"`
	ChunkSize      int64  `json:"chunk_size" type:"number" help:"chunk size while uploading (unit: MB)" omit:"true"`
	AuthUrl        string `json:"auth_url" type:"string" default:""`
	Icon           string `json:"icon" type:"string" default:"./img/driver/GoogleDrive.svg"`
	Code           string `json:"code" type:"string" help:"code from auth_url" omit:"true"`
}

var config = driver.Config{
	Name:        "GoogleDrive",
	OnlyProxy:   true,
	DefaultRoot: "root",
}

func init() {
	op.RegisterDriver(func() driver.Driver {
		google := GetConfig()
		return &google
	})
}
