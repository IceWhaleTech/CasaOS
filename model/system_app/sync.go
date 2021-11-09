package system_app

import "encoding/xml"

type SyncConfig struct {
	XMLName xml.Name `xml:"configuration"`
	Key     string   `xml:"gui>apikey"`
}
