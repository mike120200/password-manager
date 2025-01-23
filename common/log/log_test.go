package zaplog_test

import (
	"fmt"
	zaplog "password_manager/common/log"
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestLog(t *testing.T) {

	logger := zap.L()
	if logger == nil {
		fmt.Println(" create logger failed, please check zap logger")
		os.Exit(-1)
	}
	logger.Info("success")
}

func init() {
	if err := zaplog.LoggerInit(); err != nil {
		panic(err)
	}
}
