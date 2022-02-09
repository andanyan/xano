package logger

var logger = NewLogger()

var (
	SetLoggerLevel = logger.SetLoggerLevel
	SetOutput      = logger.SetOutput

	Trace  = logger.Trace
	Tracef = logger.Tracef
	Debug  = logger.Debug
	Debugf = logger.Debugf
	Info   = logger.Info
	Infof  = logger.Infof
	Warn   = logger.Warn
	Warnf  = logger.Warnf
	Error  = logger.Error
	Errorf = logger.Errorf
	Fatal  = logger.Fatal
	Fatalf = logger.Fatalf
)
