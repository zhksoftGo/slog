package slog

import (
	"io"
	"os"
)

// SugaredLogger definition.
// Is a fast and usable Logger, which already contains the default formatting and handling capabilities
type SugaredLogger struct {
	*Logger
	// output writer
	Out   io.Writer
	Level Level
	// if not set, will use DefaultFormatter
	formatter Formatter
}

// NewSugaredLogger create new SugaredLogger
func NewSugaredLogger(out io.Writer, level Level) *SugaredLogger {
	sl := &SugaredLogger{
		Out:    out,
		Level:  level,
		Logger: New(),
		// default value
		formatter: DefaultFormatter,
	}

	// NOTICE: use self as an log handler
	sl.AddHandler(sl)

	return sl
}

// Formatter get formatter
func (sl *SugaredLogger) Formatter() Formatter {
	return sl.formatter
}

// SetFormatter to handler
func (sl *SugaredLogger) SetFormatter(formatter Formatter) {
	sl.formatter = formatter
}

// IsHandling Check if the current level can be handling
func (sl *SugaredLogger) IsHandling(level Level) bool {
	return level >= sl.Level
}

// Handle log record
func (sl *SugaredLogger) Handle(record *Record) error {
	bts, err := sl.formatter.Format(record)
	if err != nil {
		return err
	}

	_, err = sl.Out.Write(bts)
	return err
}

// Flush all logs
func (sl *SugaredLogger) Flush() error {
	sl.FlushAll()
	return nil
}

//
// ------------------------------------------------------------
// Global std logger usage
// ------------------------------------------------------------
//

// std logger is an SugaredLogger
var std = NewSugaredLogger(os.Stdout, ErrorLevel)

// Std get std logger
func Std() *SugaredLogger {
	return std
}

// Exit runs all the logger exit handlers and then terminates the program using os.Exit(code)
func Exit(code int) {
	std.Exit(code)
}

func AddHandler(h Handler) {
	std.AddHandler(h)
}

func AddProcessor(p Processor) {
	std.AddProcessor(p)
}

// Trace logs a message at level Trace
func Trace(args ...interface{}) {
	std.Log(TraceLevel, args...)
}

// Trace logs a message at level Trace
func Tracef(format string, args ...interface{}) {
	std.Logf(TraceLevel, format, args...)
}

// Info logs a message at level Info
func Info(args ...interface{}) {
	std.Log(InfoLevel, args...)
}

// Info logs a message at level Info
func Infof(format string, args ...interface{}) {
	std.Logf(InfoLevel, format, args...)
}

// Warn logs a message at level Warn
func Warn(args ...interface{}) {
	std.Log(WarnLevel, args...)
}

// Warn logs a message at level Warn
func Warnf(format string, args ...interface{}) {
	std.Logf(WarnLevel, format, args...)
}

// Error logs a message at level Error
func Error(args ...interface{}) {
	std.Log(ErrorLevel, args...)
}

// Error logs a message at level Error
func Errorf(format string, args ...interface{}) {
	std.Logf(ErrorLevel, format, args...)
}

// Debug logs a message at level Debug
func Debug(args ...interface{}) {
	std.Log(DebugLevel, args...)
}

// Debug logs a message at level Debug
func Debugf(format string, args ...interface{}) {
	std.Logf(DebugLevel, format, args...)
}

// Fatal logs a message at level Fatal
func Fatal(args ...interface{}) {
	std.Log(FatalLevel, args...)
}

// Fatal logs a message at level Fatal
func Fatalf(format string, args ...interface{}) {
	std.Logf(FatalLevel, format, args...)
}
