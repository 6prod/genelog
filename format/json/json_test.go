package json

import (
	"bytes"
	"fmt"
	"time"

	"github.com/6prod/genelog"
)

func ExampleJSON() {
	buf := bytes.Buffer{}

	var context = struct {
		Date  time.Time `json:"date"`
		Level string    `json:"level"`
	}{
		Date:  time.Date(2022, 2, 1, 12, 30, 0, 0, time.UTC),
		Level: "info",
	}

	logger := genelog.New(&buf).
		WithContext(context).
		WithFormatter(JSON)

	logger.Println("mylog")

	fmt.Println(&buf)

	// Output:
	// {"context":{"date":"2022-02-01T12:30:00Z","level":"info"},"message":"mylog"}
}
