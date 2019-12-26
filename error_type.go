package zap_custom_logger

type errorType int

const (
    ServiceLogType     errorType = iota + 1
    EndpointLogType
)

var errorTypeRelations map[string]errorType

func init() {
    errorTypeRelations = map[string]errorType{
        "ServiceLogType":  ServiceLogType,
        "EndpointLogType": EndpointLogType,
    }
}

func (t *errorType) ConvertFromFloat64(data float64) {
    *t = errorType(int(data))
}

func (t errorType) String() string {
    names := [...]string{"ServiceLogType", "EndpointLogType"}
    if t < ServiceLogType || t > EndpointLogType {
        return "unknown error log type"
    }
    return names[t-1]
}

func (t *errorType) FromStrConvert(errorLogType string) *errorType {
    if result, ok := errorTypeRelations[errorLogType]; ok {
        // *t = result
        return &result
    }
    return nil
}

func (t errorType) Val() int {
    return int(t)
}
