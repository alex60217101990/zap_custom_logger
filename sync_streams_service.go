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
	Ping(context.Context) bool
}

type SyncLogsService struct {
	logger    Logger
	plugin    IndexPlugin
	readerStd *io.PipeReader
	readerErr *io.PipeReader
	globalCtx context.Context
}

func NewSyncLogsService(ctx context.Context, logger Logger, readers ...*io.PipeReader) *SyncLogsService {
	if readers == nil || len(readers) < 2 {
		log.Fatal("invalid count of readers parameter", string(debug.Stack()))
	}

	syncService := &SyncLogsService{
		readerStd: readers[0],
		readerErr: readers[1],
		globalCtx: ctx,
		logger:    logger,
	}
	switch logger.GetConfigs().Storage.LoggerStorage {
	case Elastic:
		syncService.plugin = &ElasticPlugin{}
		err := syncService.plugin.Connect(logger)
		if err != nil {
			log.Fatal(err, string(debug.Stack()))
		}
	case Loki:
		// implement Loki logging logic...
	}
	return syncService
}

func (s *SyncLogsService) Close() {
	var err error
	s.plugin.Close()
	if s.readerStd != nil {
		err = s.readerStd.Close()
		if err != nil {
			log.Printf("cancel sync plugin error: %#v, stack: %s", err, debug.Stack())
			err = nil
		}
	}
	if s.readerErr != nil {
		err = s.readerErr.Close()
		if err != nil {
			log.Printf("cancel sync plugin error: %#v, stack: %s", err, debug.Stack())
		}
	}
}

func (s *SyncLogsService) RunLogsLoops() {
	s.stdOutReadLoop()
	s.stdErrReadLoop()
}

func (s *SyncLogsService) stdOutReadLoop() {
	go func() {
		defer func() {
			r := recover()
			if err, ok := r.(error); ok {
				s.logger.ErrorService(&ServiceLog{
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
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					continue
				} else if err != nil {
					if err = s.decodeEndpointLogHelper(dec, &endpointLog); err != nil {
						s.logger.WarnService(&ServiceLog{
							Error: err.Error(),
						})
						continue
					}
				}
				if tmpLog.errorType != nil && *tmpLog.errorType == ServiceLogType {
					err = s.plugin.InsertServiceLogObjectStruct(tmpLog)
					if err != nil {
						s.logger.WarnService(&ServiceLog{
							Error: err.Error(),
						})
						continue
					}
				} else if tmpLog.errorType != nil && *tmpLog.errorType == EndpointLogType {
					if err = s.decodeEndpointLogHelper(dec, &endpointLog); err != nil {
						s.logger.WarnService(&ServiceLog{
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

func (s *SyncLogsService) Ping(ctx context.Context) bool {
	return s.plugin.Ping(ctx)
}

func (s *SyncLogsService) decodeEndpointLogHelper(dec *json.Decoder, endpointLog *EndpointLog) (err error) {
	if err = dec.Decode(endpointLog); err != nil {
		return err
	}
	return s.plugin.InsertEndpointLogObjectStruct(*endpointLog)
}

func (s *SyncLogsService) stdErrReadLoop() {
	go func() {
		defer func() {
			r := recover()
			if err, ok := r.(error); ok {
				s.logger.ErrorService(&ServiceLog{
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
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					continue
				} else if err != nil {
					if err = s.decodeEndpointLogHelper(dec, &endpointLog); err != nil {
						s.logger.WarnService(&ServiceLog{
							Error: err.Error(),
						})
						continue
					}
				}
				if tmpLog.errorType != nil && *tmpLog.errorType == ServiceLogType {
					err = s.plugin.InsertServiceLogObjectStruct(tmpLog)
					if err != nil {
						s.logger.WarnService(&ServiceLog{
							Error: err.Error(),
						})
						continue
					}
				} else if tmpLog.errorType != nil && *tmpLog.errorType == EndpointLogType {
					if err = s.decodeEndpointLogHelper(dec, &endpointLog); err != nil {
						s.logger.WarnService(&ServiceLog{
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
