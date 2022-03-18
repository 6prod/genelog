package level

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/6prod/genelog"
	"github.com/6prod/genelog/format/json"
)

type testIsActiveTestCase struct {
	Ref   Level
	Level Level
	Want  bool
}

var testIsActiveTestSuite = []testIsActiveTestCase{
	{
		Ref:   ERROR,
		Level: WARNING,
		Want:  false,
	},
	{
		Ref:   DEBUG,
		Level: WARNING,
		Want:  true,
	},
}

func TestIsActive(t *testing.T) {
	for i, tc := range testIsActiveTestSuite {
		ref := tc.Ref
		level := tc.Level
		want := tc.Want

		if got := IsActive(ref, level); got != want {
			t.Fatalf("test %d: ref: %s, level: %s, want: %t, got: %t", i, ref, level, want, got)
		}
	}
}

type exampleWithLevel struct {
	*WithLevel
}

func ExampleWithLevel() {
	buf := bytes.Buffer{}

	getContext := func(v interface{}) (exampleWithLevel, error) {
		exampleWithLevel, ok := v.(exampleWithLevel)
		if !ok {
			return exampleWithLevel, fmt.Errorf("%T: not exampleWithLevel type", v)
		}
		return exampleWithLevel, nil
	}

	context := exampleWithLevel{
		NewWithLevel(INFO),
	}

	logger := genelog.New(&buf).
		WithContext(context).
		WithFormatter(json.JSON).
		AddHook(HookLevelSkip)

	context, _ = getContext(logger.Context())
	context.WithLevel.LevelSet(ERROR)
	logger = logger.WithContext(context)

	logger.Println("mylog")

	// not displayed because min level set to info
	context, _ = getContext(logger.Context())
	context.WithLevel.LevelSet(DEBUG)
	logger = logger.WithContext(context)

	logger.Println("mylog")

	fmt.Println(&buf)

	// Output:
	// {"context":{"level":"ERR"},"message":"mylog"}
}
