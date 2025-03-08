package security

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func JsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(formatLog)
}

type LogData struct {
	StatusCode   int    `json:"status_code"`
	Path         string `json:"path"`
	Method       string `json:"method"`
	StartTime    string `json:"start_time"`
	RemoteAddr   string `json:"remote_addr"`
	ResponseTime string `json:"response_time"`
}

func createLogData(params gin.LogFormatterParams) LogData {
	return LogData{
		StatusCode:   params.StatusCode,
		Path:         params.Path,
		Method:       params.Method,
		StartTime:    params.TimeStamp.Format("2006/01/02 - 15:04:05"),
		RemoteAddr:   params.ClientIP,
		ResponseTime: params.Latency.String(),
	}
}

func formatLog(params gin.LogFormatterParams) string {
	logData := createLogData(params)
	logJSON, err := json.Marshal(logData)
	if err != nil {
		return `{"error": "failed to format log"}` + "\n"
	}
	return string(logJSON) + "\n"
}
