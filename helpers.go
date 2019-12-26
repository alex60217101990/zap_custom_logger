package zap_custom_logger

import (
	"fmt"
	"time"
)

func Int32(v int32) *int32 {
	return &v
}

func Int32Ptr(v *int32) int32 {
	if v != nil {
		return *v
	}
	return 0
}

func Bool(v bool) *bool {
	return &v
}

func BoolPtr(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}

func Int(v int) *int {
	return &v
}

func IntPtr(v *int) int {
	if v != nil {
		return *v
	}
	return 0
}

func TimeToTimePtr(t time.Time) *time.Time {
	return &t
}

func TimePtrToTime(t *time.Time) (emptyTime time.Time) {
	if t != nil {
		return *t
	}
	return emptyTime
}

func String(v string) *string {
	return &v
}

func StringPtr(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func GetServiceLogsIndexName(serviceName *string) *string {
	return String(fmt.Sprintf("index_service_logs_%s", StringPtr(serviceName)))
}

func GetEndpointsLogsIndexName(serviceName *string) *string {
	return String(fmt.Sprintf("index_endpoints_logs_%s", StringPtr(serviceName)))
}
