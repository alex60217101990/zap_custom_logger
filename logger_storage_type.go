package zap_custom_logger

var (
	logStorageTypeRelations map[string]LogStorageType
)

type LogStorageType int

const (
	Default LogStorageType = iota
	Elastic
	Loki
)

func init() {
	logStorageTypeRelations = map[string]LogStorageType{
		"Default": Default,
		"Elastic": Elastic,
		"Loki":    Loki,
	}
}

func (t LogStorageType) String() string {
	names := [...]string{"Default", "Elastic", "Loki"}
	if t < Default || t > Loki {
		return "unknown logger storage type"
	}
	return names[t]
}

func LogStorageFromStrConvert(t string) *LogStorageType {
	if result, ok := logStorageTypeRelations[t]; ok {
		return &result
	}
	return nil
}

func (t LogStorageType) Val() int {
	return int(t)
}
