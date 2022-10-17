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
        Path:  "/logs/app.log",
        Daily: true,
        Days:  7,
    })

    log := logger.New(fileLogger, logger.NewStdLogger())
    log.SetLevel(logger.DebugLevel)
    defer log.Close()

    log.Debug("debug")
    log.Info("info")
    log.Warn("warn")
    log.Error("error")
    log.Fatal("fatal")
}
```
