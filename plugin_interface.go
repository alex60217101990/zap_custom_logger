package zap_custom_logger

import "context"

type IndexPlugin interface {
	Connect(Logger) (err error)
	Close()
	RemoveIndexByName(name *string) error
	SetMapping() (err error)
	Ping(context.Context) bool
	// indexing data...
	InsertServiceLogObjectStruct(interface{}) error
	InsertServiceLogObjectString([]byte) error
	InsertEndpointLogObjectStruct(interface{}) error
	InsertEndpointLogObjectString([]byte) error
}
