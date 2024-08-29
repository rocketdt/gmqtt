package zap

type Logger struct{}
type Field = interface{}

func String(s1 string, s2 string) string {
	return ""
}

func NewProduction() (Logger, error) {
	return Logger{}, nil
}

func NewNop() *Logger {
	return nil
}

func (x *Logger) Error(args ...interface{}) {
	// do nothing
}

func (x *Logger) Debug(args ...interface{}) {
	// do nothing
}

func (x *Logger) Info(args ...interface{}) {
	// do nothing
}

func (x *Logger) Warn(args ...interface{}) {
	// do nothing
}

func (x *Logger) With(args ...interface{}) *Logger {
	return nil
}

func Uint8(args ...interface{}) uint8 {
	return 0
}

func Strings(args ...interface{}) string {
	return ""
}

func Int16(args ...interface{}) int16 {
	return 0
}
