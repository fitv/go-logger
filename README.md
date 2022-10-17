## Logger for Go

## Install
```
go get -u github.com/fitv/go-logger
```

## Usage
```go
package main

import (
    "github.com/fitv/go-logger"
)

func main() {
    fileWriter := logger.NewFileWriter(&logger.Option{
        Path:  "/app/web.log",
        Daily: true,
        Days:  7,
    })
    defer fileWriter.Close()

    log := logger.New()
    log.SetOut(fileWriter)
    log.SetLevel(logger.DebugLevel)

    log.Debug("debug")
    log.Info("info")
    log.Warn("warn")
    log.Error("error")
    log.Fatal("fatal")
}
```
