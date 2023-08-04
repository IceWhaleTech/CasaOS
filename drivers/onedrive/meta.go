package onedrive

import (
	"github.com/IceWhaleTech/CasaOS/internal/driver"
	"github.com/IceWhaleTech/CasaOS/internal/op"
)

type Host struct {
	Oauth string
	Api   string
}

type TokenErr struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type RespErr struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
type Addition struct {
	Region       string `json:"region" type:"select" required:"true" options:"global,cn,us,de" default:"global"`
	IsSharepoint bool   `json:"is_sharepoint"`
	ClientID     string `json:"client_id" required:"true"`
	ClientSecret string `json:"client_secret" required:"true"`
	RedirectUri  string `json:"redirect_uri" required:"true" default:""`
	RefreshToken string `json:"refresh_token" required:"true"`
	SiteId       string `json:"site_id"`
	ChunkSize    int64  `json:"chunk_size" type:"number" default:"5"`
	RootFolderID string `json:"root_folder_id"`
	AuthUrl      string `json:"auth_url" type:"string" default:""`
	Icon         string `json:"icon" type:"string" default:""`
	Code         string `json:"code" type:"string" help:"code from auth_url" omit:"true"`
}
type About struct {
	Total int    `json:"total"`
	Used  int    `json:"used"`
	State string `json:"state"`
}

type Info struct {
	LastModifiedBy struct {
		Application struct {
			DisplayName string `json:"displayName"`
			ID          string `json:"id"`
		} `json:"application"`
		Device struct {
			ID string `json:"id"`
		} `json:"device"`
		User struct {
			DisplayName string `json:"displayName"`
			ID          string `json:"id"`
		} `json:"user"`
	} `json:"lastModifiedBy"`
	ParentReference struct {
		DriveID   string `json:"driveId"`
		DriveType string `json:"driveType"`
	} `json:"parentReference"`
}

var config = driver.Config{
	Name:        "Onedrive",
	LocalSort:   true,
	DefaultRoot: "/",
}

func init() {
	op.RegisterDriver(func() driver.Driver {
		one := GetConfig()
		return &one
	})
}
