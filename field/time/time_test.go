package time

import (
	"bytes"
	"fmt"

	"github.com/6prod/genelog"
	"github.com/6prod/genelog/format/json"
)

type exampleWithTime struct {
	*WithTime
}

func ExampleWithTime() {
	buf := bytes.Buffer{}

	context := exampleWithTime{
		NewWithTime(),
	}

	logger := genelog.New(&buf).
		WithContext(context).
		WithFormatter(json.JSON).
		AddHook(HookUpdateTime)

	logger.Println("mylog1")
	logger.Println("mylog2")

	fmt.Println(&buf)
}
