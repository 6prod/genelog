package genelog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func ExampleLogger_Println() {
	buf := bytes.Buffer{}

	logger := New(&buf).
		WithFormatter(func(v interface{}, msg string) (string, error) {
			return msg, nil
		})

	logger.Println("mylog1")
	logger.Println("mylog2")
	logger.Println("mylog3")
	logger.Println("mylog4")

	fmt.Print(&buf)
	// Output:
	// mylog1
	// mylog2
	// mylog3
	// mylog4
}

type Context struct {
	S string
	N int
}

func updateContext(update func(Context) Context) Update {
	return func(v interface{}) (interface{}, error) {
		context, ok := v.(Context)
		if !ok {
			return nil, errors.New("not Context type")
		}
		return update(context), nil
	}
}

func ExampleHook_context() {
	buf := bytes.Buffer{}

	context := Context{}

	logger := New(&buf).
		WithContext(context).
		WithFormatter(func(v interface{}, msg string) (string, error) {
			context, ok := v.(Context)
			if !ok {
				return "", errors.New("not Context type")
			}
			return fmt.Sprintf("%s %d %s", context.S, context.N, msg), nil
		})

	_ = logger.UpdateContext(updateContext(
		func(context Context) Context {
			context.S = "A"
			context.N = 1
			return context
		}))

	logger.Print("mylog")

	fmt.Print(buf.String())

	// Output:
	// A 1 mylog
}

func ExampleHook_contextUpdateTime() {
	var context time.Time = time.Date(2021, 2, 1, 12, 30, 0, 0, time.UTC)

	hookUpdateTime := func(v interface{}, msg string) (interface{}, string, error) {
		context := time.Date(2022, 2, 1, 12, 30, 0, 0, time.UTC)
		return context, msg, nil
	}

	formatter := func(v interface{}, msg string) (string, error) {
		context, ok := v.(time.Time)
		if !ok {
			return "", fmt.Errorf("%T: not time type", v)
		}

		return fmt.Sprintf("%s: %s", context, msg), nil
	}

	buf := bytes.Buffer{}
	logger := New(&buf).
		WithContext(context).
		WithFormatter(formatter).
		AddHook(hookUpdateTime)

	logger.Println("message")
	fmt.Println(buf.String())

	// Output:
	// 2022-02-01 12:30:00 +0000 UTC: message
}

// Write to the io.Writer of the logger with no format
// and to console using a hook
func ExampleHook_multiFormatter() {
	var context string

	// To replace with os.Stdout
	var console bytes.Buffer

	hookWriteConsole := func(v interface{}, msg string) (interface{}, string, error) {
		context, ok := v.(string)
		if !ok {
			return nil, "", fmt.Errorf("%T: not string type", v)
		}

		fmt.Fprintf(&console, "%s: %s", context, msg)
		return v, msg, nil
	}

	logger := New(io.Discard).
		WithContext(context).
		AddHook(hookWriteConsole)

	_ = logger.UpdateContext(func(v interface{}) (interface{}, error) {
		context = "mycontext"
		return context, nil
	})

	logger.Println("message")
	fmt.Println(console.String())

	// Output:
	// mycontext: message
}

type benchmarkContext struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
}

func getbenchmarkContext(v interface{}) (benchmarkContext, error) {
	context, ok := v.(benchmarkContext)
	if !ok {
		return context, fmt.Errorf("%T: no context type", v)
	}
	return context, nil
}

func BenchmarkLogger(b *testing.B) {
	context := benchmarkContext{
		Level: "info",
	}

	//logger := New(os.Stdout).
	logger := New(io.Discard).
		WithContext(context).
		WithFormatter(func(v interface{}, msg string) (string, error) {
			context, err := getbenchmarkContext(v)
			if err != nil {
				return "", err
			}
			context.Message = msg
			out, err := json.Marshal(context)
			return string(out), err
		}).
		AddHook(func(v interface{}, msg string) (interface{}, string, error) {
			context, err := getbenchmarkContext(v)
			if err != nil {
				return nil, "", err
			}
			context.Time = time.Now()
			return context, msg, nil
		})

	for i := 0; i < b.N; i++ {
		logger.Println(i)
	}
	//b.FailNow()
}

func BenchmarkGoLogger(b *testing.B) {
	//logger := log.New(os.Stdout, "", 0)
	logger := log.New(io.Discard, "", 0)

	context := benchmarkContext{
		Level: "info",
	}

	for i := 0; i < b.N; i++ {
		context.Time = time.Now()
		context.Message = fmt.Sprintf("%d", i)
		out, err := json.Marshal(context)
		if err != nil {
			b.Fatal(err)
		}
		logger.Println(string(out))
	}
	//b.FailNow()
}

func BenchmarkZerolog(b *testing.B) {
	//logger := zerolog.New(os.Stdout).
	logger := zerolog.New(io.Discard).
		Level(zerolog.InfoLevel).
		With().Timestamp().
		Logger()

	for i := 0; i < b.N; i++ {
		logger.Info().Msgf("%d\n", i)
	}
	//b.FailNow()
}

func TestLogger_WithContext(t *testing.T) {
	buf := bytes.Buffer{}

	logger := New(&buf).
		WithFormatter(func(v interface{}, msg string) (string, error) {
			return fmt.Sprintf("%v: %s", v, msg), nil
		})

	logger1 := logger.WithContext(3)
	logger1.Println("logger 1 msg1")

	logger2 := logger.WithContext("string")
	logger2.Println("logger2 msg1")

	logger1.Println("logger 1 msg2")

	want := `3: logger 1 msg1
string: logger2 msg1
3: logger 1 msg2
`
	if got := buf.String(); got != want {
		t.Fatalf("\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestLogger_Writer(t *testing.T) {
	buf := bytes.Buffer{}

	prefix := "prefix "
	input := "mylog1"

	logger := New(&buf).
		WithFormatter(func(v interface{}, msg string) (string, error) {
			return fmt.Sprintf("%s%s", prefix, msg), nil
		})

	if n, err := logger.Write([]byte(input)); err != nil {
		t.Fatal(err)
	} else if want := len(input); want != n {
		t.Fatalf("want=%d, got=%d", want, n)
	}

	if n, err := io.Copy(logger, bytes.NewBufferString(input)); err != nil {
		t.Fatal(err)
	} else if want := len(input); int64(want) != n {
		t.Fatalf("want=%d, got=%d", want, n)
	}
}
