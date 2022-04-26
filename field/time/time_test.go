package time

import (
	"bufio"
	"bytes"
	libjson "encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/6prod/genelog"
	"github.com/6prod/genelog/format/json"
)

type exampleWithTime struct {
	*WithTime
}

func ExampleWithTime() {
	buf := bytes.Buffer{}

	now := time.Date(2022, time.April, 4, 26, 0, 0, 0, time.UTC)
	context := exampleWithTime{
		NewWithTime(now),
	}

	logger := genelog.New(&buf).
		WithContext(context).
		WithFormatter(json.JSON).
		AddHook(HookUpdateTime)

	logger.Print("mylog1")
	logger.Print("mylog2")

	fmt.Println(&buf)
}

func TestWithType(t *testing.T) {
	buf := bytes.Buffer{}

	now := time.Date(2022, time.April, 4, 26, 0, 0, 0, time.UTC)
	context := exampleWithTime{
		NewWithTime(now),
	}

	logger := genelog.New(&buf).
		WithContext(context).
		WithFormatter(func(v interface{}, msg string) (string, error) {
			context, ok := v.(exampleWithTime)
			if !ok {
				return "", errors.New("not exampleWithTime type")
			}

			return fmt.Sprintf("%s|%s", context.Time().Format(time.RFC3339), msg), nil
		}).
		AddHook(HookUpdateTime)

	testSuite := []string{
		"mylog1",
		"mylog2",
	}

	for _, tc := range testSuite {
		input := tc
		logger.Println(input)
	}

	var previousDate time.Time
	var i int

	// test each line of the output buffer
	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "|")
		if n := len(fields); n != 2 {
			t.Fatalf("%s: expecting %d fields", line, n)
		}

		dateStr := fields[0]
		msg := fields[1]
		want := testSuite[i]

		date, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			t.Fatal(err)
		}

		if date.Before(previousDate) {
			t.Fatalf("msg: %s, date %s should not be before previous date %s", msg, date, previousDate)
		}

		if want != msg {
			t.Fatalf("want msg %s, got: %s", want, msg)
		}

		i++
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("reading standard input: %v", err)
	}
}

func TestWithType_UnmarshalJSON(t *testing.T) {
	buf := bytes.Buffer{}

	now := time.Date(2022, time.April, 4, 26, 0, 0, 0, time.UTC)
	context := exampleWithTime{
		NewWithTime(now),
	}

	logger := genelog.New(&buf).
		WithContext(context).
		WithFormatter(json.JSON)

	jsonContext := struct {
		Context exampleWithTime
		Message string
	}{
		// init pointer
		Context: exampleWithTime{
			NewWithTime(time.Time{}),
		},
	}

	want := "mylog"
	logger.Print(want)

	// test each line of the output buffer
	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		line := scanner.Bytes()
		if err := libjson.Unmarshal(line, &jsonContext); err != nil {
			t.Fatal(err)
		}

		gotTime := jsonContext.Context.Time()
		msg := jsonContext.Message

		if !gotTime.Equal(now) {
			t.Fatalf("%s: expecting time zero", gotTime)
		}

		if msg != want {
			t.Fatalf("message: want: %s, got: %s", want, msg)
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("reading standard input: %v", err)
	}
}
