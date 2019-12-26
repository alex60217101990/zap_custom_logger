package zap_custom_logger

type Configs struct {
	App App
	Encoder EncoderType
}

type App struct {
	LoggerStorage LogStorageType
	PublicIP string
	Version string
	ServiceName string
	Namespace string
}