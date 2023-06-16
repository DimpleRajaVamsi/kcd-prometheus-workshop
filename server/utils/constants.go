package utils

type ContextKey string

const (
	LoggerKey        ContextKey = "logger"
	HttpErrorMessage string     = "Failed to process"
)
