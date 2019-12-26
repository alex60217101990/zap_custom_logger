package zap_custom_logger

type IndexPlugin interface {
	Connect(hosts []string) (err error)
	Close()
	RemoveIndexByName(name *string) error
	SetMapping() (err error)
	// indexing data...
	InsertServiceLogObjectStruct(log interface{}) error
	InsertServiceLogObjectString(log []byte) error
	InsertEndpointLogObjectStruct(log interface{}) error
	InsertEndpointLogObjectString(log []byte) error
}
