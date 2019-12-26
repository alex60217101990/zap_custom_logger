package zap_custom_logger

import (
	"context"
	"io"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	readStdErr, writeStdErr := io.Pipe()
	readStdOut, writeStdOut := io.Pipe()
	ctx, cancel := context.WithCancel(context.Background())
	log := NewZapLogger(
		SetConfigs(&Configs{
			App: App{
				PublicIP: "127.0.0.1",
			},
			Storage: Storage{
				Hosts:         []string{"http://192.168.64.29:31558"},
				LoggerStorage: Elastic,
			},
			Encoder: Console,
		}),
		SetWriters(writeStdOut, writeStdErr),
		SetContext(context.Background()),
	)
	log.Connect()
	syncService := NewSyncLogsService(ctx, log, readStdOut, readStdErr)
	syncService.RunLogsLoops()
	log.ErrorEndpoint(&EndpointLog{
		Error: String("djbvdvbbvrfkvbruv"),
	})
	defer func() {
		syncService.Close()
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
				/*log.ErrorEndpoint(&zap_custom_logger.EndpointLog{
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
				})*/
				log.ErrorService(&ServiceLog{
					ServiceName:      "efbfjfvefvrgfrujg",
					InstancePublicIP: "fejfefuefeifgrigur",
					Msg:              nil,
					Error:            String("shcjdvdfefyefru"),
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
