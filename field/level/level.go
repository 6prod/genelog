package level

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/6prod/genelog"
	"github.com/fatih/color"
)

// Type list the available type of logger
type Level int

func (l Level) Color() string {
	return LevelColor[l]
}

func (l Level) String() string {
	return LevelString[l]
}

func (l Level) MarshalText() (text []byte, err error) {
	return []byte(l.String()), nil
}

func (l *Level) UnmarshalText(text []byte) error {
	level, ok := NewLevelFromString(string(text))
	if !ok {
		return fmt.Errorf("%s: unknown level", string(text))
	}
	*l = level
	return nil
}

const (
	UNSET Level = iota
	// DEBUG is the type to log DEBUG message
	DEBUG
	// INFO is the type to log INFO message
	INFO
	// WARNING is the type to log WARNING message
	WARNING
	// ERROR is the type to log ERROR message
	ERROR
	// FATAL is the type to log FATAL message
	FATAL
	// OFF turns off logging
	OFF
)

// Colors
var (
	// ErrorColor defines the color of the ERROR label
	ErrorColor = color.New(color.Bold, color.FgHiRed)
	// InfoColor defines the color of the INFO label
	InfoColor = color.New(color.Bold, color.FgHiMagenta)
	// DebugColor defines the color of the DEBUG label
	DebugColor = color.New(color.Bold, color.FgHiBlue)
	// WarningColor defines the color of the DEBUG label
	WarningColor = color.New(color.Bold, color.FgHiYellow)
	// FatalColor defines the color of the DEBUG label
	FatalColor = color.New(color.Bold, color.FgHiRed)
)

var (
	ErrLevelerNotImplemented = errors.New("Leveler interface not implemented")
)

var LevelString = map[Level]string{
	UNSET:   "unset",
	DEBUG:   "debug",
	INFO:    "info",
	WARNING: "warning",
	ERROR:   "error",
	FATAL:   "fatal",
	OFF:     "off",
}

var LevelFromString = map[string]Level{
	"unset":   UNSET,
	"debug":   DEBUG,
	"info":    INFO,
	"warning": WARNING,
	"error":   ERROR,
	"fatal":   FATAL,
	"off":     OFF,
}

var LevelColor = map[Level]string{
	UNSET:   LevelString[UNSET],
	DEBUG:   DebugColor.Sprint(LevelString[DEBUG]),
	INFO:    InfoColor.Sprint(LevelString[INFO]),
	WARNING: WarningColor.Sprint(LevelString[WARNING]),
	ERROR:   ErrorColor.Sprint(LevelString[ERROR]),
	FATAL:   FatalColor.Sprint(LevelString[FATAL]),
	OFF:     LevelString[OFF],
}

func NewLevelFromString(s string) (Level, bool) {
	s = strings.ToLower(s)
	v, ok := LevelFromString[s]
	return v, ok
}

// IsActive returns true if l includes ref
func IsActive(ref, l Level) bool {
	return l >= ref
}

type WithLevel struct {
	levelMin Level
	level    Level
}

func NewWithLevel(min Level) *WithLevel {
	return &WithLevel{
		levelMin: min,
	}
}

func (w WithLevel) Level() Level {
	return w.level
}

func (w WithLevel) LevelMin() Level {
	return w.levelMin
}

func (w *WithLevel) LevelSet(level Level) {
	w.level = level
}

func (w WithLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Level Level `json:"level"`
	}{
		Level: w.level,
	})
}

// Leveler is the interface to access the level field
type Leveler interface {
	// LevelMin returns the minimum level to print
	LevelMin() Level
	// Level returns the current level
	Level() Level
	// LevelSet changes the level
	LevelSet(Level)
}

// GetLeveler converts v into Leveler.
// Returns false if not possible.
func GetLeveler(v interface{}) (Leveler, bool) {
	leveler, ok := v.(Leveler)
	if !ok {
		return nil, false
	}
	return leveler, true
}

func HookLevelSkip(v interface{}, msg string) (interface{}, string, error) {
	context, ok := v.(Leveler)
	if !ok {
		return nil, "", fmt.Errorf("%T: not implementing the Leveler interface", v)
	}
	if !IsActive(context.LevelMin(), context.Level()) {
		return v, msg, genelog.ErrSkip
	}
	return v, msg, nil
}
