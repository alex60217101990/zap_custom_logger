package zap_custom_logger

import (
	"time"

	"go.uber.org/zap"
)

func (l *ZapLogger) encodeEndpointLogForConsole(methodLevel *string, log *EndpointLog) []zap.Field {
	return []zap.Field{
		zap.String("level", *l.getLevel(methodLevel, &log.Level)),
		zap.String("url", StringPtr(log.Url)),
		zap.Any("error", log.Error),
		zap.String("service_name", l.Configs.App.ServiceName),
		zap.String("instance_public_ip", l.Configs.App.PublicIP),
		zap.Any("trace_id", log.TraceID),
		zap.Any("request_body", log.RequestBody),
		zap.Any("response_body", log.ResponseBody),
		zap.Duration("duration", log.Duration),
		zap.Time("timestamp", *l.getTimestamp(&log.Timestamp)),
		zap.String("version", l.Configs.App.Version),
		zap.String("namespace", l.getNamespace(log.Namespace)),
		zap.Int("error_type", log.errorType.Val()),
	}
}

func (l *ZapLogger) encodeServiceLogForConsole(methodLevel *string, log *ServiceLog) []zap.Field {
	return []zap.Field{
		zap.String("level", *l.getLevel(methodLevel, &log.Level)),
		zap.Any("error", log.Error),
		zap.String("service_name", l.Configs.App.ServiceName),
		zap.String("instance_public_ip", l.Configs.App.PublicIP),
		zap.Time("timestamp", *l.getTimestamp(&log.Timestamp)),
		zap.String("version", l.Configs.App.Version),
		zap.String("namespace", l.getNamespace(log.Namespace)),
		zap.Int("error_type", log.errorType.Val()),
	}
}

func (l *ZapLogger) getEndpointLogBody(methodLevel *string, log *EndpointLog) []interface{} {
	return []interface{}{
		"level", l.getLevel(methodLevel, &log.Level),
		"url", log.Url,
		"error", log.Error,
		"service_name", l.Configs.App.ServiceName,
		"instance_public_ip", l.Configs.App.PublicIP,
		"trace_id", log.TraceID,
		"request_body", log.RequestBody,
		"response_body", log.ResponseBody,
		"duration", log.Duration,
		"timestamp", l.getTimestamp(&log.Timestamp),
		"version", l.Configs.App.Version,
		"namespace", l.getNamespace(log.Namespace),
		"error_type", log.errorType,
	}
}

func (l *ZapLogger) getNamespace(n *string) string {
	if n == nil {
		return l.Configs.App.Namespace
	}
	return StringPtr(n)
}

func (l *ZapLogger) getTimestamp(t *time.Time) *time.Time {
	if TimePtrToTime(t).IsZero() {
		return TimeToTimePtr(time.Now())
	}
	return t
}

func (l *ZapLogger) getLevel(methodLevel *string, level *string) *string {
	if level != nil && StringPtr(level) != "" {
		return level
	}
	return methodLevel
}

func (l *ZapLogger) getServiceLogBody(methodLevel *string, log *ServiceLog) []interface{} {
	return []interface{}{
		"level", l.getLevel(methodLevel, &log.Level),
		"error", log.Error,
		"service_name", l.Configs.App.ServiceName,
		"instance_public_ip", l.Configs.App.PublicIP,
		"timestamp", l.getTimestamp(&log.Timestamp),
		"version", l.Configs.App.Version,
		"namespace", l.getNamespace(log.Namespace),
		"error_type", log.errorType,
	}
}

func (l *ZapLogger) PanicEndpoint(log *EndpointLog) {
	t := EndpointLogType
	log.errorType = &t
	if l.loggerJson != nil {
		l.loggerJson.Panicw(StringPtr(log.Msg),
			l.getEndpointLogBody(String("panic"), log)...,
		)
	}
	if l.loggerConsole != nil {
		l.loggerConsole.Panic(
			StringPtr(log.Msg),
			l.encodeEndpointLogForConsole(String("panic"), log)...,
		)
	}
}

func (l *ZapLogger) PanicService(log *ServiceLog) {
	t := ServiceLogType
	log.errorType = &t
	if l.loggerJson != nil {
		l.loggerJson.Panicw(StringPtr(log.Msg),
			l.getServiceLogBody(String("panic"), log)...,
		)
	}
	if l.loggerConsole != nil {
		l.loggerConsole.Panic(
			StringPtr(log.Msg),
			l.encodeServiceLogForConsole(String("panic"), log)...,
		)
	}
}

func (l *ZapLogger) InfoEndpoint(log *EndpointLog) {
	t := EndpointLogType
	log.errorType = &t
	if l.loggerJson != nil {
		l.loggerJson.Infow(StringPtr(log.Msg),
			l.getEndpointLogBody(String("info"), log)...,
		)
	}
	if l.loggerConsole != nil {
		l.loggerConsole.Info(StringPtr(log.Msg),
			l.encodeEndpointLogForConsole(String("info"), log)...,
		)
	}
}

func (l *ZapLogger) InfoService(log *ServiceLog) {
	t := ServiceLogType
	log.errorType = &t
	if l.loggerJson != nil {
		l.loggerJson.Infow(StringPtr(log.Msg),
			l.getServiceLogBody(String("info"), log)...,
		)
	}
	if l.loggerConsole != nil {
		l.loggerConsole.Info(StringPtr(log.Msg),
			l.encodeServiceLogForConsole(String("info"), log)...,
		)
	}
}

func (l *ZapLogger) WarnEndpoint(log *EndpointLog) {
	t := EndpointLogType
	log.errorType = &t
	if l.loggerJson != nil {
		l.loggerJson.Warnw(StringPtr(log.Msg),
			l.getEndpointLogBody(String("warn"), log)...,
		)
	}
	if l.loggerConsole != nil {
		l.loggerConsole.Warn(StringPtr(log.Msg),
			l.encodeEndpointLogForConsole(String("warn"), log)...,
		)
	}
}

func (l *ZapLogger) WarnService(log *ServiceLog) {
	t := ServiceLogType
	log.errorType = &t
	if l.loggerJson != nil {
		l.loggerJson.Warnw(StringPtr(log.Msg),
			l.getServiceLogBody(String("warn"), log)...,
		)
	}
	if l.loggerConsole != nil {
		l.loggerConsole.Warn(StringPtr(log.Msg),
			l.encodeServiceLogForConsole(String("warn"), log)...,
		)
	}
}

func (l *ZapLogger) ErrorEndpoint(log *EndpointLog) {
	t := EndpointLogType
	log.errorType = &t
	if l.loggerJson != nil {
		l.loggerJson.Errorw(StringPtr(log.Msg),
			l.getEndpointLogBody(String("error"), log)...,
		)
	}
	if l.loggerConsole != nil {
		ll := l.encodeEndpointLogForConsole(String("error"), log)
		l.loggerConsole.Error(StringPtr(log.Msg),
			ll...,
		)
	}
}

func (l *ZapLogger) ErrorService(log *ServiceLog) {
	t := ServiceLogType
	log.errorType = &t
	if l.loggerJson != nil {
		l.loggerJson.Errorw(StringPtr(log.Msg),
			l.getServiceLogBody(String("error"), log)...,
		)
	}
	if l.loggerConsole != nil {
		l.loggerConsole.Error(StringPtr(log.Msg),
			l.encodeServiceLogForConsole(String("error"), log)...,
		)
	}
}

func (l *ZapLogger) FatalEndpoint(log *EndpointLog) {
	t := EndpointLogType
	log.errorType = &t
	if l.loggerJson != nil {
		l.loggerJson.Fatalw(StringPtr(log.Msg),
			l.getEndpointLogBody(String("fatal"), log)...,
		)
	}
	if l.loggerConsole != nil {
		l.loggerConsole.Fatal(StringPtr(log.Msg),
			l.encodeEndpointLogForConsole(String("fatal"), log)...,
		)
	}
}

func (l *ZapLogger) FatalService(log *ServiceLog) {
	t := ServiceLogType
	log.errorType = &t
	if l.loggerJson != nil {
		l.loggerJson.Fatalw(StringPtr(log.Msg),
			l.getServiceLogBody(String("fatal"), log)...,
		)
	}
	if l.loggerConsole != nil {
		l.loggerConsole.Fatal(StringPtr(log.Msg),
			l.encodeServiceLogForConsole(String("fatal"), log)...,
		)
	}
}
