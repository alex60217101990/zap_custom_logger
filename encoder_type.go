package zap_custom_logger

var (
	encodeTypeRelations map[string]EncoderType
)

type EncoderType int

const (
	Console EncoderType = iota
	Json
)

func init() {
	encodeTypeRelations = map[string]EncoderType{
		"Console": Console,
		"Json": Json,
	}
}

func (t EncoderType) String() string {
	names := [...]string{"Console", "Json"}
	if t < Console || t > Json {
		return "unknown encoder type"
	}
	return names[t-1]
}

func EncoderFromStrConvert(t string) *EncoderType {
	if result, ok := encodeTypeRelations[t]; ok {
		return &result
	}
	return nil
}

func (t EncoderType) Val() int {
	return int(t)
}
