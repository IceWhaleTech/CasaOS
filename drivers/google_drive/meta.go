package google_drive

import (
	"github.com/IceWhaleTech/CasaOS/internal/driver"
	"github.com/IceWhaleTech/CasaOS/internal/op"
)

const ICONURL = "./img/driver/GoogleDrive.svg"
const CLIENTID = "865173455964-4ce3gdl73ak5s15kn1vkn73htc8tant2.apps.googleusercontent.com"
const CLIENTSECRET = "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"

type Addition struct {
	driver.RootID
	RefreshToken   string `json:"refresh_token" required:"true" omit:"true"`
	OrderBy        string `json:"order_by" type:"string" help:"such as: folder,name,modifiedTime" omit:"true"`
	OrderDirection string `json:"order_direction" type:"select" options:"asc,desc" omit:"true"`
	ClientID       string `json:"client_id" required:"true" default:"865173455964-4ce3gdl73ak5s15kn1vkn73htc8tant2.apps.googleusercontent.com" omit:"true"`
	ClientSecret   string `json:"client_secret" required:"true" default:"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" omit:"true"`
	ChunkSize      int64  `json:"chunk_size" type:"number" help:"chunk size while uploading (unit: MB)" omit:"true"`
	AuthUrl        string `json:"auth_url" type:"string" default:"https://accounts.google.com/o/oauth2/auth/oauthchooseaccount?response_type=code&client_id=865173455964-4ce3gdl73ak5s15kn1vkn73htc8tant2.apps.googleusercontent.com&redirect_uri=http%3A%2F%2Ftest-get.casaos.io&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fdrive&access_type=offline&approval_prompt=force&state=${HOST}%2Fv1%2Frecover%2FGoogleDrive&service=lso&o2v=1&flowName=GeneralOAuthFlow"`
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
		return &GoogleDrive{}
	})
}
