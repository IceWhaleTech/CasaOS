package google_drive

import (
	"github.com/IceWhaleTech/CasaOS/internal/driver"
	"github.com/IceWhaleTech/CasaOS/internal/op"
)

type Addition struct {
	driver.RootID
	RefreshToken   string `json:"refresh_token" required:"true"`
	OrderBy        string `json:"order_by" type:"string" help:"such as: folder,name,modifiedTime"`
	OrderDirection string `json:"order_direction" type:"select" options:"asc,desc"`
	ClientID       string `json:"client_id" required:"true" default:"865173455964-4ce3gdl73ak5s15kn1vkn73htc8tant2.apps.googleusercontent.com"`
	ClientSecret   string `json:"client_secret" required:"true" default:"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"`
	ChunkSize      int64  `json:"chunk_size" type:"number" default:"5" help:"chunk size while uploading (unit: MB)"`
	AuthUrl        string `json:"auth_url" type:"string" default:"https://accounts.google.com/o/oauth2/auth/oauthchooseaccount?response_type=code&client_id=865173455964-4ce3gdl73ak5s15kn1vkn73htc8tant2.apps.googleusercontent.com&redirect_uri=http%3A%2F%2Ftest-get.casaos.io&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fdrive&access_type=offline&approval_prompt=force&state=${HOST}%2Fv1%2Frecover%2FGoogleDrive&service=lso&o2v=1&flowName=GeneralOAuthFlow"`
	Icon           string `json:"icon" type:"string" default:"https://i.pcmag.com/imagery/reviews/02PHW91bUvLOs36qNbBzOiR-12.fit_scale.size_760x427.v1569471162.png"`
	Code           string `json:"code" type:"string" help:"code from auth_url"`
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
