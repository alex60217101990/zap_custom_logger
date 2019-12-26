package zap_custom_logger

import (
	"context"
	"testing"
)

func TestConnect(t *testing.T) {
	log := NewZapLogger(
		SetConfigs(&Configs{
			App: App{
				LoggerStorage: Elastic,
				PublicIP:      "127.0.0.1",
			},
			Encoder: Console,
		}),
		SetContext(context.Background()),
	)
	log.Connect()
	log.InfoEndpoint(&EndpointLog{
		Msg: String("djbvdvbbvrfkvbruv"),
	})
}
