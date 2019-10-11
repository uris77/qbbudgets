package budgets

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func ProductionLogger() {

	cf := zap.NewProductionConfig()
	cf.OutputPaths = []string{"/var/log/qbbudgets/logs"}
	cf.ErrorOutputPaths = []string{"/var/log/qbbudgets/logs"}
	cf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cf.Encoding = "json"
	logg, err := cf.Build()
	if err != nil {
		panic(err)
	}

	Logger = logg.Sugar()
}

func DevelopmentLogger() {
	Logger = zap.NewExample().Sugar()
}
