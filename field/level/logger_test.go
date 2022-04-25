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
	// warning: mylog
	// error: mylog
	// fatal: mylog
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

	w1 := logger.Writer(WARNING)
	_, _ = io.WriteString(w1, "w1")
	_, _ = io.WriteString(w1, "w1")

	fmt.Print(buf.String())
	// Output:
	// {"context":{"level":"warning"},"message":"w1"}
	// {"context":{"level":"warning"},"message":"w1"}
}
