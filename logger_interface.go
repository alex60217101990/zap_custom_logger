package zap_custom_logger

type Logger interface {
	GetConfigs() *Configs
	Connect()
	Close()
	// methods
	PanicEndpoint(log *EndpointLog)
	PanicService(log *ServiceLog)
	InfoEndpoint(log *EndpointLog)
	InfoService(log *ServiceLog)
	WarnEndpoint(log *EndpointLog)
	WarnService(log *ServiceLog)
	ErrorEndpoint(log *EndpointLog)
	ErrorService(log *ServiceLog)
	FatalEndpoint(log *EndpointLog)
	FatalService(log *ServiceLog)
}
