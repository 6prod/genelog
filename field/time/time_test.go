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

	getContext := func(v interface{}) (exampleWithTime, error) {
		exampleWithTime, ok := v.(exampleWithTime)
		if !ok {
			return exampleWithTime, fmt.Errorf("%T: not exampleWithTime type", v)
		}
		return exampleWithTime, nil
	}

	context := exampleWithTime{
		NewWithTime(),
	}

	logger := genelog.New(&buf).
		WithContext(context).
		WithFormatter(json.JSON).
		AddHook(HookUpdateTime)

	context, _ = getContext(logger.Context())
	logger.Println("mylog1")
	logger.Println("mylog2")

	fmt.Println(&buf)
}
