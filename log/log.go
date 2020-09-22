package log

import "log"

const (
	// LevelFatal stands for fatal log level. Logs with level other than fatal
	// will be ignored.
	LevelFatal = iota

	// LevelError stands for error log level. Logs with level other than error
	// and fatal will be ignored.
	LevelError

	// LevelPrint stands for print log level. Logs with level debug will be
	// ignored.
	LevelPrint

	// LevelDebug stands for debug log level. No logs are ignored.
	LevelDebug
)

var levelToIndicator = map[int]string{
	LevelFatal: "[F] ",
	LevelError: "[E] ",
	LevelPrint: "[P] ",
	LevelDebug: "[D] ",
}

type logFunc func(...interface{})
type logFmtFunc func(string, ...interface{})

var (
	logLevel = LevelFatal
)

// SetLogLevel sets the range of logs which must be/not be ignored. Call will
// result in log levels <= l to be processed and everything else to be ignored.
func SetLogLevel(l int) {
	logLevel = l
}

// Fatalln calls log.Fatalln if log level is at least Fatal
func Fatalln(v ...interface{}) {
	logWithFunc(LevelFatal, log.Fatalln, v...)
}

// Fatalf calls log.Fatalf if log level is at least Fatal
func Fatalf(format string, v ...interface{}) {
	logFmtWithFunc(LevelFatal, log.Fatalf, format, v...)
}

// Errorln calls log.Errorln if log level is at least Error
func Errorln(v ...interface{}) {
	logWithFunc(LevelError, log.Println, v...)
}

// Errorf calls log.Errorf if log level is at least Error
func Errorf(format string, v ...interface{}) {
	logFmtWithFunc(LevelError, log.Printf, format, v...)
}

// Println calls log.Println if log level is at least Print
func Println(v ...interface{}) {
	logWithFunc(LevelPrint, log.Println, v...)
}

// Printf calls log.Printf if log level is at least Print
func Printf(format string, v ...interface{}) {
	logFmtWithFunc(LevelPrint, log.Printf, format, v...)
}

// Debugln calls log.Debugln if log level is at least Debug
func Debugln(v ...interface{}) {
	logWithFunc(LevelDebug, log.Print, v...)
}

// Debugf calls log.Debugf if log level is at least Debug
func Debugf(format string, v ...interface{}) {
	logFmtWithFunc(LevelDebug, log.Printf, format, v...)
}

func logWithFunc(level int, logger logFunc, v ...interface{}) {
	if logLevel < level {
		return
	}
	newV := []interface{}{levelToIndicator[level]}
	newV = append(newV, v...)
	logger(newV...)
}

func logFmtWithFunc(level int, fmtLogger logFmtFunc, format string, v ...interface{}) {
	if logLevel < level {
		return
	}
	fmtLogger(levelToIndicator[level]+format, v...)
}
