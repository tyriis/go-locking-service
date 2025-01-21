package utils

import (
	"fmt"
	"runtime"
	"strings"
)

func GetFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	fullName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fullName, ".")
	if len(parts) < 2 {
		return fullName
	}
	// Extract service and method names
	service := parts[len(parts)-2]
	method := parts[len(parts)-1]
	return fmt.Sprintf("%s.%s", service, method)
}
