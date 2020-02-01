package zap_custom_logger

import (
	"context"

	"go.uber.org/zap"
)

type Logger interface {
	GetConfigs() *Configs
	Ping(context.Context) bool
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
	GetZapNativeLogger() *zap.Logger
}
