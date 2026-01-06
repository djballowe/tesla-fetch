package ui

type StatusLoggerMethods interface {
	Log(message string)
	Done()
}

type ProgressUpdate struct {
	Message string
	Done    bool
}

type StatusLogger struct {
	status chan ProgressUpdate
}

func NewStatusLogger(ch chan ProgressUpdate) *StatusLogger {
	return &StatusLogger{status: ch}
}

func (s *StatusLogger) Log(message string) {
	s.status <- ProgressUpdate{Message: message}
}
func (s *StatusLogger) Done() {
	s.status <- ProgressUpdate{Done: true}
}

// nil status reporter
type NoopLogger struct{}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

func (n *NoopLogger) Log(message string) {}
func (n *NoopLogger) Done()              {}
