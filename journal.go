package journal

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"
)

////////////////////////////////////////////////////////////////////////////////

var (
	appname string //TODO
	notime  bool
)

// NoTimestamp removes the timestamp in all journals.
func NoTimestamp() {
	notime = true
}

////////////////////////////////////////////////////////////////////////////////

// Journal is used to write informations, warnings or debugging entries. Each
// package should have its own Journal, created with New.
type Journal struct {
	pkgname                 string
	buf                     *bytes.Buffer
	infofd, warnfd, debugfd io.Writer
	pc                      []uintptr
}

// New returns a journal for the package of the caller.
func New() *Journal {
	_, file, _, ok := runtime.Caller(1)
	var pkgname string
	if ok {
		split1 := strings.LastIndexByte(file, '/')
		split2 := split1
		if split1 != -1 {
			split2 = strings.LastIndexByte(file[:split1], '/')
		}
		pkgname = fmt.Sprintf("%s: ", file[split2+1:split1])
	} else {
		pkgname = "???: "
	}
	return &Journal{
		pkgname: pkgname,
		buf:     bytes.NewBuffer(make([]byte, 0, 256)),
		infofd:  Stdout,
		warnfd:  Stderr,
		debugfd: Stderr,
		pc:      make([]uintptr, 10),
	}
}

// InfoTo changes the file descriptor used to print informations.
func (j *Journal) InfoTo(o io.Writer) *Journal {
	j.infofd = o
	return j
}

// WarnTo changes the file descriptor used to print warnings.
func (j *Journal) WarnTo(o io.Writer) *Journal {
	j.warnfd = o
	return j
}

// DebugTo changes the file descriptor used to print debugging entries.
func (j *Journal) DebugTo(o io.Writer) *Journal {
	j.debugfd = o
	return j
}

////////////////////////////////////////////////////////////////////////////////

// Info writes an information entry in the journal.
func (j *Journal) Info(format string, v ...interface{}) {
	if j.infofd != nil {
		j.timestamp()
		j.buf.WriteString(j.pkgname)
		fmt.Fprintf(j.buf, format, v...)
		if format == "" || format[len(format)-1] != '\n' {
			j.buf.WriteByte('\n')
		}
		j.infofd.Write(j.buf.Bytes())
		j.buf.Reset()
	}
}

// Warn writes a warning entry in the journal.
func (j *Journal) Warn(format string, v ...interface{}) {
	if j.warnfd != nil {
		j.timestamp()
		j.buf.WriteString(j.pkgname)
		fmt.Fprintf(j.buf, format, v...)
		if format == "" || format[len(format)-1] != '\n' {
			j.buf.WriteByte('\n')
		}
		j.warnfd.Write(j.buf.Bytes())
		j.buf.Reset()
	}
}

// Debug writes a debugging entry in the journal.
func (j *Journal) Debug(format string, v ...interface{}) {
	if j.debugfd != nil {
		j.timestamp()
		_, file, line, ok := runtime.Caller(1)
		if ok {
			split := strings.LastIndexByte(file, '/')
			if split != -1 {
				split = strings.LastIndexByte(file[:split], '/')
			}
			fmt.Fprintf(j.buf, "%s:%d: ", file[split+1:], line)
		} else {
			j.buf.WriteString(j.pkgname)
		}
		fmt.Fprintf(j.buf, format, v...)
		if format == "" || format[len(format)-1] != '\n' {
			j.buf.WriteByte('\n')
		}
		j.debugfd.Write(j.buf.Bytes())
		j.buf.Reset()
	}
}

// Check returns true (and writes a debugging entry) if err is not nil.
func (j *Journal) Check(err error) bool {
	if err == nil {
		return false
	}
	if j.debugfd != nil {
		j.timestamp()
		_, file, line, ok := runtime.Caller(1)
		if ok {
			split := strings.LastIndexByte(file, '/')
			if split != -1 {
				split = strings.LastIndexByte(file[:split], '/')
			}
			fmt.Fprintf(j.buf, "%s:%d: ", file[split+1:], line)
		} else {
			j.buf.WriteString(j.pkgname)
		}
		j.buf.WriteString(err.Error())
		j.buf.WriteByte('\n')
		j.debugfd.Write(j.buf.Bytes())
		j.buf.Reset()
	}
	return true
}

func (j *Journal) timestamp() {
	if notime {
		return
	}
	now := time.Now()
	year, month, day := now.Date()
	hour, min, sec := now.Clock()
	j.writeInt(year, 4)
	j.buf.WriteByte('/')
	j.writeInt(int(month), 2)
	j.buf.WriteByte('/')
	j.writeInt(day, 2)
	j.buf.WriteByte(' ')
	j.writeInt(hour, 2)
	j.buf.WriteByte(':')
	j.writeInt(min, 2)
	j.buf.WriteByte(':')
	j.writeInt(sec, 2)
	j.buf.WriteByte(' ')
}

func (j *Journal) writeInt(number int, width int) {
	// Assemble decimal in reverse order.
	var buf [16]byte
	i := len(buf) - 1
	for number >= 10 || width > 1 {
		width--
		q := number / 10
		buf[i] = byte('0' + number - q*10)
		i--
		number = q
	}
	buf[i] = byte('0' + number)
	for _, b := range buf[i:] {
		j.buf.WriteByte(b)
	}
}
