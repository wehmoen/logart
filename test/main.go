package main

import (
	"github.com/wehmoen/logart/client"
	"go.uber.org/zap/zapcore"
)

func main() {
	c := client.New("6fd46d9f-078c-47c1-8c9c-677f2c166a1b", "client_test", zapcore.DebugLevel)
	c.Info("Hello World!")
	c.Warn("This is a warning!")
}
