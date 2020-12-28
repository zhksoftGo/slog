# slog

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gookit/slog?style=flat-square)
[![GoDoc](https://godoc.org/github.com/gookit/slog?status.svg)](https://pkg.go.dev/github.com/gookit/slog)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/slog)](https://goreportcard.com/report/github.com/gookit/slog)
[![Unit-Tests](https://github.com/gookit/slog/workflows/Unit-Tests/badge.svg)](https://github.com/gookit/slog/actions)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/gookit/slog)](https://github.com/gookit/slog)

A simple log library for Go.

> Inspired the projects [Seldaek/monolog](https://github.com/Seldaek/monolog) and [sirupsen/logrus](https://github.com/sirupsen/logrus). Thank you very much

## [中文说明](README.zh-CN.md)

中文说明请阅读 [README.zh-CN](README.zh-CN.md)

## Features

- Simple, directly available without configuration
- Multiple `Handler` log handlers can be added at the same time to output logs to different places
- You can arbitrarily extend the `Handler` `Formatter` you need
- Support to custom log messages `Handler`
- Support to custom log message `Formatter`
- Built-in commonly used log write handlers
  - `console`
  - `buffer file`
  - `rotate_file`
  - `size_rotate_file`
  - `time_rotate_file` 

## GoDoc

- [godoc for github](https://pkg.go.dev/github.com/gookit/slog?tab=doc)

## Install

```bash
go get github.com/gookit/slog
```

## Usage

`slog` is very simple to use and can be used without any configuration

## Quick Start

```go
package main

import (
	"github.com/gookit/slog"
)

func main() {
	slog.Info("info log message")
	slog.Warn("warning log message")
	slog.Infof("info log %s", "message")
	slog.Debugf("debug %s", "message")
}
```

**Output:**

```text
[2020/07/16 12:19:33] [application] [INFO] [main.go:7] info log message  
[2020/07/16 12:19:33] [application] [WARNING] [main.go:8] warning log message  
[2020/07/16 12:19:33] [application] [INFO] [main.go:9] info log message  
[2020/07/16 12:19:33] [application] [DEBUG] [main.go:10] debug message  
```

### Console Color

You can enable color on output logs to console. _This is default_

```go
package main

import (
	"github.com/gookit/slog"
)

func main() {
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
	})

	slog.Trace("this is a simple log message")
	slog.Debug("this is a simple log message")
	slog.Info("this is a simple log message")
	slog.Notice("this is a simple log message")
	slog.Warn("this is a simple log message")
	slog.Error("this is a simple log message")
	slog.Fatal("this is a simple log message")
}
```

**Output:**

![](_example/images/console-color-log.png)

- Change output format

Change the default logger output format.

```go
slog.GetFormatter().(*slog.TextFormatter).Template = slog.NamedTemplate
```

**Output:**

![](_example/images/console-color-log1.png)

### Use JSON Format

```go
package main

import (
	"github.com/gookit/slog"
)

func main() {
	// use JSON formatter
	slog.SetFormatter(slog.NewJSONFormatter())

	slog.Info("info log message")
	slog.Warn("warning log message")
	slog.WithData(slog.M{
		"key0": 134,
		"key1": "abc",
	}).Infof("info log %s", "message")

	r := slog.WithFields(slog.M{
		"category": "service",
		"IP": "127.0.0.1",
	})
	r.Infof("info %s", "message")
	r.Debugf("debug %s", "message")
}
```

**Output:**

```text
{"channel":"application","data":{},"datetime":"2020/07/16 13:23:33","extra":{},"level":"INFO","message":"info log message"}
{"channel":"application","data":{},"datetime":"2020/07/16 13:23:33","extra":{},"level":"WARNING","message":"warning log message"}
{"channel":"application","data":{"key0":134,"key1":"abc"},"datetime":"2020/07/16 13:23:33","extra":{},"level":"INFO","message":"info log message"}
{"IP":"127.0.0.1","category":"service","channel":"application","datetime":"2020/07/16 13:23:33","extra":{},"level":"INFO","message":"info message"}
{"IP":"127.0.0.1","category":"service","channel":"application","datetime":"2020/07/16 13:23:33","extra":{},"level":"DEBUG","message":"debug message"}
```

## Introduction

slog handle workflow(like monolog):

```text
         Processors
Logger -{
         Handlers -{ With Formatters
```

### Processor

`Processor` - Logging `Record` processor.

You can use it to perform additional operations on the record before the log record reaches the Handler for processing, such as adding fields, adding extended information, etc.

`Processor` definition:

```go
// Processor interface definition
type Processor interface {
	// Process record
	Process(record *Record)
}
```

`Processor` func wrapper definition:

```go
// ProcessorFunc wrapper definition
type ProcessorFunc func(record *Record)

// Process record
func (fn ProcessorFunc) Process(record *Record) {
	fn(record)
}
```

Here we use the built-in processor `slog.AddHostname` as an example, it can add a new field `hostname` to each log record.

```go
slog.AddProcessor(slog.AddHostname())

slog.Info("message")
```

**Output:**

```json
{"channel":"application","level":"INFO","datetime":"2020/07/17 12:01:35","hostname":"InhereMac","data":{},"extra":{},"message":"message"}
```

### Handler

`Handler` - Log processor, each log will be processed by `Handler.Handle()`, where you can send the log to the console, file, or remote server.

> You can customize any `Handler` you want, you only need to implement the `slog.Handler` interface.

```go
// Handler interface definition
type Handler interface {
	// Close handler
	io.Closer
	// Flush logs to disk
	Flush() error
	// IsHandling Checks whether the given record will be handled by this handler.
	IsHandling(level Level) bool
	// Handle a log record.
	// all records may be passed to this method, and the handler should discard
	// those that it does not want to handle.
	Handle(*Record) error
}
```

> Note: Remember to add the `Handler` to the logger instance before the log records will be processed by the `Handler`.

### Formatter

`Formatter` - Log data formatting.

It is usually set in `Handler`, which can be used to format log records, convert records into text, JSON, etc., `Handler` then writes the formatted data to the specified place.

`Formatter` definition:

```go
// Formatter interface
type Formatter interface {
	Format(record *Record) ([]byte, error)
}
```

`Formatter` function wrapper:

```go
// FormatterFunc wrapper definition
type FormatterFunc func(r *Record) ([]byte, error)

// Format an record
func (fn FormatterFunc) Format(r *Record) ([]byte, error) {
	return fn(r)
}
```

## Custom Logger

### Create New Logger

You can create a new instance of `slog.Logger`:

- Method 1：

```go
l := slog.New()
// add handlers ...
h1 := handler.NewConsoleHandler(slog.AllLevels)
l.AddHandlers(h1)
```

- Method 2：

```go
l := slog.NewWithName("myLogger")
// add handlers ...
h1 := handler.NewConsoleHandler(slog.AllLevels)
l.AddHandlers(h1)
```

- Method 3：

```go
package main

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func main() {
	l := slog.NewWithHandlers(handler.NewConsoleHandler(slog.AllLevels))
	l.Info("message")
}
```

### Create New Handler

you only need implement the `slog.Handler` interface:

```go
package mypkg

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

type MyHandler struct {
	handler.LevelsWithFormatter
}

func (h *MyHandler) Handle(r *slog.Record) error {
	// you can write log message to file or send to remote.
}
```

add handler to default logger:

```go
slog.AddHander(&MyHandler{})
```

or add to custom logger:

```go
l := slog.New()
l.AddHander(&MyHandler{})
```

### Create New Processor

you only need implement the `slog.Processor` interface:

```go
package mypkg

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

// AddHostname to record
func AddHostname() slog.Processor {
	hostname, _ := os.Hostname()

	return slog.ProcessorFunc(func(record *slog.Record) {
		record.AddField("hostname", hostname)
	})
}
```

add the processor:

```go
slog.AddProcessor(mypkg.AddHostname())
```

or:

```go
l := slog.New()
l.AddProcessor(mypkg.AddHostname())
```

### Create New Formatter

you only need implement the `slog.Formatter` interface:

```go
package mypkg

import (
	"github.com/gookit/slog"
)

type MyFormatter struct {
}

func (f *Formatter) Format(r *slog.Record) error {
	// format Record to text/JSON or other format.
}
```

add the formatter:

```go
slog.SetFormatter(&mypkg.MyFormatter{})
```

OR:

```go
l := slog.New()
h := &MyHandler{}
h.SetFormatter(&mypkg.MyFormatter{})

l.AddHander(h)
```

## Refer

- https://github.com/golang/glog
- https://github.com/sirupsen/logrus
- https://github.com/Seldaek/monolog
- https://github.com/syyongx/llog
- https://github.com/uber-go/zap
- https://github.com/rs/zerolog

## LICENSE

[MIT](LICENSE)
