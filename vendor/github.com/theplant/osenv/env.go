package osenv

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
)

func Get(key string, msg string, defaultValue string) string {
	return getWithCallerSkip(key, msg, defaultValue, 3)
}

func GetBool(key string, msg string, defaultValue bool) (r bool) {
	sv := getWithCallerSkip(key, msg, fmt.Sprint(defaultValue), 3)
	if sv == "" {
		return defaultValue
	}
	r, _ = strconv.ParseBool(sv)
	return
}

func GetInt64(key string, msg string, defaultValue int64) (r int64) {
	sv := getWithCallerSkip(key, msg, fmt.Sprint(defaultValue), 3)
	if sv == "" {
		return defaultValue
	}
	r, _ = strconv.ParseInt(sv, 10, 64)
	return
}

func getWithCallerSkip(key string, msg string, defaultValue string, callerSkip int) string {
	defer func() {
		pc, _, _, _ := runtime.Caller(callerSkip)
		funcName := runtime.FuncForPC(pc).Name()

		defaultMessage := ""
		if defaultValue != "" {
			defaultMessage = fmt.Sprintf(" (default: %s)", defaultValue)
		}
		fmt.Printf("%s: %s,%s, at: %s\n", key, msg, defaultMessage, funcName)
	}()

	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
