package boot

import (
	"fmt"

	"go.uber.org/zap"
)

func InitializeApp() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("could not initialize zap!")
		return nil
	}

	sugar := logger.Sugar()
	sugar.Info("Logger Initialized")

	return sugar
}
