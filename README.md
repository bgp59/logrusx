# A [logrus](https://pkg.go.dev/github.com/sirupsen/logrus) Extension Module

`logrusx` provides the following additional features:

* file logging via [lumberjack](https://pkg.go.dev/gopkg.in/natefinch/lumberjack.v2)

* source file path logging relative to the module root. See [internal](internal)

* YAML loadable configuration

* command line loadable configuration

* support for testing whereby the log output is collected and it is displayed via testing.T.Log at the end, only in case of error or enabled verbosity. See [testutils](testutils)

Although anyone is welcome to use it, this module is not intended for public consumption, hence the lack of polished documentation. See [example](example) in lieu of reference documentation.
