package configs

import (
	"github.com/alex6021710/zap_custom_logger/enums"
)

type Configs struct {
	App App
}

type App struct {
	LoggerStorage enums.LogStorageType
}