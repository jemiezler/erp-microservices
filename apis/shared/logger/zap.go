package logger

import (
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const (
	ColorReset  = "\033[0m"
	ColorBlue   = "\033[1;34m"
	ColorCyan   = "\033[0;36m"
	ColorYellow = "\033[1;33m"
	ColorGreen  = "\033[0;32m"
	ColorRed    = "\033[0;31m"
)

func GetConfig(serviceName string) logger.Config {
	return logger.Config{
		Format: fmt.Sprintf("%s[%s]%s ${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} ${error}\n", ColorCyan, serviceName, ColorReset),
		TimeFormat: "15:04:05",
	}
}

func Info(service, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s[%s] INFO: %s%s\n", ColorBlue, service, msg, ColorReset)
}

func Success(service, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s[%s] SUCCESS: %s%s\n", ColorGreen, service, msg, ColorReset)
}

func Error(service, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s[%s] ERROR: %s%s\n", ColorRed, service, msg, ColorReset)
}
