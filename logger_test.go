package logrusx_test

import (
	"testing"

	"github.com/bgp59/logrusx"

	logrusx_testutils "github.com/bgp59/logrusx/testutils"
)

func testLogConfig(t *testing.T, cfg *logrusx.LoggerConfig) {
	rootLogger := logrusx.NewCollectableLogger()
	err := rootLogger.SetLogger(cfg)
	if err != nil {
		t.Fatal(err)
	}
	rootLogger.AddCallerSrcPathPrefix(1)

	tlc := logrusx_testutils.NewTestCollectableLogger(t, rootLogger, nil)
	defer tlc.RestoreLog()

	log1 := rootLogger.NewCompLogger("Comp1")
	log2 := rootLogger.NewCompLogger("Comp2")

	log1.Debug("debug test")
	log1.Info("info test")
	log1.Warn("warn test")
	log1.Error("error test")

	log2.Debug("debug test")
	log2.Info("info test")
	log2.Warn("warn test")
	log2.Error("error test")
}

func TestLogConfig(t *testing.T) {
	for _, cfg := range []*logrusx.LoggerConfig{
		{
			UseJson: false,
		},
		{
			UseJson: false,
			Level:   "debug",
		},
		{
			UseJson: true,
		},
	} {
		t.Run("", func(t *testing.T) { testLogConfig(t, cfg) })
	}
}
