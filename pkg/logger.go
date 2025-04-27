package pkg

import (
	"context"
	"errors"
	"sync"
)

// logMu is a mutex to control changes on watson.LogData
var logMu sync.RWMutex

// LogData return a new LoggerData with all fields.
// The copy is returned to avoid clients change the original values
func LogData(ctx context.Context) map[string]interface{} {
	logMu.RLock()
	defer logMu.RUnlock()

	logData, ok := ctx.Value("log-data").(map[string]interface{})
	if !ok {
		return make(map[string]interface{}, 0)
	}

	d := make(map[string]interface{}, len(logData))
	for k, v := range logData {
		d[k] = v
	}

	return d
}

// AddLogData adds to the log data field that will be kept during the whole request.
// Passing data as nill would delete the corresponding field.
func AddLogData(ctx context.Context, field string, data interface{}) error {
	logMu.Lock()
	defer logMu.Unlock()

	logData, ok := ctx.Value("log-data").(map[string]interface{})
	if !ok {
		logData = make(map[string]interface{})
	}

	if data == nil {
		delete(logData, field)
		return nil
	}
	logData[field] = data

	if !ok {
		return errors.New("LogData was not initialized in this context")
	}
	return nil
}

// ContextWithEmptyLogData returns a copy of parent context with an empty map of log data
func ContextWithEmptyLogData(ctx context.Context) context.Context {
	return context.WithValue(ctx, "log-data", make(map[string]interface{}))
}

// ContextWithLogData adds to the log data map to the existing log data in context, or create the log data if it hasn't been initialized yet.
func ContextWithLogData(ctx context.Context, additionalLogData map[string]interface{}) context.Context {
	logMu.Lock()
	defer logMu.Unlock()

	logData, ok := ctx.Value("log-data").(map[string]interface{})
	if !ok {
		return context.WithValue(ctx, "log-data", additionalLogData)
	}
	for k, v := range additionalLogData {
		logData[k] = v
	}
	return context.WithValue(ctx, "log-data", logData)
}
