package zap_custom_logger

type Configs struct {
	App     App
	Encoder EncoderType
	Storage Storage
}

type App struct {
	PublicIP    string
	Version     string
	ServiceName string
	Namespace   string
}

type Storage struct {
	LoggerStorage LogStorageType
	Hosts         []string
}
