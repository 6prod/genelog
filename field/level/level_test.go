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

type testLevelTestCase struct {
	Description string
	Test        func() error
}

var testLevelTestSuite = []testLevelTestCase{
	{
		Description: "convert from string for all levels",
		Test: func() error {
			want := len(LevelString)
			got := len(LevelFromString)
			if want != got {
				return fmt.Errorf("want: %d, got: %d", want, got)
			}
			return nil
		},
	},
	{
		Description: "colors for all level",
		Test: func() error {
			want := len(LevelString)
			got := len(LevelColor)
			if want != got {
				return fmt.Errorf("want: %d, got: %d", want, got)
			}
			return nil
		},
	},
}

func TestLevel(t *testing.T) {
	for _, tc := range testLevelTestSuite {
		if err := tc.Test(); err != nil {
			t.Fatalf("%s: %s", tc.Description, err)
		}
	}
}

type exampleWithLevel struct {
	*WithLevel
}

func ExampleWithLevel() {
	// output buffer
	buf := bytes.Buffer{}

	// simply converts v into exampleWithLevel
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
	context.LevelSet(ERROR)
	logger = logger.WithContext(context)

	logger.Println("mylog")

	// not displayed because min level set to info
	context, _ = getContext(logger.Context())
	context.LevelSet(DEBUG)
	logger = logger.WithContext(context)

	logger.Println("mylog")

	fmt.Println(&buf)

	// Output:
	// {"context":{"level":"error"},"message":"mylog"}
}
