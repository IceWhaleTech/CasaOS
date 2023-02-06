package common

import (
	"fmt"

	"github.com/IceWhaleTech/CasaOS/codegen/message_bus"
)

var (
	// devtype -> action -> event
	EventTypes map[string]map[string]message_bus.EventType

	PropertyNameLookupMaps = map[string]map[string]string{
		"system": {
			fmt.Sprintf("%s:%s", SERVICENAME, "utilization"): "ID_BUS",
		},
	}

	ActionPastTense = map[string]string{
		"add":    "added",
		"remove": "removed",
	}
)
