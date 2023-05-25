# mlog
An minimal log library for golang. It's simple to use and only meant for stdout logging.

    - No fancy coloring
    - Supports log levels
    - Supports json to stdout
    - Supports custom log format
 
## Install
```
go get github.com/dinifarb/mlog
```

## Default Logger Usage

### Text
```
[MLOG] <timestamp> | <level> | <message>
```

```go
package main

import "github.com/dinifarb/mlog"

func main(){
    mlog.Trace("Hello World, Trace!")	
    mlog.Debug("Hello World, Debug!")
    mlog.Info("Hello World, Info!")
    mlog.Warn("Hello World, Warn!")
    mlog.Error("Hello World, Error!")
}
```

### Json
```json
{"timestamp":"<timestamp>","appname":"MLOG","level":"<level>","message":"<message>"}
```

```go
package main

import "github.com/dinifarb/mlog"

func main(){
    mlog.SetFormat(mlog.Fjson)
    mlog.Trace("Hello World, Trace!")	
    mlog.Debug("Hello World, Debug!")
    mlog.Info("Hello World, Info!")
    mlog.Warn("Hello World, Warn!")
    mlog.Error("Hello World, Error!")
}
```

## Settings

### Set level

Set the level with:
```go
mlog.SetLevel(mlog.Ltrace)
```
Supported Levels:
```
Ltrace
Ldebug
Linfo
Lwarn
Lerror
```

### Set app name
You can set the app name with:
```go
mlog.SetAppName("MyApp")
```
This will change the beginning of the logline like:
```
[MyApp] <timestamp> | <level> | <message>
```

### Set time format
You can set the time format with:
```go
mlog.SetTimeFormat("2006-01-02 15:04:05")
```
golang time format: https://golang.org/pkg/time/#pkg-constants

### Set format
There are two formats available:
```
Ftext
Fjson
```
You can set the format with:
```go
mlog.SetFormat(mlog.Fjson)
```

### Set custom format
You can set a custom format with:
```go
	mlog.SetCustomFormat(func(l mlog.LogLine) {
		fmt.Println("CUSTOM:", l.Level, l.Message)
	})
```
The `LogLine` struct looks like:
```go
type LogLine struct {
    Timestamp string
    AppName   string
    Level     string
    Message   string
}
```