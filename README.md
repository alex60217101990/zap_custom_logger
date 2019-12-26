# ZAP custom logger

## Installation

Run command on you [$GOPATH/src] path:

```bash
go get -u github.com/alex60217101990/zap_custom_logger
```

If you use: ```bash 
go mod``` 
the package will load automatically with other dependencies.

## Usage

The package provides two basic models of the logbook, according to which indexes are created in Elasticsearch: 

```go
type EndpointLog struct {
	Level            string        `json:"level"`
	Msg              *string       `json:"msg,omitempty"`
	Timestamp        time.Time     `json:"timestamp,omitempty"`
	Error            interface{}   `json:"error,omitempty"`
	ServiceName      string        `json:"service_name"`
	InstancePublicIP string        `json:"instance_public_ip"`
	Version          string        `json:"version"`
	Namespace        *string       `json:"namespace"`
	StackTrace       *string       `json:"stacktrace,omitempty"`
	Url              *string       `json:"url"`
	TraceID          interface{}   `json:"trace_id"`
	RequestBody      interface{}   `json:"request_body"`
	ResponseBody     interface{}   `json:"response_body"`
	Duration         time.Duration `json:"duration"`
}
``` 
and
```go
type ServiceLog struct {
	Level            string      `json:"level"`
	Msg              *string     `json:"msg,omitempty"`
	Timestamp        time.Time   `json:"timestamp,omitempty"`
	Error            interface{} `json:"error,omitempty"`
	ServiceName      string      `json:"service_name"`
	InstancePublicIP string      `json:"instance_public_ip"`
	StackTrace       *string     `json:"stacktrace,omitempty"`
	Version          string      `json:"version"`
	Namespace        *string     `json:"namespace"`
}
```


An example of creating a logger with some configuration parameters (for correct closure, it is desirable to convey a global context):

```go
	ctx, cancel := context.WithCancel(context.Background())
	log := zap_custom_logger.NewZapLogger(
		zap_custom_logger.SetConfigs(&zap_custom_logger.Configs{
			App: zap_custom_logger.App{
                PublicIP: "<you local or global IP>",
                Version: "<some version for example: 0.0.1>",
                Namespace: "<some namespac: default - "default">",
                ServiceName: "<some service name>",
			},
			Storage: zap_custom_logger.Storage{
				Hosts:         []string{"<some host>"},
				LoggerStorage: zap_custom_logger.Elastic,
			},
			Encoder: zap_custom_logger.Console,
		}),
		zap_custom_logger.SetContext(context.Background()),
	)
	log.Connect()
```
Correct close forexample:
```go
    defer func() {
		log.Close()
		cancel()
	}()
```

## License
[MIT](https://choosealicense.com/licenses/mit/)
