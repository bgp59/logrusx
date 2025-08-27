module github.com/bgp59/logrusx/example

replace github.com/bgp59/logrusx => ../

go 1.23.5

require github.com/bgp59/logrusx v0.0.0-00010101000000-000000000000

require (
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/sys v0.33.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
