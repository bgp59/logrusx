package logrusx_internal

import (
	"fmt"
	"path"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	LOGGER_TIMESTAMP_FORMAT = time.RFC3339
	// Extra field added for component sub loggers:
	LOGGER_COMPONENT_FIELD_NAME = "comp"
)

// When files are logged, the file name is converted to a relative path,
// generally based on the root dir its the module. If a logger is used across
// multiple modules, each module should be able its root.
type ModuleDirPathCache struct {
	// List of prefixes to be removed from the file path when logging, sorted in
	// reverse order by length.
	prefixList []string
	// If no prefix match is found, the number of directories to keep from the
	// end of the path.
	keepNDirs int
}

func NewModuleDirPathCache() *ModuleDirPathCache {
	return &ModuleDirPathCache{
		prefixList: make([]string, 0),
		keepNDirs:  1,
	}
}

func (c *ModuleDirPathCache) addPrefix(prefix string) error {
	i := len(c.prefixList) - 1
	for i >= 0 {
		if c.prefixList[i] == prefix {
			return nil // already there
		}
		if len(c.prefixList[i]) > len(prefix) {
			break
		}
		i--
	}
	i++
	if i >= len(c.prefixList) {
		c.prefixList = append(c.prefixList, prefix)
	} else {
		c.prefixList = append(c.prefixList[:i+1], c.prefixList[i:]...)
		c.prefixList[i] = prefix
	}
	return nil
}

func (c *ModuleDirPathCache) stripPrefix(filePath string) string {
	// Check if the file name starts with any of the prefixes:
	for _, prefix := range c.prefixList {
		if strings.HasPrefix(filePath, prefix) {
			// Strip the prefix and return the rest:
			return filePath[len(prefix):]
		}
	}
	// No prefix match, keep the last `keepNDirs` directories:
	pathComp := strings.Split(filePath, "/")
	keepNComps := c.keepNDirs + 1
	if keepNComps < 1 {
		keepNComps = 1
	}
	if keepNComps < len(pathComp) {
		filePath = path.Join(pathComp[len(pathComp)-keepNComps:]...)
	}
	return filePath
}

func (c *ModuleDirPathCache) setKeepNDirs(n int) {
	c.keepNDirs = n
}

func (c *ModuleDirPathCache) addCallerSrcPathPrefix(upNDirs int, skip int) error {
	skip += 1 // skip this function
	_, file, _, ok := runtime.Caller(skip)
	if !ok {
		return fmt.Errorf("cannot determine source root: runtime.Caller(%d) failed", skip)
	}
	prefix := path.Dir(file)
	for i := 0; i < upNDirs; i++ {
		prefix = path.Dir(prefix)
	}
	// The prefix should end with a slash, so that it matches a complete path
	// from a file name starting with it (e.g. "/path/to/module/" will match
	// "/path/to/module/pkg/file.go" but not "/path/to/module2/pkg/file.go")
	if prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}
	c.addPrefix(prefix)
	return nil
}

// Maintain a cache for caller PC -> (file:line#, function) to speed up the
// formatting:
type LogFuncFilePair struct {
	function string
	file     string
}

type CallerPrettyfier struct {
	m                  *sync.Mutex
	funcFileCache      map[uintptr]*LogFuncFilePair
	moduleDirPathCache *ModuleDirPathCache
}

func NewCallerPrettyfier() *CallerPrettyfier {
	return &CallerPrettyfier{
		m:                  &sync.Mutex{},
		funcFileCache:      make(map[uintptr]*LogFuncFilePair),
		moduleDirPathCache: NewModuleDirPathCache(),
	}
}

// Return the function name and filename:line# info from the frame. The filename is
// relative to the source root dir.
func (p *CallerPrettyfier) Pretiffy(f *runtime.Frame) (function string, file string) {
	p.m.Lock()
	defer p.m.Unlock()
	funcFile := p.funcFileCache[f.PC]
	if funcFile == nil {
		funcFile = &LogFuncFilePair{
			"", //f.Function,
			fmt.Sprintf("%s:%d", p.moduleDirPathCache.stripPrefix(f.File), f.Line),
		}
		p.funcFileCache[f.PC] = funcFile
	}
	return funcFile.function, funcFile.file
}

func (p *CallerPrettyfier) AddCallerSrcPathPrefix(upNDirs int, skip int) error {
	p.m.Lock()
	defer p.m.Unlock()
	return p.moduleDirPathCache.addCallerSrcPathPrefix(upNDirs, skip+1)
}

func (p *CallerPrettyfier) SetKeepNDirs(n int) {
	p.m.Lock()
	defer p.m.Unlock()
	p.moduleDirPathCache.setKeepNDirs(n)
}

var LogFieldKeySortOrder = map[string]int{
	// The desired order is time, level, file, func, other fields sorted
	// alphabetically and msg. Use negative numbers for the fields preceding
	// `other' to capitalize on the fact that any of the latter will return 0 at
	// lookup.
	logrus.FieldKeyTime:         -5,
	logrus.FieldKeyLevel:        -4,
	LOGGER_COMPONENT_FIELD_NAME: -3,
	logrus.FieldKeyFile:         -2,
	logrus.FieldKeyFunc:         -1,
	logrus.FieldKeyMsg:          1,
}

type LogFieldKeySortable struct {
	keys []string
}

func (d *LogFieldKeySortable) Len() int {
	return len(d.keys)
}

func (d *LogFieldKeySortable) Less(i, j int) bool {
	key_i, key_j := d.keys[i], d.keys[j]
	order_i, order_j := LogFieldKeySortOrder[key_i], LogFieldKeySortOrder[key_j]
	if order_i != 0 || order_j != 0 {
		return order_i < order_j
	}
	return strings.Compare(key_i, key_j) == -1
}

func (d *LogFieldKeySortable) Swap(i, j int) {
	d.keys[i], d.keys[j] = d.keys[j], d.keys[i]
}

func LogSortFieldKeys(keys []string) {
	sort.Sort(&LogFieldKeySortable{keys})
}

func NewTextFormatter(pretyffier *CallerPrettyfier) *logrus.TextFormatter {
	return &logrus.TextFormatter{
		DisableColors:    true,
		DisableQuote:     false,
		FullTimestamp:    true,
		TimestampFormat:  LOGGER_TIMESTAMP_FORMAT,
		CallerPrettyfier: pretyffier.Pretiffy,
		DisableSorting:   false,
		SortingFunc:      LogSortFieldKeys,
	}
}

func NewJsonFormatter(pretyffier *CallerPrettyfier) *logrus.JSONFormatter {
	return &logrus.JSONFormatter{
		TimestampFormat:  LOGGER_TIMESTAMP_FORMAT,
		CallerPrettyfier: pretyffier.Pretiffy,
	}
}
