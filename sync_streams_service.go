package zap_custom_logger

import (
    "context"
    "encoding/json"
    "io"
    "log"
    "runtime/debug"
)

type SyncService interface {
    Close()
    RunLogsLoops()
}

type SyncLogsService struct {
    plugin    IndexPlugin
    readerStd *io.PipeReader
    readerErr *io.PipeReader
    globalCtx context.Context
}

func NewSyncLogsService(ctx context.Context, storageType LogStorageType, readers ...*io.PipeReader) *SyncLogsService {
    if readers == nil || len(readers) < 2 {
        log.Fatal("invalid count of readers parameter", string(debug.Stack()))
	}
	syncService := &SyncLogsService{
        readerStd: readers[0],
        readerErr: readers[1],
        globalCtx: ctx,
    }
	switch storageType {
	case Elastic:
		syncService.plugin = &ElasticPlugin{}
		err := syncService.plugin.Connect()
		if err != nil {
			log.Fatal(err, string(debug.Stack()))
		}
	case Loki:
		//
	}
    return syncService
}

func (s *SyncLogsElasticService) Close() {
    err := s.readerStd.Close()
    if err != nil {
        CustomLogger.WarnService(&ServiceLog{
            Error: err,
        })
    }
    err = s.readerErr.Close()
    if err != nil {
        CustomLogger.WarnService(&ServiceLog{
            Error: err,
        })
    }
    s.engine.Close()
}

func (s *SyncLogsElasticService) RunLogsLoops() {
   s.InfoHighLevelLogsLoop()
   s.InfoLowLevelLogsLoop()
}

func (s *SyncLogsElasticService) InfoHighLevelLogsLoop() {
    go func() {
        defer func() {
            r := recover()
            if err, ok := r.(error); ok {
                CustomLogger.ErrorService(&ServiceLog{
                    Error: err.Error(),
                })
            }
        }()
        dec := json.NewDecoder(s.readerStd)
        var tmpLog ServiceLog
        var endpointLog EndpointLog
    Loop:
        for {
            select {
            default:
                err := dec.Decode(&tmpLog)
                if err == io.EOF || err == io.ErrUnexpectedEOF  {
                    continue
                } else if err != nil {
                    if err = s.decodeEndpointLogHelper(dec, &endpointLog); err != nil {
                        CustomLogger.WarnService(&ServiceLog{
                            Error: err.Error(),
                        })
                        continue
                    }
                }
                if tmpLog.errorType != nil && *tmpLog.errorType == ServiceLogType {
                    err = s.engine.InsertServiceLogObjectStruct(tmpLog)
                    if err != nil {
                        CustomLogger.WarnService(&ServiceLog{
                            Error: err.Error(),
                        })
                        continue
                    }
                } else if tmpLog.errorType != nil && *tmpLog.errorType == EndpointLogType {
                    if err = s.decodeEndpointLogHelper(dec, &endpointLog); err != nil {
                        CustomLogger.WarnService(&ServiceLog{
                            Error: err.Error(),
                        })
                        continue
                    }
                }
            case <-s.globalCtx.Done():
                break Loop
            }
        }
    }()
}

func (s *SyncLogsElasticService) decodeEndpointLogHelper(dec *json.Decoder, endpointLog *EndpointLog) (err error) {
    if err = dec.Decode(endpointLog); err != nil {
        return err
    }
    return s.engine.InsertEndpointLogObjectStruct(*endpointLog)
}

func (s *SyncLogsElasticService) InfoLowLevelLogsLoop() {
    go func() {
        defer func() {
            r := recover()
            if err, ok := r.(error); ok {
                CustomLogger.ErrorService(&ServiceLog{
                    Error: err.Error(),
                })
            }
        }()
        dec := json.NewDecoder(s.readerErr)
        var tmpLog ServiceLog
        var endpointLog EndpointLog
    Loop:
        for {
            select {
            default:
                err := dec.Decode(&tmpLog)
                if err == io.EOF || err == io.ErrUnexpectedEOF  {
                    continue
                } else if err != nil {
                    if err = s.decodeEndpointLogHelper(dec, &endpointLog); err != nil {
                        CustomLogger.WarnService(&ServiceLog{
                            Error: err.Error(),
                        })
                        continue
                    }
                }
                if tmpLog.errorType != nil && *tmpLog.errorType == ServiceLogType {
                    err = s.engine.InsertServiceLogObjectStruct(tmpLog)
                    if err != nil {
                        CustomLogger.WarnService(&ServiceLog{
                            Error: err.Error(),
                        })
                        continue
                    }
                } else if tmpLog.errorType != nil && *tmpLog.errorType == EndpointLogType {
                    if err = s.decodeEndpointLogHelper(dec, &endpointLog); err != nil {
                        CustomLogger.WarnService(&ServiceLog{
                            Error: err.Error(),
                        })
                        continue
                    }
                }
            case <-s.globalCtx.Done():
                break Loop
            }
        }
    }()
}