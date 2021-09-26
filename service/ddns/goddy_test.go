package ddns

import (
	"testing"
)

func TestSetConfig(t *testing.T) {
	var model GoDaddy
	model.IPV4 = "180.164.179.198"
	model.Domain = "link-liang.xyz"
	model.Secret = "secret"
	model.Key = "key"
	//model.Type=ddns.GOGADDY
	//model.SetConfig()
}
