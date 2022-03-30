package level

import (
	"bytes"
	"fmt"
	"io"

	"github.com/6prod/genelog/format/json"
)

func ExampleLevelLogger() {
	buf := bytes.Buffer{}

	context := exampleWithLevel{
		NewWithLevel(WARNING),
	}

	logger := NewLevelLogger(&buf).
		WithContext(context).
		WithFormatter(json.JSON)

	logger = logger.WithFormatter(func(v interface{}, msg string) (string, error) {
		context, _ := v.(exampleWithLevel)
		return fmt.Sprintf("%s: %s\n", context.Level(), msg), nil
	})

	logger.AddHook(func(v interface{}, msg string) (interface{}, string, error) {
		return v, msg, nil
	})

	logger.Print("mylog")
	logger.Debug("mylog")
	logger.Info("mylog")
	logger.Warning("mylog")
	logger.Error("mylog")
	logger.Fatalf("%s", "mylog")

	fmt.Print(buf.String())
	// Output:
	// WARNING: mylog
	// ERR: mylog
	// FATAL: mylog
}

func ExampleLevelLogger_Writer() {
	buf := bytes.Buffer{}

	var context = struct {
		*WithLevel
	}{
		NewWithLevel(INFO),
	}

	logger := NewLevelLogger(&buf).
		WithContext(context).
		WithFormatter(json.JSON)

	w := logger.Writer(WARNING)
	_, _ = io.WriteString(w, "mylog1")
	_, _ = io.WriteString(w, "mylog2")

	w = logger.Writer(DEBUG)
	if _, err := io.WriteString(w, "mylog3"); err != nil {
		buf.WriteString(err.Error())
	}
	_, _ = io.WriteString(w, "mylog4")

	fmt.Print(buf.String())
	// Output:
	// {"context":{"level":"WARNING"},"message":"mylog1"}
	// {"context":{"level":"WARNING"},"message":"mylog2"}
}
