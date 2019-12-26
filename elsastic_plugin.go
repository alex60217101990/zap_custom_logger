package zap_custom_logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/olivere/elastic"
)

type ElasticPlugin struct {
	Client *elastic.Client
	logger Logger
}

func (e *ElasticPlugin) Connect(logger Logger) (err error) {
	e.logger = logger
	e.Client, err = elastic.NewClient(
		elastic.SetSniff(false),
		//elastic.SetSniff(global_helpers.BoolPtr(configs.Conf.ElasticSearch.Sniff)),
		//elastic.SetSnifferInterval(time.Minute * 5),
		elastic.SetURL(strings.Join(e.logger.GetConfigs().Storage.Hosts, ", ")),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(log.New(os.Stderr, "[ELASTIC ERROR] ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "[ELASTIC INFO] ", log.LstdFlags)),
		//elastic.SetTraceLog(log.New(os.Stderr, "[ELASTIC TRACE] ", log.LstdFlags)),
		elastic.SetGzip(true),
	)
	if err != nil {
		return err
	}
	// Ping the ElasticSearch server to get e.g. the version number
	info, code, err := e.Client.Ping(e.logger.GetConfigs().Storage.Hosts[0]).Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("ElasticSearch run with code %d and version %s...", code, info.Version.Number)
	if err = e.SetMapping(); err != nil {
		return err
	}
	return nil
}

// Ping the ElasticSearch server
func (e *ElasticPlugin) Ping(ctx context.Context) bool {
	_, code, err := e.Client.Ping(e.logger.GetConfigs().Storage.Hosts[0]).Do(ctx)
	if err != nil || code > 220 || code < 200 {
		return false
	}
	return true
}

func (e *ElasticPlugin) Close() {
	e.Client.Stop()
}

func (e *ElasticPlugin) RemoveIndexByName(name *string) error {
	exists, err := e.Client.IndexExists(*name).Do(context.Background())
	if err != nil {
		return err
	}
	if exists {
		deleteIndex, err := e.Client.DeleteIndex(*name).Do(context.Background())
		if err != nil {
			return err
		}
		e.logger.ErrorService(&ServiceLog{
			Error: fmt.Sprintf("elasticsearch delete index: %s, with response: %v\n", *name, deleteIndex),
		})
	}
	return fmt.Errorf("index: %s is not exists", *name)
}

func (e *ElasticPlugin) SetMapping() (err error) {
	err = e.SetServiceLogsMap()
	if err != nil {
		return err
	}
	return e.SetEndpointsLogsMap()
}

func (e *ElasticPlugin) SetServiceLogsMap() error {
	exists, err := e.Client.IndexExists(
		StringPtr(
			GetServiceLogsIndexName(
				String(e.logger.GetConfigs().App.ServiceName),
			),
		),
	).Do(context.Background())
	if err != nil {
		return err
	}
	if !exists {
		mappingServiceLog := `
{
   "settings":{
      "number_of_shards":5,
      "number_of_replicas":2,
      "analysis":{
         "analyzer":{
            "rebuilt_whitespace":{
               "tokenizer":"whitespace",
               "filter":[
                  "lowercase"
               ]
            }
         }
      }
   },
   "mappings":{
	"_source": {
		"excludes": ["error_type"]
	},
      "properties":{
         "level":{
            "type":"keyword"
         },
         "msg":{
            "type":"text",
            "fields":{
               "keyword":{
                  "type":"keyword",
                  "ignore_above":1500
               }
            },
            "analyzer":"rebuilt_whitespace"
         },
         "timestamp":{
            "type":"date"
         },
         "error":{
            "type":"text",
            "fields":{
               "keyword":{
                  "type":"keyword",
                  "ignore_above":1500
               }
            },
            "analyzer":"rebuilt_whitespace"
         },
         "service_name":{
            "type":"keyword"
         },
         "instance_public_ip":{
            "type":"keyword"
         },
         "version":{
            "type":"keyword"
         },
         "namespace":{
            "type":"keyword"
         },
         "stacktrace":{
            "type":"text",
            "fields":{
               "keyword":{
                  "type":"keyword",
                  "ignore_above":1500
               }
            },
            "analyzer":"rebuilt_whitespace"
		 },
		 "error_type": {
			"type": "object",
			"enabled": false,
			"include_in_all": false
		 }
      }
   }
}
`
		createIndex, err := e.Client.CreateIndex(
			StringPtr(
				GetServiceLogsIndexName(
					String(
						e.logger.GetConfigs().App.ServiceName),
				),
			),
		).Body(mappingServiceLog).Do(context.Background())
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			e.logger.ErrorService(&ServiceLog{
				Error: fmt.Sprintf("'%s' index not acknowledged.", StringPtr(GetServiceLogsIndexName(String(e.logger.GetConfigs().App.ServiceName)))),
			})

		}
	}
	return nil
}

func (e *ElasticPlugin) SetEndpointsLogsMap() error {
	exists, err := e.Client.IndexExists(
		StringPtr(
			GetEndpointsLogsIndexName(
				String(e.logger.GetConfigs().App.ServiceName),
			),
		),
	).Do(context.Background())
	if err != nil {
		return err
	}
	if !exists {
		mappingEndpointsLog := `
{
   "settings":{
      "number_of_shards":5,
      "number_of_replicas":2,
      "analysis":{
         "analyzer":{
            "rebuilt_whitespace_1":{
               "tokenizer":"whitespace",
               "filter":[
                  "lowercase"
               ]
            }
         }
      }
   },
   "mappings":{
	"_source": {
		"excludes": ["error_type"]
	},
      "properties":{
         "level":{
            "type":"keyword"
         },
         "msg":{
            "type":"text",
            "fields":{
               "keyword":{
                  "type":"keyword",
                  "ignore_above":1500
               }
            },
            "analyzer":"rebuilt_whitespace_1"
         },
         "timestamp":{
            "type":"date"
         },
         "error":{
            "type":"text",
            "fields":{
               "keyword":{
                  "type":"keyword",
                  "ignore_above":1500
               }
            },
            "analyzer":"rebuilt_whitespace_1"
         },
         "stacktrace":{
            "type":"text",
            "fields":{
               "keyword":{
                  "type":"keyword",
                  "ignore_above":1500
               }
            },
            "analyzer":"rebuilt_whitespace_1"
         },
         "service_name":{
            "type":"keyword"
         },
         "instance_public_ip":{
            "type":"keyword"
         },
         "url":{
            "type":"keyword"
         },
         "trace_id":{
            "type":"keyword"
         },
         "request_body":{
            "type":"object"
         },
         "response_body":{
            "type":"object"
         },
         "namespace":{
            "type":"keyword"
         },
         "version":{
            "type":"keyword"
         },
         "duration":{
            "type":"long"
		 },
		 "error_type": {
			"type": "object",
			"enabled": false,
			"include_in_all": false
		 }
      }
   }
}
`
		createIndex, err := e.Client.CreateIndex(
			StringPtr(
				GetEndpointsLogsIndexName(
					String(e.logger.GetConfigs().App.ServiceName),
				),
			),
		).Body(mappingEndpointsLog).Do(context.Background())
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			e.logger.ErrorService(&ServiceLog{
				Error: fmt.Sprintf("'%s' index not acknowledged.", StringPtr(GetEndpointsLogsIndexName(String(e.logger.GetConfigs().App.ServiceName)))),
			})

		}
	}
	return nil
}
