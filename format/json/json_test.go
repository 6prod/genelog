package json

import (
	"bytes"
	"fmt"
	"testing"
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
	logger.Println("mylog2")

	fmt.Println(&buf)

	// Output:
	// {"context":{"date":"2022-02-01T12:30:00Z","level":"info"},"message":"mylog"}
	// {"context":{"date":"2022-02-01T12:30:00Z","level":"info"},"message":"mylog2"}
}

type C struct {
	A string
	B string
}

func TestJSON(t *testing.T) {
	c := C{"a", "b"}
	m := "coucou"
	want := `{"context":{"A":"a","B":"b"},"message":"coucou"}`

	got, err := JSON(c, m)
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Fatalf("\nwant:\n%s\ngot:\n%s\n", want, got)
	}
}
