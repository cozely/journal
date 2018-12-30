package journal

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"

	"golang.org/x/sys/unix"
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
	infofd, warnfd, debugfd int
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
		infofd:  unix.Stdout,
		warnfd:  unix.Stderr,
		debugfd: unix.Stderr,
		pc:      make([]uintptr, 10),
	}
}

// InfoTo changes the file descriptor used to print informations.
func (j *Journal) InfoTo(fd int) *Journal {
	//TODO check fd
	j.infofd = fd
	return j
}

// WarnTo changes the file descriptor used to print warnings.
func (j *Journal) WarnTo(fd int) *Journal {
	//TODO check fd
	j.warnfd = fd
	return j
}

// DebugTo changes the file descriptor used to print debugging entries.
func (j *Journal) DebugTo(fd int) *Journal {
	//TODO check fd
	j.debugfd = fd
	return j
}

////////////////////////////////////////////////////////////////////////////////

// Info writes an information in the journal.
func (j *Journal) Info(format string, v ...interface{}) {
	if j.infofd != 0 {
		j.timestamp()
		write(j.buf, j.pkgname)
		fmt.Fprintf(j.buf, format, v...)
		if format == "" || format[len(format)-1] != '\n' {
			j.buf.WriteByte('\n')
		}
		unix.Write(j.infofd, j.buf.Bytes())
		j.buf.Reset()
	}
}

// Warn writes a warning in the journal.
func (j *Journal) Warn(format string, v ...interface{}) {
	if j.warnfd != 0 {
		j.timestamp()
		write(j.buf, j.pkgname)
		fmt.Fprintf(j.buf, format, v...)
		if format == "" || format[len(format)-1] != '\n' {
			j.buf.WriteByte('\n')
		}
		unix.Write(j.warnfd, j.buf.Bytes())
		j.buf.Reset()
	}
}

// Debug writes a debugging entry in the journal.
func (j *Journal) Debug(format string, v ...interface{}) {
	if j.debugfd != 0 {
		j.timestamp()
		_, file, line, ok := runtime.Caller(1)
		if ok {
			split := strings.LastIndexByte(file, '/')
			if split != -1 {
				split = strings.LastIndexByte(file[:split], '/')
			}
			fmt.Fprintf(j.buf, "%s:%d: ", file[split+1:], line)
		} else {
			write(j.buf, j.pkgname)
		}
		fmt.Fprintf(j.buf, format, v...)
		if format == "" || format[len(format)-1] != '\n' {
			j.buf.WriteByte('\n')
		}
		unix.Write(j.debugfd, j.buf.Bytes())
		j.buf.Reset()
	}
}

// Check returns true (and writes a debugging entry) if err is not nil.
func (j *Journal) Check(err error) bool {
	if err == nil {
		return false
	}
	if j.debugfd != 0 {
		j.timestamp()
		_, file, line, ok := runtime.Caller(1)
		if ok {
			split := strings.LastIndexByte(file, '/')
			if split != -1 {
				split = strings.LastIndexByte(file[:split], '/')
			}
			fmt.Fprintf(j.buf, "%s:%d: ", file[split+1:], line)
		} else {
			write(j.buf, j.pkgname)
		}
		write(j.buf, err.Error())
		j.buf.WriteByte('\n')
		unix.Write(j.debugfd, j.buf.Bytes())
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
	itoa(j.buf, year, 4)
	j.buf.WriteByte('/')
	itoa(j.buf, int(month), 2)
	j.buf.WriteByte('/')
	itoa(j.buf, day, 2)
	j.buf.WriteByte(' ')
	itoa(j.buf, hour, 2)
	j.buf.WriteByte(':')
	itoa(j.buf, min, 2)
	j.buf.WriteByte(':')
	itoa(j.buf, sec, 2)
	j.buf.WriteByte(' ')
}

func write(w io.ByteWriter, s string) {
	for i := 0; i < len(s); i++ {
		w.WriteByte(s[i])
	}
}

func writeln(w io.ByteWriter, msg ...string) {
	var b byte
	for _, s := range msg {
		for i := 0; i < len(s); i++ {
			b = s[i]
			w.WriteByte(b)
		}
	}
	if b != '\n' {
		w.WriteByte('\n')
	}
}

func itoa(w io.ByteWriter, number int, width int) {
	// Assemble decimal in reverse order.
	var buf [20]byte
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
		w.WriteByte(b)
	}
}
