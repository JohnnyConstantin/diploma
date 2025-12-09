package drivers

import (
	"diploma/pkg/logger/message"
	"testing"
)

func TestMakeStdoutLogger(t *testing.T) {
	logger := MakeStdoutLogger(message.InfoLevel)
	if logger == nil {
		t.Fatal("MakeStdoutLogger returned nil")
	}

	msg := &message.LogMessage{Message: "test"}

	logger.Debug(msg)
	logger.Info(msg)
	logger.Warn(msg)
	logger.Error(msg)
	logger.Fatal(msg)
}
