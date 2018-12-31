package journal

import (
	"fmt"
	"io"
	"os"
)

////////////////////////////////////////////////////////////////////////////////

var std = New(Stderr, "", Lshortfile)

////////////////////////////////////////////////////////////////////////////////

// SetFlags changes the file descriptor used to print informations.
func SetFlags(flags int) {
	std.flags = flags
}

// SetPrefix changes the file descriptor used to print informations.
func SetPrefix(prefix string) {
	std.prefix = prefix
}

// SetOutput changes the file descriptor used to print informations.
func SetOutput(output io.Writer) {
	std.output = output
}

////////////////////////////////////////////////////////////////////////////////

// Printf writes an entry in the journal. Arguments are handled in the manner of
// fmt.Printf.
func Printf(format string, v ...interface{}) {
	if std.output == nil {
		return
	}
	std.pre(2)
	fmt.Fprintf(std.buf, format, v...)
	std.post()
}

// Print writes an entry in the journal. Arguments are handled in the manner of
// fmt.Print.
func Print(v ...interface{}) {
	if std.output == nil {
		return
	}
	std.pre(2)
	fmt.Fprint(std.buf, v...)
	std.post()
}

// Println writes an entry in the journal. Arguments are handled in the manner of
// fmt.Println.
func Println(v ...interface{}) {
	if std.output == nil {
		return
	}
	std.pre(2)
	fmt.Fprintln(std.buf, v...)
	std.post()
}

////////////////////////////////////////////////////////////////////////////////

// Panicf is equivalent to Printf followed by a call to panic.
func Panicf(format string, v ...interface{}) {
	if std.output == nil {
		return
	}
	std.pre(2)
	fmt.Fprintf(std.buf, format, v...)
	std.post()
	panic(std.buf.String())
}

// Panic is equivalent to Print followed by a call to panic.
func Panic(v ...interface{}) {
	if std.output == nil {
		return
	}
	std.pre(2)
	fmt.Fprint(std.buf, v...)
	std.post()
	panic(std.buf.String())
}

// Panicln is equivalent to Println followed by a call to panic.
func Panicln(v ...interface{}) {
	if std.output == nil {
		return
	}
	std.pre(2)
	fmt.Fprintln(std.buf, v...)
	std.post()
	panic(std.buf.String())
}

////////////////////////////////////////////////////////////////////////////////

// Fatalf is equivalent to Printf followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	if std.output == nil {
		return
	}
	std.pre(2)
	fmt.Fprintf(std.buf, format, v...)
	std.post()
	os.Exit(1)
}

// Fatal is equivalent to Print followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	if std.output == nil {
		return
	}
	std.pre(2)
	fmt.Fprint(std.buf, v...)
	std.post()
	os.Exit(1)
}

// Fatalln is equivalent to Println followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	if std.output == nil {
		return
	}
	std.pre(2)
	fmt.Fprintln(std.buf, v...)
	std.post()
	os.Exit(1)
}

////////////////////////////////////////////////////////////////////////////////

// Check is equivalent to Print(err) if err is not nil.
func Check(err error) {
	if err == nil || std.output == nil {
		return
	}
	std.pre(2)
	std.buf.WriteString(err.Error())
	std.post()
}
