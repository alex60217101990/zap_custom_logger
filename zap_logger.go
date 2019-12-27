package zap_custom_logger

import (
	"context"
	"io"
	"log"
	"os"
	"runtime/debug"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	loggerJson    *zap.SugaredLogger
	loggerConsole *zap.Logger
	Configs       *Configs
	ctx           context.Context
	writeStdErr   *io.PipeWriter
	writeStdOut   *io.PipeWriter
	syncService   SyncService
}

func NewZapLogger(options ...func(*ZapLogger) error) *ZapLogger {
	logger := &ZapLogger{}
	for _, op := range options {
		err := op(logger)
		if err != nil {
			log.Fatalf("error: %+v, stack: %s\n", err, string(debug.Stack()))
		}
	}
	if logger.Configs == nil {
		log.Fatalf("error: empty logger configs, stack: %s\n", string(debug.Stack()))
	}
	if logger.ctx == nil {
		logger.ctx = context.Background()
	}
	if len(logger.Configs.App.PublicIP) == 0 {
		logger.Configs.App.PublicIP = "localhost"
	}
	if len(logger.Configs.App.Namespace) == 0 {
		logger.Configs.App.Namespace = "default"
	}
	if len(logger.Configs.App.Version) == 0 {
		logger.Configs.App.Version = "0.0.0"
	}
	return logger
}

func (l *ZapLogger) Connect() {
	highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.ErrorLevel //zapcore.WarnLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zapcore.ErrorLevel //zapcore.WarnLevel
	})
	// add stack tracing
	var stackOptions []zap.Option
	stackOptions = append(stackOptions, zap.AddStacktrace(zap.FatalLevel))
	stackOptions = append(stackOptions, zap.AddStacktrace(zap.PanicLevel))
	stackOptions = append(stackOptions, zap.AddStacktrace(zap.ErrorLevel))
	stackOptions = append(stackOptions, zap.AddCallerSkip(1))

	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// sync ZapCore write streams with parameter write streams
	var core zapcore.Core

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	consoleEncoderDebug := zapcore.NewConsoleEncoder(config)
	consoleEncoderProd := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	consoleDebugging := zapcore.AddSync(os.Stderr)
	consoleErrors := zapcore.AddSync(os.Stdout)

	if l.Configs.Storage.LoggerStorage == Default {
		if l.Configs.Encoder == Console {
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
	} else {
		var readStdErr *io.PipeReader
		var readStdOut *io.PipeReader
		readStdErr, l.writeStdErr = io.Pipe()
		readStdOut, l.writeStdOut = io.Pipe()
		syncService := NewSyncLogsService(l.ctx, l, readStdOut, readStdErr)

		topicLogsHigh := zapcore.AddSync(l.writeStdErr)
		topicLogsLow := zapcore.AddSync(l.writeStdOut)

		indexEngineEncoder := zapcore.NewJSONEncoder(encoderCfg)
		if l.Configs.Encoder == Console {
			core = zapcore.NewTee(
				zapcore.NewCore(consoleEncoderDebug, consoleErrors, lowPriority),
				zapcore.NewCore(consoleEncoderDebug, consoleDebugging, highPriority),
				zapcore.NewCore(indexEngineEncoder, topicLogsHigh, highPriority),
				zapcore.NewCore(indexEngineEncoder, topicLogsLow, lowPriority),
			)
		} else {
			core = zapcore.NewTee(
				zapcore.NewCore(consoleEncoderProd, consoleErrors, lowPriority),
				zapcore.NewCore(consoleEncoderProd, consoleDebugging, highPriority),
				zapcore.NewCore(indexEngineEncoder, topicLogsHigh, highPriority),
				zapcore.NewCore(indexEngineEncoder, topicLogsLow, lowPriority),
			)
		}
		syncService.RunLogsLoops()
	}

	if l.Configs.Encoder == Console {
		l.loggerConsole = zap.New(core, stackOptions...)
	} else if l.Configs.Encoder == Json {
		l.loggerJson = zap.New(core, stackOptions...).Sugar()
	}
}

func (l *ZapLogger) GetConfigs() *Configs {
	return l.Configs
}

func (l *ZapLogger) Ping(ctx context.Context) bool {
	return l.syncService.Ping(ctx)
}

func (l *ZapLogger) Close() {
	var err error
	if l.loggerJson != nil {
		if err = l.loggerJson.Sync(); err != nil {
			log.Printf("cancel zap logger error: %v", err)
			err = nil
		}
	}
	if l.loggerConsole != nil {
		if err = l.loggerConsole.Sync(); err != nil {
			log.Printf("cancel zap logger error: %v", err)
			err = nil
		}
	}
	err = l.writeStdErr.Close()
	if err != nil {
		log.Printf("cancel zap logger error: %v", err)
	}
	err = l.writeStdOut.Close()
	if err != nil {
		log.Printf("cancel zap logger error: %v", err)
	}
	l.syncService.Close()
}

func SetConfigs(conf *Configs) func(*ZapLogger) error {
	return func(logger *ZapLogger) error {
		logger.Configs = conf
		return nil
	}
}

func SetContext(ctx context.Context) func(*ZapLogger) error {
	return func(logger *ZapLogger) error {
		logger.ctx = ctx
		return nil
	}
}
