package common

import (
	"github.com/IceWhaleTech/CasaOS/codegen/message_bus"
)

// devtype -> action -> event
var EventTypes = []message_bus.EventType{
	{Name: "casaos:system:utilization", SourceID: SERVICENAME, PropertyTypeList: []message_bus.PropertyType{}},
	{Name: "casaos:file:recover", SourceID: SERVICENAME, PropertyTypeList: []message_bus.PropertyType{}},
	{Name: "casaos:file:operate", SourceID: SERVICENAME, PropertyTypeList: []message_bus.PropertyType{}},
}
