package main

import (
	"context"
	"net/http"
	"time"

	"github.com/alex6021710/zap_custom_logger"
)

func main() {
	// readStdErr, writeStdErr := io.Pipe()
	// readStdOut, writeStdOut := io.Pipe()
	ctx, cancel := context.WithCancel(context.Background())
	log := zap_custom_logger.NewZapLogger(
		zap_custom_logger.SetConfigs(&zap_custom_logger.Configs{
			App: zap_custom_logger.App{
				PublicIP: "127.0.0.1",
			},
			Storage: zap_custom_logger.Storage{
				Hosts:         []string{"http://192.168.64.29:31558"},
				LoggerStorage: zap_custom_logger.Elastic,
			},
			Encoder: zap_custom_logger.Console,
		}),
		// zap_custom_logger.SetWriters(writeStdOut, writeStdErr),
		zap_custom_logger.SetContext(context.Background()),
	)
	log.Connect()
	// syncService := zap_custom_logger.NewSyncLogsService(ctx, log, readStdOut, readStdErr)
	// syncService.RunLogsLoops()
	defer func() {
		// syncService.Close()
		log.Close()
		cancel()
	}()
	go func() {
		intervalDump := time.Duration(1) * time.Second
		tickerDump := time.NewTicker(intervalDump)
		defer tickerDump.Stop()
	Exit:
		for {
			select {
			case <-tickerDump.C:
				log.ErrorEndpoint(&zap_custom_logger.EndpointLog{
					Level:     "panic",
					Timestamp: time.Now(),
					Error:     "cdjvbrugrigburbvrigbruvrb fbvruigbr vrgburgr",
					Namespace: zap_custom_logger.String("cjdbfeufrf"),
					Url:       zap_custom_logger.String("/dcbdjhbhcd/djvdbv/dvdvk"),
					TraceID:   "dufefirfbrfirbruig",
					RequestBody: http.Request{
						Method: "POST",
						Proto:  "cdefefrrg",
					},
					ResponseBody: http.Response{
						Status:     "OK",
						StatusCode: 500,
					},
					Duration: 500,
				})
				log.ErrorService(&zap_custom_logger.ServiceLog{
					ServiceName:      "efbfjfvefvrgfrujg",
					InstancePublicIP: "fejfefuefeifgrigur",
					Msg:              nil,
					Error:            zap_custom_logger.String("shcjdvdfefyefru"),
					Timestamp:        time.Now(),
				})
				/* logging.CustomLogger.InfoService(&logging.ServiceLog{
				    ServiceName:      "efbfjfvefvrgfrujg",
				    InstancePublicIP: "fejfefuefeifgrigur",
				    Msg:          global_helpers.String("dbjrbrugrgrgriurgr"),
				    Error:            nil,
				    Timestamp:        time.Now(),
				})*/
			case <-ctx.Done():
				break Exit
			}
		}
	}()
	time.Sleep(30 * time.Second)
}
