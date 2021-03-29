package utils

import (
	"context"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"

	log "k8s.io/klog"
)

// LogBuilderContextKey is the key for the log builder on the context object
const LogBuilderContextKey = "lb"

// LogBuilderImpl implements the LogBuilder interface for logging during the handling of a workspace
type LogBuilderImpl struct {
	wsID string
	ns   string
}

// NewLogBuilder returns a new log builder object
func NewLogBuilder(wsID, namespace string) LogBuilder {
	return &LogBuilderImpl{
		wsID: wsID,
		ns:   namespace,
	}
}

// LogBuilder provides a method for consistent logging in a specific context
type LogBuilder interface {
	Msg(format string, args ...interface{}) string
}

// Msg wraps a message with relevant contextual formatting for log printing
func (lb *LogBuilderImpl) Msg(format string, args ...interface{}) string {
	inputMsg := fmt.Sprintf(format, args...)
	return fmt.Sprintf("[WS: %s/%s] %s", lb.ns, lb.wsID, inputMsg)
}

// GetLogBuilder returns the log builder stored in the context or a new one if it doesn't exist
func GetLogBuilder(ctx context.Context) LogBuilder {
	lb, ok := ctx.Value(LogBuilderContextKey).(LogBuilder)
	if !ok { // this isn't a logbuilder
		log.Errorf("could not retrieve logBuilder from context")
		lbImpl := NewLogBuilder("-N/A-", "-N/A-")
		return lbImpl
	}
	return lb
}

func GenerateNewContext(req ctrl.Request) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, LogBuilderContextKey, NewLogBuilder(req.Name, req.Namespace))
	return ctx
}
