package logging

import (
	"go.uber.org/zap"
)

func LogGrpcFuncCall(l *zap.Logger, serviceName zap.Field, funcName zap.Field) {
	l.Info("func called",
		serviceName,
		funcName,
	)
}

func LogGrpcFuncFinish(l *zap.Logger, serviceName zap.Field, funcName zap.Field) {
	l.Info("func finished",
		serviceName,
		funcName,
	)
}

func LogError(l *zap.Logger, serviceName zap.Field, funcName zap.Field, message string, err error) {
	l.Info("error",
		serviceName,
		funcName,
		zap.String("message", message),
		zap.Errors("detail", []error{err}),
	)
}
