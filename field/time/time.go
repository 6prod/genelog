package time

import (
	"encoding/json"
	"fmt"
	"time"
)

type WithTime struct {
	time time.Time
}

// withTimeJSON is an helper structure
// to decode json from WithTime with
// private attributes
type withTimeJSON struct {
	Time time.Time `json:"time"`
}

func NewWithTime() *WithTime {
	return &WithTime{}
}

func (w WithTime) Time() time.Time {
	return w.time
}

func (w *WithTime) TimeSet(t time.Time) {
	w.time = t
}

func (w WithTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(withTimeJSON{Time: w.time})
}

func (w *WithTime) UnmarshalJSON(b []byte) error {
	var withTime withTimeJSON
	if err := json.Unmarshal(b, &withTime); err != nil {
		return err
	}

	*w = WithTime{
		time: withTime.Time,
	}

	return nil
}

type Timer interface {
	Time() time.Time
	TimeSet(time.Time)
}

func HookUpdateTime(v interface{}, msg string) (interface{}, string, error) {
	context, ok := v.(Timer)
	if !ok {
		return nil, "", fmt.Errorf("%T: not implementing the Timer interface", v)
	}
	context.TimeSet(time.Now())
	return context, msg, nil
}
