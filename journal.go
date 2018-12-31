package journal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////

var notime bool

// EnableTimestamps enables (or disables) the timestamps in all journals.
func EnableTimestamps(b bool) {
	notime = !b
}

////////////////////////////////////////////////////////////////////////////////

// Flags to specify the prefix of the journal
const (
	Ldate         = 1 << iota                  // the date in the local time zone: 2009/01/23
	Ltime                                      // the time in the local time zone: 01:23:23
	Lmicroseconds                              // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                                  // full file name and line number: /a/b/c/d.go:23
	Lshortfile                                 // final file name element and line number: d.go:23. overrides Llongfile
	Ldirectory                                 // directory name of the file where the Logger was created.
	LUTC                                       // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime | Ldirectory // initial values for the standard logger
)

////////////////////////////////////////////////////////////////////////////////

// Logger is used to write informations, warnings or debugging entries. Each
// package should have its own Logger, created with New.
type Logger struct {
	prefix    string
	directory string
	flags     int
	buf       *bytes.Buffer
	output    io.Writer
	pc        []uintptr
}

// New returns a journal for the package of the caller.
func New(output io.Writer, prefix string, flags ...int) *Logger {
	var dir string
	_, file, _, ok := runtime.Caller(1)
	if ok {
		split1 := strings.LastIndexByte(file, '/')
		split2 := split1
		if split1 != -1 {
			split2 = strings.LastIndexByte(file[:split1], '/')
		}
		dir = fmt.Sprintf("%s: ", file[split2+1:split1])
	} else {
		dir = "???: "
	}

	f := 0
	for i := range flags {
		f |= flags[i]
	}
	if len(flags) == 0 {
		f = LstdFlags
	}

	return &Logger{
		prefix:    prefix,
		directory: dir,
		flags:     f,
		buf:       bytes.NewBuffer(make([]byte, 0, 256)),
		output:    output,
		pc:        make([]uintptr, 10),
	}
}

////////////////////////////////////////////////////////////////////////////////

// SetFlags changes the file descriptor used to print informations.
func (l *Logger) SetFlags(flags int) {
	l.flags = flags
}

// SetPrefix changes the file descriptor used to print informations.
func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

// SetOutput changes the file descriptor used to print informations.
func (l *Logger) SetOutput(output io.Writer) {
	l.output = output
}

////////////////////////////////////////////////////////////////////////////////

// Printf writes an entry in the journal. Arguments are handled in the manner of
// fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	if l.output == nil {
		return
	}
	l.pre(2)
	fmt.Fprintf(l.buf, format, v...)
	l.post()
}

// Print writes an entry in the journal. Arguments are handled in the manner of
// fmt.Print.
func (l *Logger) Print(v ...interface{}) {
	if l.output == nil {
		return
	}
	l.pre(2)
	fmt.Fprint(l.buf, v...)
	l.post()
}

// Println writes an entry in the journal. Arguments are handled in the manner of
// fmt.Println.
func (l *Logger) Println(v ...interface{}) {
	if l.output == nil {
		return
	}
	l.pre(2)
	fmt.Fprintln(l.buf, v...)
	l.post()
}

////////////////////////////////////////////////////////////////////////////////

// Panicf is equivalent to Printf followed by a call to panic.
func (l *Logger) Panicf(format string, v ...interface{}) {
	if l.output == nil {
		return
	}
	l.pre(2)
	fmt.Fprintf(l.buf, format, v...)
	l.post()
	panic(l.buf.String())
}

// Panic is equivalent to Print followed by a call to panic.
func (l *Logger) Panic(v ...interface{}) {
	if l.output == nil {
		return
	}
	l.pre(2)
	fmt.Fprint(l.buf, v...)
	l.post()
	panic(l.buf.String())
}

// Panicln is equivalent to Println followed by a call to panic.
func (l *Logger) Panicln(v ...interface{}) {
	if l.output == nil {
		return
	}
	l.pre(2)
	fmt.Fprintln(l.buf, v...)
	l.post()
	panic(l.buf.String())
}

////////////////////////////////////////////////////////////////////////////////

// Fatalf is equivalent to Printf followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.output == nil {
		return
	}
	l.pre(2)
	fmt.Fprintf(l.buf, format, v...)
	l.post()
	os.Exit(1)
}

// Fatal is equivalent to Print followed by a call to os.Exit(1).
func (l *Logger) Fatal(v ...interface{}) {
	if l.output == nil {
		return
	}
	l.pre(2)
	fmt.Fprint(l.buf, v...)
	l.post()
	os.Exit(1)
}

// Fatalln is equivalent to Println followed by a call to os.Exit(1).
func (l *Logger) Fatalln(v ...interface{}) {
	if l.output == nil {
		return
	}
	l.pre(2)
	fmt.Fprintln(l.buf, v...)
	l.post()
	os.Exit(1)
}

////////////////////////////////////////////////////////////////////////////////

// Check is equivalent to Print(err) if err is not nil.
func (l *Logger) Check(err error) {
	if err == nil || l.output == nil {
		return
	}
	l.pre(2)
	l.buf.WriteString(err.Error())
	l.post()
}
