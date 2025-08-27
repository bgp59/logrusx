// Command line args for logger

package logrusx

import (
	"flag"
	"fmt"
)

var loggerFlags = make(map[string]any)

func EnableLoggerArgs() {
	loggerFlags[LOGGER_ARGS_USE_JSON] = flag.Bool(
		LOGGER_ARGS_USE_JSON,
		LOGGER_CONFIG_USE_JSON_DEFAULT,
		"Structure the logged record in JSON",
	)

	loggerFlags[LOGGER_ARGS_LEVEL] = flag.String(
		LOGGER_ARGS_LEVEL,
		LOGGER_CONFIG_LEVEL_DEFAULT,
		fmt.Sprintf("Log level name, one of %v", GetLogLevelNames()),
	)

	loggerFlags[LOGGER_ARGS_DISABALE_SRC_FILE] = flag.Bool(
		LOGGER_ARGS_DISABALE_SRC_FILE,
		LOGGER_CONFIG_DISBALE_SRC_FILE_DEFAULT,
		"Disable the reporting of the source file:line# info",
	)

	loggerFlags[LOGGER_ARGS_LOG_FILE] = flag.String(
		LOGGER_ARGS_LOG_FILE,
		LOGGER_CONFIG_LOG_FILE_DEFAULT,
		"Log to a file or use stdout/stderr",
	)

	loggerFlags[LOGGER_ARGS_LOG_FILE_MAX_SIZE_MB] = flag.Int(
		LOGGER_ARGS_LOG_FILE_MAX_SIZE_MB,
		LOGGER_CONFIG_LOG_FILE_MAX_SIZE_MB_DEFAULT,
		"Log file max size, in MB, before rotation, use 0 to disable",
	)

	loggerFlags[LOGGER_ARGS_LOG_FILE_MAX_BACKUP_NUM] = flag.Int(
		LOGGER_ARGS_LOG_FILE_MAX_BACKUP_NUM,
		LOGGER_CONFIG_LOG_FILE_MAX_BACKUP_NUM_DEFAULT,
		"How many older log files to keep upon rotation",
	)
}

func applyFlag(name string, cfg *LoggerConfig) {
	if flagPtr, ok := loggerFlags[name]; ok {
		switch name {
		case LOGGER_ARGS_USE_JSON:
			cfg.UseJson = *(flagPtr.(*bool))
		case LOGGER_ARGS_LEVEL:
			cfg.Level = *(flagPtr.(*string))
		case LOGGER_ARGS_DISABALE_SRC_FILE:
			cfg.DisableSrcFile = *(flagPtr.(*bool))
		case LOGGER_ARGS_LOG_FILE:
			cfg.LogFile = *(flagPtr.(*string))
		case LOGGER_ARGS_LOG_FILE_MAX_SIZE_MB:
			cfg.LogFileMaxSizeMB = *(flagPtr.(*int))
		case LOGGER_ARGS_LOG_FILE_MAX_BACKUP_NUM:
			cfg.LogFileMaxBackupNum = *(flagPtr.(*int))
		}
	}
}

// Apply logger args to config. If onlySet is defined, then apply only those
// that were set on the command line and not their default value.
func ApplyLoggerArgs(cfg *LoggerConfig, onlySet bool) *LoggerConfig {
	if cfg == nil {
		cfg = DefaultLoggerConfig()
	}
	if onlySet {
		flag.Visit(func(f *flag.Flag) { applyFlag(f.Name, cfg) })
	} else {
		for name := range loggerFlags {
			applyFlag(name, cfg)
		}
	}
	return cfg
}

func ApplySetLoggerArgs(cfg *LoggerConfig) {
	ApplyLoggerArgs(cfg, true)
}

func LoggerConfigFromArgs() *LoggerConfig {
	return ApplyLoggerArgs(nil, false)
}
