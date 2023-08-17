package main

import (
	"github.com/wehmoen/logart/client"
	"go.uber.org/zap/zapcore"
)

func main() {
	c := client.New("6023771e-e97e-480d-885b-b1f2c3de421f", "client_test", zapcore.DebugLevel)
	c.Info("Hello World!")
	c.Warn("This is a warning!")
}
