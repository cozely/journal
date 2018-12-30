package journal

import (
	"fmt"
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

// Info writes an information entry in the journal.
func (j *Journal) Info(format string, v ...interface{}) {
	if j.infofd != 0 {
		j.timestamp()
		j.buf.WriteString(j.pkgname)
		fmt.Fprintf(j.buf, format, v...)
		if format == "" || format[len(format)-1] != '\n' {
			j.buf.WriteByte('\n')
		}
		j.commitInfo()
		j.buf.Reset()
	}
}

// Warn writes a warning entry in the journal.
func (j *Journal) Warn(format string, v ...interface{}) {
	if j.warnfd != 0 {
		j.timestamp()
		j.buf.WriteString(j.pkgname)
		fmt.Fprintf(j.buf, format, v...)
		if format == "" || format[len(format)-1] != '\n' {
			j.buf.WriteByte('\n')
		}
		j.commitWarn()
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
			j.buf.WriteString(j.pkgname)
		}
		fmt.Fprintf(j.buf, format, v...)
		if format == "" || format[len(format)-1] != '\n' {
			j.buf.WriteByte('\n')
		}
		j.commitDebug()
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
			j.buf.WriteString(j.pkgname)
		}
		j.buf.WriteString(err.Error())
		j.buf.WriteByte('\n')
		j.commitDebug()
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
