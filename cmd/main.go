package main

import (
	"context"

	"github.com/alex6021710/zap_custom_logger"
)

func main() {
	log := zap_custom_logger.NewZapLogger(
		zap_custom_logger.SetConfigs(&zap_custom_logger.Configs{
			App: zap_custom_logger.App{
				LoggerStorage: zap_custom_logger.Elastic,
				PublicIP:      "127.0.0.1",
			},
			Encoder: zap_custom_logger.Console,
		}),
		zap_custom_logger.SetContext(context.Background()),
	)
	log.Connect()
	log.InfoEndpoint(&zap_custom_logger.EndpointLog{
		Msg: zap_custom_logger.String("djbvdvbbvrfkvbruv"),
	})
}
