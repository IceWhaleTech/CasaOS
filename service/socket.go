package service

import (
	"net/http"
	"strings"

	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/mileusna/useragent"
)

type Name struct {
	Model       string `json:"model"`
	OS          string `json:"os"`
	Browser     string `json:"browser"`
	Type        string `json:"type"`
	DeviceName  string `json:"deviceName"`
	DisplayName string `json:"displayName"` //随机生成
}

func GetPeerId(request *http.Request, id string) string {

	cookiePree, err := request.Cookie("peerid")
	if err != nil {
		return id
	}
	if len(cookiePree.Value) > 0 {
		return cookiePree.Value
	}
	return id

}

func GetIP(request *http.Request) string {
	ip := ""
	if len(request.Header.Get("x-forwarded-for")) > 0 {
		ip = strings.Split(request.Header.Get("x-forwarded-for"), ",")[0]
	} else {
		ip = request.RemoteAddr
	}

	if ip == "::1" || ip == "::ffff:127.0.0.1" {
		ip = "127.0.0.1"
	}
	return ip
}

func GetName(request *http.Request) Name {
	us := useragent.Parse(request.Header.Get("user-agent"))

	device := ""
	if len(us.Device) > 0 {
		device += us.Device
	} else {
		device += us.Name
	}

	return Name{
		Model:       us.Device,
		OS:          us.OS,
		Browser:     us.Name,
		DeviceName:  device,
		DisplayName: Generate(),
	}
}
func GetNameByDB(m model2.PeerDriveDBModel) Name {
	device := ""
	if len(m.DeviceName) > 0 {
		device += m.DeviceName
	} else {
		device += m.Browser
	}
	return Name{
		Model:       m.DeviceName,
		OS:          m.OS,
		Browser:     m.Browser,
		DeviceName:  device,
		DisplayName: m.DisplayName,
	}
}
func Generate() string {
	nameParts := strings.Split(namesgenerator.GetRandomName(0), "_")

	for i := 0; i < len(nameParts); i++ {
		nameParts[i] = strings.Title(nameParts[i])
	}

	return strings.Join(nameParts, " ")
}

func GenerateMultiple(count int) []string {
	s := make([]string, count)
	for i := 0; i <= count; i++ {
		s[i] = Generate()
	}

	return s
}
