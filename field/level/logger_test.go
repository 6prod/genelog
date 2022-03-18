package level

import (
	"bytes"
	"fmt"

	"github.com/6prod/genelog"
	"github.com/6prod/genelog/format/json"
)

func ExampleLevelLogger() {
	buf := bytes.Buffer{}

	context := exampleWithLevel{
		NewWithLevel(INFO),
	}

	logger := NewLevelLogger(genelog.New(&buf).
		WithContext(context).
		WithFormatter(json.JSON))

	logger.Print("mylog")
	logger.Info("mylog")
	logger.Error("mylog")
	logger.Debug("mylog")
	logger.Warning("mylog")
	logger.Fatalf("%s", "mylog")

	fmt.Print(buf.String())
	// Output:
	// {"context":{"level":"UNSET"},"message":"mylog"}
	// {"context":{"level":"INFO"},"message":"mylog"}
	// {"context":{"level":"ERR"},"message":"mylog"}
	// {"context":{"level":"WARNING"},"message":"mylog"}
	// {"context":{"level":"FATAL"},"message":"mylog"}
}
