// Reference code on how to use logrusx

package main

import (
	"flag"
	"os"

	"github.com/bgp59/logrusx"

	// The app logger will be defined in some package:
	"github.com/bgp59/logrusx/example/pkg"
)

var mainLogger = pkg.RootLogger.NewCompLogger("main")

func main() {
	logrusx.EnableLoggerArgs()
	flag.Parse()
	cfg := logrusx.LoggerConfigFromArgs()

	if err := pkg.RootLogger.SetLogger(cfg); err != nil {
		mainLogger.Errorf("%v\n", err)
		os.Exit(1)
	}

	// Log some messages:
	mainLogger.Error("Error")
	mainLogger.Warn("Warn")
	mainLogger.Info("Info")
	mainLogger.Debug("Debug")
	mainLogger.Trace("Trace")

	// Invoke functions that use the logger:
	pkg.Func1()
	pkg.Func2()
}
