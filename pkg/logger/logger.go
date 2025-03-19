package logger

import "go.uber.org/zap"

type Logger struct {
	logger *zap.Logger
}

func NewLogger() *Logger {
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return &Logger{logger: l}

}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}
