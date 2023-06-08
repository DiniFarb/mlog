# mlog
An minimal log library for golang.

    - No fancy coloring
    - Supports log levels
    - Supports json format
    - Supports custom format
    - Supports custom output with in memory queue
 
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
	mlog.SetCustomFormat(func(logline mlog.LogLine) string {
		return fmt.Sprintf("CUSTOM |%s|%s", logline.Level, logline.Message)
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

## SetCustomOutput
If a custom output is added, the logLine struct will be additionally written to an in memory queue. The queue will be processed in a separate go routine. If the custom output function returns `false` the logline will be retried later. If the custom output function returns `true` the logline will be removed from the queue.

You can set a custom output with:
```go
    mlog.AddCustomOutput(func(logline mlog.LogLine) bool {
        msg := strings.NewReader(mlog.ApplyFormat(logline))
        resp, err := http.Post("http://localhost:8080", "text/plain", msg)
		if err != nil {
            mlog.MLogError("Error while sending logline to server: %s", err.Error())
			return false
		}
		return resp.StatusCode == http.StatusOK
	})
```

