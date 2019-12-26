package zap_custom_logger

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type ServiceLog struct {
	Level            string      `json:"level"`
	Msg              *string     `json:"msg,omitempty"`
	Timestamp        time.Time   `json:"timestamp,omitempty"`
	Error            interface{} `json:"error,omitempty"`
	ServiceName      string      `json:"service_name"`
	InstancePublicIP string      `json:"instance_public_ip"`
	StackTrace       *string     `json:"stacktrace,omitempty"`
	Version          string      `json:"version"`
	Namespace        *string     `json:"namespace"`
	errorType        *errorType  `json:"error_type"`
}

func (e ServiceLog) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		Level            string      `json:"level"`
		Msg              *string     `json:"msg,omitempty"`
		Timestamp        time.Time   `json:"timestamp,omitempty"`
		Error            interface{} `json:"error,omitempty"`
		ServiceName      string      `json:"service_name"`
		InstancePublicIP string      `json:"instance_public_ip"`
		StackTrace       *string     `json:"stacktrace,omitempty"`
		Version          string      `json:"version"`
		Namespace        *string     `json:"namespace"`
		ErrorType        *errorType  `json:"error_type"`
	}{
		Level:            e.Level,
		Msg:              e.Msg,
		Timestamp:        e.Timestamp,
		Error:            e.Error,
		ServiceName:      e.ServiceName,
		InstancePublicIP: e.InstancePublicIP,
		Version:          e.Version,
		Namespace:        e.Namespace,
		StackTrace:       e.StackTrace,
		ErrorType:        e.errorType,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (e *ServiceLog) UnmarshalJSON(data []byte) (err error) {
	var dat map[string]interface{}
	if err := json.Unmarshal(data, &dat); err != nil {
		return err
	}
	if tmp, ok := dat["error_type"]; ok && tmp != nil {
		if errorType, ok := tmp.(float64); ok {
			if errorType == 0 {
				return fmt.Errorf("invalid 'error_type' field value")
			}
			e.errorType = ConvertErrTypeFromFloat64(errorType)
		}
		if errorType, ok := tmp.(string); ok {
			e.errorType = e.errorType.FromStrConvert(errorType)
		}
	}
	if tmp, ok := dat["error"]; ok && tmp != nil {
		if interfaceObj, ok := tmp.(interface{}); ok && interfaceObj != nil {
			e.Error = interfaceObj
		}
	}
	if tmp, ok := dat["stacktrace"]; ok {
		if strPointer, ok := tmp.(string); ok && len(strPointer) > 0 {
			e.StackTrace = &strPointer
		}
	}
	if tmp, ok := dat["namespace"]; ok {
		if strPointer, ok := tmp.(string); ok && len(strPointer) > 0 {
			e.Namespace = &strPointer
		}
	}
	if tmp, ok := dat["version"]; ok {
		e.Version, ok = tmp.(string)
	}
	if tmp, ok := dat["service_name"]; ok {
		e.ServiceName, ok = tmp.(string)
	}
	if tmp, ok := dat["instance_public_ip"]; ok {
		e.InstancePublicIP = tmp.(string)
	}
	if tmp, ok := dat["level"]; ok {
		e.Level = tmp.(string)
	}
	if tmp, ok := dat["timestamp"]; ok {
		if timeStr, ok := tmp.(string); ok {
			e.Timestamp, err = time.Parse(time.RFC3339, timeStr)
			if err != nil {
				log.Println(err)
			}
		}
	}
	if tmp, ok := dat["msg"]; ok {
		if tmpStr, ok := tmp.(string); ok && len(tmpStr) > 0 {
			e.Msg = &tmpStr
		}
	}
	return nil
}
