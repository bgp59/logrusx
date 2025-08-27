// A logrus extension w/ file support and src file formatter

package logrusx

import (
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"

	"gopkg.in/natefinch/lumberjack.v2"

	logrusx_internal "github.com/bgp59/logrusx/internal"
)

const (
	LOGGER_CONFIG_USE_JSON_DEFAULT                = true
	LOGGER_CONFIG_LEVEL_DEFAULT                   = "info"
	LOGGER_CONFIG_DISBALE_SRC_FILE_DEFAULT        = false
	LOGGER_CONFIG_LOG_FILE_DEFAULT                = "" // i.e. stderr
	LOGGER_CONFIG_LOG_FILE_MAX_SIZE_MB_DEFAULT    = 10
	LOGGER_CONFIG_LOG_FILE_MAX_BACKUP_NUM_DEFAULT = 1

	LOGGER_ARGS_USE_JSON                = "log-use-json"
	LOGGER_ARGS_LEVEL                   = "log-level"
	LOGGER_ARGS_DISABALE_SRC_FILE       = "log-disable-src-file"
	LOGGER_ARGS_LOG_FILE                = "log-file"
	LOGGER_ARGS_LOG_FILE_MAX_SIZE_MB    = "log-file-max-size-mb"
	LOGGER_ARGS_LOG_FILE_MAX_BACKUP_NUM = "log-file-max-backup-num"

	LOGGER_DEFAULT_LEVEL = logrus.InfoLevel
)

// Collectable logger interface for logurs.Log (see testutils/log_collector.go):
type CollectableLogger struct {
	// The actual logger:
	logrus.Logger

	// Cache the condition of being enabled for debug or not. Various sections
	// of  the code may test this condition before doing more expensive actions,
	// such as formatting debug info, so it pays off to make it as efficient as
	// possible:
	IsEnabledForDebug bool

	// Caller prettyfier:
	prettyfier *logrusx_internal.CallerPrettyfier
}

func (logger *CollectableLogger) GetOutput() io.Writer {
	return logger.Out
}

func (logger *CollectableLogger) GetLevel() any {
	return logger.Logger.GetLevel()
}

func (logger *CollectableLogger) SetLevel(level any) {
	if level, ok := level.(logrus.Level); ok {
		logger.Logger.SetLevel(level)
		logger.IsEnabledForDebug = logger.IsLevelEnabled(logrus.DebugLevel)
	}
}

type LoggerConfig struct {
	// Whether to structure the logged record in JSON:
	UseJson bool `yaml:"use_json"`
	// Log level name: info, warn, ...:
	Level string `yaml:"level"`
	// Whether to disable the reporting of the source file:line# info:
	DisableSrcFile bool `yaml:"disable_src_file"`
	// Whether to log to a file or, if empty, to stderr:
	LogFile string `yaml:"log_file"`
	// Log file max size, in MB, before rotation, use 0 to disable:
	LogFileMaxSizeMB int `yaml:"log_file_max_size_mb"`
	// How many older log files to keep upon rotation:
	LogFileMaxBackupNum int `yaml:"log_file_max_backup_num"`
}

func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		UseJson:             LOGGER_CONFIG_USE_JSON_DEFAULT,
		Level:               LOGGER_CONFIG_LEVEL_DEFAULT,
		DisableSrcFile:      LOGGER_CONFIG_DISBALE_SRC_FILE_DEFAULT,
		LogFile:             LOGGER_CONFIG_LOG_FILE_DEFAULT,
		LogFileMaxSizeMB:    LOGGER_CONFIG_LOG_FILE_MAX_SIZE_MB_DEFAULT,
		LogFileMaxBackupNum: LOGGER_CONFIG_LOG_FILE_MAX_BACKUP_NUM_DEFAULT,
	}
}

func NewCollectableLogger() *CollectableLogger {
	prettyfier := logrusx_internal.NewCallerPrettyfier()
	return &CollectableLogger{
		Logger: logrus.Logger{
			Out: os.Stderr,
			//Hooks:        make(logrus.LevelHooks),
			Formatter:    logrusx_internal.NewTextFormatter(prettyfier),
			Level:        LOGGER_DEFAULT_LEVEL,
			ReportCaller: true,
		},
		prettyfier: prettyfier,
	}
}

// Set the logger based on config, post creation. This may be necessary since an
// app may start with the default logger and later, after loading loading a
// configuration and/or parsing the command line args, it may need to amend the logger.
func (logger *CollectableLogger) SetLogger(cfg *LoggerConfig) error {
	if cfg == nil {
		cfg = DefaultLoggerConfig()
	}

	levelName := cfg.Level
	if levelName != "" {
		level, err := logrus.ParseLevel(levelName)
		if err != nil {
			return err
		}
		logger.SetLevel(level)
	}

	if cfg.UseJson {
		logger.SetFormatter(logrusx_internal.NewJsonFormatter(logger.prettyfier))
	} else {
		logger.SetFormatter(logrusx_internal.NewTextFormatter(logger.prettyfier))
	}

	logger.SetReportCaller(!cfg.DisableSrcFile)

	switch logFile := cfg.LogFile; logFile {
	case "stderr":
		logger.SetOutput(os.Stderr)
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "":
	default:
		// Create log dir as needed:
		logDir := path.Dir(cfg.LogFile)
		_, err := os.Stat(logDir)
		if err != nil {
			err = os.MkdirAll(logDir, os.ModePerm)
			if err != nil {
				return err
			}
		}
		// Check if the log file exists, in which case force rotate it before
		// the 1st use:
		_, err = os.Stat(cfg.LogFile)
		forceRotate := err == nil
		logFile := &lumberjack.Logger{
			Filename:   cfg.LogFile,
			MaxSize:    cfg.LogFileMaxSizeMB,
			MaxBackups: cfg.LogFileMaxBackupNum,
		}
		if forceRotate {
			err := logFile.Rotate()
			if err != nil {
				return err
			}
		}
		logger.SetOutput(logFile)
	}

	return nil
}

func (logger *CollectableLogger) NewCompLogger(compName string) *logrus.Entry {
	return logger.WithField(logrusx_internal.LOGGER_COMPONENT_FIELD_NAME, compName)
}

// Add the prefix based on the caller's stack, going back `upNDirs` directories
// using the caller's file path. The prefix is added to the list of prefixes to
// be stripped from the file path when logging.
func (logger *CollectableLogger) AddCallerSrcPathPrefix(upNDirs int) error {
	return logger.prettyfier.AddCallerSrcPathPrefix(upNDirs, 1)
}

// Set how many sub-dirs to keep, starting from the filename towards the root,
// in case there is no prefix match. (the fallback, that is). For instance if
// the caller's path is /a/b/c/f.go and n == 2, the source will be logged as
// b/c/f.go. The builtin default is 1.
func (logger *CollectableLogger) SetKeepNDirs(n int) {
	logger.prettyfier.SetKeepNDirs(n)
}

// Get the list of supported level names:
func GetLogLevelNames() []string {
	levelNames := make([]string, len(logrus.AllLevels))
	for i, level := range logrus.AllLevels {
		levelNames[i] = level.String()
	}
	return levelNames
}
