package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	logger := New()

	if logger == nil {
		t.Fatal("New() returned nil")
	}

	if logger.Logger == nil {
		t.Fatal("Logger.Logger is nil")
	}
}

func TestLogger_Println(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "[APP] ", log.LstdFlags),
	}

	testMessage := "test message"
	logger.Println(testMessage)

	output := buf.String()
	if !strings.Contains(output, testMessage) {
		t.Errorf("Expected output to contain '%s', got: %s", testMessage, output)
	}

	if !strings.Contains(output, "[APP]") {
		t.Errorf("Expected output to contain '[APP]', got: %s", output)
	}
}

func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "[APP] ", log.LstdFlags),
	}

	testMessage := "info message"
	logger.Info(testMessage)

	output := buf.String()
	if !strings.Contains(output, testMessage) {
		t.Errorf("Expected output to contain '%s', got: %s", testMessage, output)
	}

	if !strings.Contains(output, "[INFO]") {
		t.Errorf("Expected output to contain '[INFO]', got: %s", output)
	}

	if !strings.Contains(output, "[APP]") {
		t.Errorf("Expected output to contain '[APP]', got: %s", output)
	}
}

func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "[APP] ", log.LstdFlags),
	}

	testMessage := "error message"
	logger.Error(testMessage)

	output := buf.String()
	if !strings.Contains(output, testMessage) {
		t.Errorf("Expected output to contain '%s', got: %s", testMessage, output)
	}

	if !strings.Contains(output, "[ERROR]") {
		t.Errorf("Expected output to contain '[ERROR]', got: %s", output)
	}

	if !strings.Contains(output, "[APP]") {
		t.Errorf("Expected output to contain '[APP]', got: %s", output)
	}
}

func TestLogger_Debug(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "[APP] ", log.LstdFlags),
	}

	testMessage := "debug message"
	logger.Debug(testMessage)

	output := buf.String()
	if !strings.Contains(output, testMessage) {
		t.Errorf("Expected output to contain '%s', got: %s", testMessage, output)
	}

	if !strings.Contains(output, "[DEBUG]") {
		t.Errorf("Expected output to contain '[DEBUG]', got: %s", output)
	}

	if !strings.Contains(output, "[APP]") {
		t.Errorf("Expected output to contain '[APP]', got: %s", output)
	}
}

func TestLogger_MultipleArguments(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "[APP] ", log.LstdFlags),
	}

	arg1 := "first"
	arg2 := "second"
	arg3 := 123

	logger.Info(arg1, arg2, arg3)

	output := buf.String()
	if !strings.Contains(output, "first") {
		t.Errorf("Expected output to contain 'first', got: %s", output)
	}

	if !strings.Contains(output, "second") {
		t.Errorf("Expected output to contain 'second', got: %s", output)
	}

	if !strings.Contains(output, "123") {
		t.Errorf("Expected output to contain '123', got: %s", output)
	}

	if !strings.Contains(output, "[INFO]") {
		t.Errorf("Expected output to contain '[INFO]', got: %s", output)
	}
}

func TestLogger_EmptyMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "[APP] ", log.LstdFlags),
	}

	logger.Info()

	output := buf.String()
	if !strings.Contains(output, "[INFO]") {
		t.Errorf("Expected output to contain '[INFO]', got: %s", output)
	}

	if !strings.Contains(output, "[APP]") {
		t.Errorf("Expected output to contain '[APP]', got: %s", output)
	}
}

func TestLogger_LogLevels(t *testing.T) {
	tests := []struct {
		name     string
		logFunc  func(*Logger, ...interface{})
		expected string
	}{
		{
			name:     "Info level",
			logFunc:  (*Logger).Info,
			expected: "[INFO]",
		},
		{
			name:     "Error level",
			logFunc:  (*Logger).Error,
			expected: "[ERROR]",
		},
		{
			name:     "Debug level",
			logFunc:  (*Logger).Debug,
			expected: "[DEBUG]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := &Logger{
				Logger: log.New(&buf, "[APP] ", log.LstdFlags),
			}

			testMessage := "test message"
			tt.logFunc(logger, testMessage)

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected output to contain '%s', got: %s", tt.expected, output)
			}

			if !strings.Contains(output, testMessage) {
				t.Errorf("Expected output to contain '%s', got: %s", testMessage, output)
			}
		})
	}
}

func TestLogger_ConcurrentAccess(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "[APP] ", log.LstdFlags),
	}

	done := make(chan bool, 3)

	go func() {
		logger.Info("concurrent info message")
		done <- true
	}()

	go func() {
		logger.Error("concurrent error message")
		done <- true
	}()

	go func() {
		logger.Debug("concurrent debug message")
		done <- true
	}()

	for i := 0; i < 3; i++ {
		<-done
	}

	output := buf.String()

	if !strings.Contains(output, "concurrent info message") {
		t.Error("Expected output to contain 'concurrent info message'")
	}

	if !strings.Contains(output, "concurrent error message") {
		t.Error("Expected output to contain 'concurrent error message'")
	}

	if !strings.Contains(output, "concurrent debug message") {
		t.Error("Expected output to contain 'concurrent debug message'")
	}
}

func BenchmarkLogger_Info(b *testing.B) {
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "[APP] ", log.LstdFlags),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message", i)
	}
}

func BenchmarkLogger_Error(b *testing.B) {
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "[APP] ", log.LstdFlags),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error("benchmark error message", i)
	}
}

func BenchmarkLogger_Debug(b *testing.B) {
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "[APP] ", log.LstdFlags),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("benchmark debug message", i)
	}
}
