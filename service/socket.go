package service

import (
	"net/http"
	"strconv"
	"strings"

	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/mileusna/useragent"
)

type Name struct {
	Model       string `json:"model"`
	OS          string `json:"os"`
	Browser     string `json:"browser"`
	DeviceName  string `json:"deviceName"`
	DisplayName string `json:"displayName"`
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

	display := ""
	if len(us.Device) > 0 {
		display = us.Device + " " + us.Name
	} else {
		display = us.OS + " " + us.Name
	}

	model := "desktop"
	if us.Mobile {
		model = "mobile"
	}
	if us.Tablet {
		model = "tablet"
	}
	peer := MyService.Peer().GetPeerByName(display)
	if len(peer.ID) > 0 {
		for i := 0; true; i++ {
			peer = MyService.Peer().GetPeerByName(display + "_" + strconv.Itoa(i+1))
			if len(peer.ID) == 0 {
				display = display + "_" + strconv.Itoa(i+1)
				break
			}
		}
	}

	return Name{
		Model:       model,
		OS:          us.OS,
		Browser:     us.Name,
		DeviceName:  device,
		DisplayName: display,
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
