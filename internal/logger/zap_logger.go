package logger

import (
	"go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
 }

 func NewZapLogger(writeStream io.Writer) *ZapLogger {
    if global_helpers.StringPtr(configs.App.LoggerStorage) != consts.NoneLoggerStorage &&
        writeStream == nil && len(writeStream) < 2 {
        log.Fatal("invalid count of io.Write streams for log data synchronization", string(debug.Stack()))
    }
    global_helpers.InitPublicIP()
    global_helpers.SetGlobalVariablesForApmTransport()
    zapLogger := &ZapLogger{}

    highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
        return level >= zapcore.WarnLevel
    })
    lowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
        return level < zapcore.WarnLevel
    })

    var stackOptions []zap.Option
    // add stack tracing
    stackOptions = append(stackOptions, zap.AddStacktrace(zap.FatalLevel))
    stackOptions = append(stackOptions, zap.AddStacktrace(zap.PanicLevel))
    stackOptions = append(stackOptions, zap.AddStacktrace(zap.ErrorLevel))
    stackOptions = append(stackOptions, zap.AddCallerSkip(1))
    // add filebeat hook
    stackOptions = append(stackOptions, zap.WrapCore((&apmzap.Core{}).WrapCore))

    config := zap.NewDevelopmentEncoderConfig()
    if configs.Conf.IsDebug {
        config.EncodeLevel = zapcore.CapitalColorLevelEncoder
    }
    // sync ZapCore write streams with parameter write streams
    var core zapcore.Core
    var coreIndexHigh zapcore.Core
    if global_helpers.StringPtr(configs.Conf.App.LoggerStorage) == consts.NoneLoggerStorage {
        encoderCfg := zap.NewProductionEncoderConfig()
        encoderCfg.TimeKey = "timestamp"
        encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

        consoleEncoderDebug := zapcore.NewConsoleEncoder(config)
        consoleEncoderProd := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
        consoleDebugging := zapcore.Lock(os.Stdout)
        consoleErrors := zapcore.Lock(os.Stderr)
        if configs.Conf.IsDebug {
            core = zapcore.NewTee(
                zapcore.NewCore(consoleEncoderDebug, consoleErrors, highPriority),
                zapcore.NewCore(consoleEncoderDebug, consoleDebugging, lowPriority),
            )
        } else {
            core = zapcore.NewTee(
                zapcore.NewCore(consoleEncoderProd, consoleErrors, highPriority),
                zapcore.NewCore(consoleEncoderProd, consoleDebugging, lowPriority),
            )
        }
    } else if len(writeStream) < 2 || writeStream == nil {
        log.Fatal("invalid count of io.Write or io.PipeWriter streams for log data synchronization", string(debug.Stack()))
    } else {
        topicLogsHigh := zapcore.AddSync(writeStream[0])

        encoderCfg := zap.NewProductionEncoderConfig()
        encoderCfg.TimeKey = "timestamp"
        encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

        indexEngineEncoder := zapcore.NewJSONEncoder(encoderCfg)

        consoleEncoderDebug := zapcore.NewConsoleEncoder(config)
        consoleEncoderProd := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
        consoleDebugging := zapcore.Lock(os.Stdout)
        consoleErrors := zapcore.Lock(os.Stderr)

        coreIndexHigh = zapcore.NewTee(
            zapcore.NewCore(indexEngineEncoder, topicLogsHigh, highPriority),
        )
        core = zapcore.NewTee(
            zapcore.NewCore(consoleEncoderProd, consoleErrors, lowPriority),
            zapcore.NewCore(consoleEncoderDebug, consoleDebugging, highPriority),
        )
    }
    zapLogger.logger = zap.New(core, stackOptions...).Sugar()
    zapLogger.loggerToHigh = zap.New(coreIndexHigh, stackOptions...).Sugar()
    return zapLogger
}

func (l *ZapLogger) Close() {
    if err := l.logger.Sync(); err != nil {
        log.Printf("cancel zap logger error: %v", err)
    }
    if err := l.loggerToHigh.Sync(); err != nil {
        log.Printf("cancel zap logger error: %v", err)
    }
}