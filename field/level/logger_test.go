package level

import (
	"bytes"
	"fmt"

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

	context = func() exampleWithLevel {
		c, _ := logger.Context().(exampleWithLevel)
		return c
	}()

	logger = logger.WithContext(context)
	logger.WithFormatter(func(v interface{}, msg string) (string, error) {
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
