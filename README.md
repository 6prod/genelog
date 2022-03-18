# genelog

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/6prod/genelog) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/6prod/genelog/master/LICENSE)

Small (200 loc) and composable logger with context, hooks and formatter.

## Features
- Composable
- Support any context
- Support any formatter
- Support hook functions to update context and message on writes
- Extensions:
  - Fields:
    - Level
    - Time
  - Formatter:
    - JSON

## Usage
### Basic
```go
logger := New(&buf).
  WithFormatter(func(context interface{}, msg string) (string, error) {
    return msg, nil
  })

logger.Println("mylog1")
logger.Println("mylog2")

// Output:
// mylog1
// mylog2
```

### With context
```go
type Context struct {
  S string
  N int
}

// updateContext is an helper function to update
// this specific context
func updateContext(update func(Context) Context) Update {
  return func(v interface{}) (interface{}, error) {
    context, ok := v.(Context)
    if !ok {
      return nil, errors.New("not Context type")
    }
    return update(context), nil
  }
}

func main() {
  // our context
  context := Context{}

  // create logger
  logger := New(os.Stdout).
    WithContext(context).
    WithFormatter(func(v interface{}, msg string) (string, error) {
      context, ok := v.(Context)
      if !ok {
        return "", errors.New("not Context type")
      }
      return fmt.Sprintf("%s %d %s", context.S, context.N, msg), nil
    })

  // update context
  logger.UpdateContext(updateContext(
    func(context Context) Context {
      context.S = "A"
      context.N = 1
      return context
    }))

  // write log items
  logger.Print("mylog")

  // Output:
  // A 1 mylog
```

### With level
```go
package main

import (
  "os"

  "github.com/6prod/genelog"
  "github.com/6prod/genelog/format/json"
  "github.com/6prod/genelog/field/level"
)

type Context struct {
  *level.WithLevel
}

func main() {
  context := Context{
    NewWithLevel(INFO),
  }

  logger := NewLevelLogger(genelog.New(os.Stdout).
    WithContext(context).
    WithFormatter(json.JSON))

  logger.Print("mylog")
  logger.Info("mylog")
  logger.Error("mylog")
  logger.Debug("mylog")
  logger.Warning("mylog")
  logger.Fatalf("%s", "mylog")

  // Output:
  // {"context":{"level":"UNSET"},"message":"mylog"}
  // {"context":{"level":"INFO"},"message":"mylog"}
  // {"context":{"level":"ERR"},"message":"mylog"}
  // {"context":{"level":"WARNING"},"message":"mylog"}
  // {"context":{"level":"FATAL"},"message":"mylog"}
}
```

See documentation for more examples.

## Benchmark
```sh
$ ci/bench
goos: openbsd
goarch: amd64
pkg: github.com/6prod/genelog
cpu: Intel(R) Core(TM) i7-5600U CPU @ 2.60GHz
BenchmarkLogger-2         350896              3142 ns/op
BenchmarkGoLogger-2       276690              4432 ns/op
BenchmarkZerolog-2        591278              2223 ns/op
PASS
ok      github.com/6prod/genelog        3.845s
```
