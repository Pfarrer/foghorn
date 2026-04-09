package logger

import (
	"bytes"
	"strings"
	"testing"
)

func newTestLogger(level LogLevel, verbose bool) (*Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	l := &Logger{
		level:   level,
		verbose: verbose,
		output:  &buf,
	}
	return l, &buf
}

func TestDebugLevelShowsAll(t *testing.T) {
	l, buf := newTestLogger(LevelDebug, false)

	l.log(LevelDebug, "debug msg")
	l.log(LevelInfo, "info msg")
	l.log(LevelWarn, "warn msg")
	l.log(LevelError, "error msg")

	output := buf.String()
	for _, level := range []string{"[DEBUG]", "[INFO]", "[WARN]", "[ERROR]"} {
		if !strings.Contains(output, level) {
			t.Fatalf("expected %s in output, got: %s", level, output)
		}
	}
}

func TestInfoLevelSuppressesDebug(t *testing.T) {
	l, buf := newTestLogger(LevelInfo, false)

	l.log(LevelDebug, "debug msg")
	l.log(LevelInfo, "info msg")
	l.log(LevelWarn, "warn msg")
	l.log(LevelError, "error msg")

	output := buf.String()
	if strings.Contains(output, "[DEBUG]") {
		t.Fatal("debug should be suppressed at info level")
	}
	for _, level := range []string{"[INFO]", "[WARN]", "[ERROR]"} {
		if !strings.Contains(output, level) {
			t.Fatalf("expected %s in output, got: %s", level, output)
		}
	}
}

func TestErrorLevelOnlyShowsErrors(t *testing.T) {
	l, buf := newTestLogger(LevelError, false)

	l.log(LevelDebug, "debug msg")
	l.log(LevelInfo, "info msg")
	l.log(LevelWarn, "warn msg")
	l.log(LevelError, "error msg")

	output := buf.String()
	for _, level := range []string{"[DEBUG]", "[INFO]", "[WARN]"} {
		if strings.Contains(output, level) {
			t.Fatalf("%s should be suppressed at error level", level)
		}
	}
	if !strings.Contains(output, "[ERROR]") {
		t.Fatal("expected [ERROR] in output")
	}
}

func TestParseLevelValid(t *testing.T) {
	tests := []struct {
		input string
		want  LogLevel
	}{
		{"debug", LevelDebug},
		{"info", LevelInfo},
		{"warn", LevelWarn},
		{"error", LevelError},
	}
	for _, tt := range tests {
		got, err := ParseLevel(tt.input)
		if err != nil {
			t.Fatalf("ParseLevel(%q) returned error: %v", tt.input, err)
		}
		if got != tt.want {
			t.Fatalf("ParseLevel(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestParseLevelInvalid(t *testing.T) {
	_, err := ParseLevel("trace")
	if err == nil {
		t.Fatal("expected error for invalid level 'trace'")
	}
}

func TestOutputFormat(t *testing.T) {
	l, buf := newTestLogger(LevelInfo, false)
	l.log(LevelInfo, "check started")

	output := strings.TrimSpace(buf.String())
	if output != "[INFO] check started" {
		t.Fatalf("unexpected output: %q", output)
	}
}

func TestVerboseTimestamps(t *testing.T) {
	l, buf := newTestLogger(LevelInfo, true)
	l.log(LevelInfo, "test")

	output := strings.TrimSpace(buf.String())
	if !strings.HasPrefix(output, "20") {
		t.Fatalf("expected timestamp prefix in verbose mode, got: %q", output)
	}
	if !strings.Contains(output, "[INFO] test") {
		t.Fatalf("expected [INFO] test in output, got: %q", output)
	}
}

func TestNonVerboseNoTimestamp(t *testing.T) {
	l, buf := newTestLogger(LevelInfo, false)
	l.log(LevelInfo, "test")

	output := strings.TrimSpace(buf.String())
	if output != "[INFO] test" {
		t.Fatalf("expected no timestamp, got: %q", output)
	}
}

func TestSetGetGlobal(t *testing.T) {
	original := global
	global = nil
	defer func() { global = original }()

	custom := New(LevelDebug, true)
	SetGlobal(custom)

	if GetGlobal() != custom {
		t.Fatal("GetGlobal should return the logger set via SetGlobal")
	}
}

func TestGetGlobalDefault(t *testing.T) {
	original := global
	global = nil
	defer func() { global = original }()

	l := GetGlobal()
	if l == nil {
		t.Fatal("GetGlobal should return a default logger")
	}
	if l.GetLevel() != LevelInfo {
		t.Fatalf("default level should be info, got %d", l.GetLevel())
	}
	if l.IsVerbose() {
		t.Fatal("default verbose should be false")
	}
}

func TestGlobalPackageFunctions(t *testing.T) {
	original := global
	defer func() { global = original }()

	var buf bytes.Buffer
	l := &Logger{
		level:   LevelDebug,
		verbose: false,
		output:  &buf,
	}
	SetGlobal(l)

	Info("hello %s", "world")

	output := strings.TrimSpace(buf.String())
	if output != "[INFO] hello world" {
		t.Fatalf("unexpected output: %q", output)
	}
}
