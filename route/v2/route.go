package v2

import (
	"github.com/IceWhaleTech/CasaOS/codegen"
)

type CasaOS struct{}

func NewCasaOS() codegen.ServerInterface {
	return &CasaOS{}
}
