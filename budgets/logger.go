package budgets

import (
	"log"

	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func SetupLogger() {

	cf := zap.NewProductionConfig()
	cf.OutputPaths = []string{"stdout", "/tmp/qbbudgets"}
	cf.ErrorOutputPaths = []string{"stderr", "/tmp/qbbudgets"}
	cf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cf.Encoding = "json"
	logg, err := cf.Build()
	if err != nil {
		panic(err)
	}

	Logger = logg.Sugar()
	log.Printf("Logger is: %v", Logger)
}
