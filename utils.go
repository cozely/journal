package journal

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

////////////////////////////////////////////////////////////////////////////////

func (l *Logger) pre(skip int) {
	l.buf.Reset()
	l.buf.WriteString(l.prefix)
	// Timestamp
	if !notime {
		//TODO: UTC
		now := time.Now()
		if l.flags&Ldate != 0 {
			year, month, day := now.Date()
			l.writeInt(year, 4)
			l.buf.WriteByte('/')
			l.writeInt(int(month), 2)
			l.buf.WriteByte('/')
			l.writeInt(day, 2)
			l.buf.WriteByte(' ')
		}
		if l.flags&Ltime != 0 {
			hour, min, sec := now.Clock()
			l.writeInt(hour, 2)
			l.buf.WriteByte(':')
			l.writeInt(min, 2)
			l.buf.WriteByte(':')
			l.writeInt(sec, 2)
			if l.flags&Lmicroseconds != 0 {
				//TODO
			}
			l.buf.WriteByte(' ')
		}
	}
	// File and line
	if l.flags&Ldirectory != 0 {
		l.buf.WriteString(l.directory)
	} else if l.flags&Lshortfile != 0 {
		_, file, line, ok := runtime.Caller(skip)
		if ok {
			split := strings.LastIndexByte(file, '/')
			if split != -1 {
				split = strings.LastIndexByte(file[:split], '/')
			}
			file = file[split+1:]
		} else {
			file = "???"
			line = 0
		}
		fmt.Fprintf(l.buf, "%s:%d: ", file, line)
	} else if l.flags&Llongfile != 0 {
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			file = "???"
			line = 0
		}
		fmt.Fprintf(l.buf, "%s:%d: ", file, line)
	}
}

func (l *Logger) writeInt(number int, width int) {
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
		l.buf.WriteByte(b)
	}
}

func (l *Logger) post() {
	b := l.buf.Bytes()
	if b[len(b)-1] != '\n' {
		l.buf.WriteByte('\n')
	}
	l.output.Write(l.buf.Bytes())
}
