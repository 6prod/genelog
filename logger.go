package genelog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"
)

var (
	// ErrLogger is the general logger error
	ErrLogger = errors.New("logger error")
	// ErrSkip tells the logger to skip a log entry
	ErrSkip = errors.New("skip")
)

// Format assembles context and msg into a single string
type Format func(context interface{}, msg string) (out string, err error)

// Update returns a new updated context
type Update func(context interface{}) (newcontext interface{}, err error)

// Hook is called before each write to update the context or message
//
// If err is ErrSkip, just return without writing anything
type Hook func(context interface{}, msg string) (newcontext interface{}, newmsg string, err error)

type Logger struct {
	// mu synchronizes writes
	mu *sync.Mutex
	// w writes the logs somewhere
	w io.Writer
	// context adds metadata to logs
	context interface{}
	// formatter is a function to shape the log output
	format Format
	// hooks updates the context and message on every writes
	hooks []Hook
}

func New(w io.Writer) *Logger {
	return &Logger{
		mu:    new(sync.Mutex),
		w:     w,
		hooks: make([]Hook, 0),
	}
}

// Print uses fmt.Print to write to the logger
func (l *Logger) Print(v ...interface{}) {
	msg := fmt.Sprint(v...)
	l.write(msg, func(w io.Writer, s string) {
		fmt.Fprint(w, s)
	})
}

// Println uses fmt.Println to write to the logger
func (l *Logger) Println(v ...interface{}) {
	msg := fmt.Sprint(v...)
	l.write(msg, func(w io.Writer, s string) {
		fmt.Fprintln(w, s)
	})
}

// Printf uses fmt.Printf to write to the logger
func (l *Logger) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.write(msg, func(w io.Writer, s string) {
		fmt.Fprint(w, s)
	})
}

// WithContext adds a context to the logger
func (l Logger) WithContext(v interface{}) *Logger {
	logger := l
	logger.context = v
	return &logger
}

// Context returns the context
func (l Logger) Context() interface{} {
	return l.context
}

// WithFormatter adds a formatter function to the logger
func (l Logger) WithFormatter(f Format) *Logger {
	logger := l
	logger.format = f
	return &logger
}

// AddHook adds a hook function to the list of hooks of the logger.
//
// Hooks are called in the added order
func (l Logger) AddHook(h Hook) *Logger {
	logger := l
	logger.hooks = append(logger.hooks, h)
	return &logger
}

// UpdateContext updates the logger context with the update function
func (l *Logger) UpdateContext(update Update) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	context, err := update(l.context)
	if err != nil {
		return err
	}
	l.context = context
	return nil
}

// Write writes into the logger
func (l *Logger) Write(p []byte) (n int, err error) {
	msg := bytes.NewBuffer(p).String()
	l.write(msg, func(w io.Writer, s string) {
		n, err = fmt.Fprint(w, s)
	})
	return
}

func (l *Logger) write(msg string, fn func(w io.Writer, s string)) {
	var err error
	l.mu.Lock()
	defer l.mu.Unlock()

	// Apply hook if any
	for _, hook := range l.hooks {
		l.context, msg, err = hook(l.context, msg)
		if err != nil {
			if errors.Is(err, ErrSkip) {
				return
			}
			fmt.Fprintf(l.w, "%s: %s\n", ErrLogger, err)
			return
		}
	}

	// Apply formatter if any
	if l.format != nil {
		msg, err = l.format(l.context, msg)
		if err != nil {
			fmt.Fprintf(l.w, "%s: %s\n", ErrLogger, err)
			return
		}
	}

	// Process output
	fn(l.w, msg)
}
