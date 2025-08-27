package pkg

// Define a component sub-logger:
var compLogger = RootLogger.NewCompLogger("comp")

// Define a specialized logger with additional information:
var func2Logger = compLogger.WithField("extra_info", "text")

// Pretend function(s) illustrating component logging:
func Func1() {
	compLogger.Error("Error")
	compLogger.Warn("Warn")
	compLogger.Info("Info")
	compLogger.Debug("Debug")
	compLogger.Trace("Trace")
}

func Func2() {
	func2Logger.Error("Error")
	func2Logger.Warn("Warn")
	func2Logger.Info("Info")
	func2Logger.Debug("Debug")
	func2Logger.Trace("Trace")
}
