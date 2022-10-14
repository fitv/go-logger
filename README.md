## Logger for Go

## Install
```
go get -u github.com/fitv/go-logger
```

## Usage
```go
package main

import "github.com/fitv/go-logger"

func main() {
    fileLogger := logger.NewFileLogger(&logger.Option{
        Path:  "/var/log",
        Name:  "app",
        Daily: true,
        Days:  15,
    })

    log := logger.New(fileLogger, logger.NewStdLogger())
    defer log.Close()
    log.SetLevel(logger.DebugLevel)

    log.Debug("debug")
    log.Info("info")
    log.Warn("warn")
    log.Error("error")
    log.Fatal("fatal")
}
```
