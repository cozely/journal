package journal

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/sys/unix"
)

////////////////////////////////////////////////////////////////////////////////

var (
	buf     = bytes.NewBuffer(make([]byte, 0, 256))
	printfd = unix.Stdout
	warnfd  = unix.Stderr
	debugfd = unix.Stderr
)

////////////////////////////////////////////////////////////////////////////////

func PrintTo(fd int) {
	//TODO check fd
	printfd = fd
}

func WarnTo(fd int) {
	//TODO check fd
	warnfd = fd
}

func DebugTo(fd int) {
	//TODO check fd
	debugfd = fd
}

////////////////////////////////////////////////////////////////////////////////

type Journal struct {
	prefix string
	info   bool
}

func New(prefix string, info bool) Journal {
	return Journal{
		prefix: prefix,
		info:   info,
	}
}

func (j Journal) Put(msg string) {
	if j.info && printfd !=0 {
		write(buf, j.prefix)
		write(buf, msg)
		if msg == "" || msg[len(msg)-1] != '\n' {
			buf.WriteByte('\n')
		}
		unix.Write(printfd, buf.Bytes())
		buf.Reset()
	}
}

func (j Journal) Print(msg ...string) {
	if j.info && printfd !=0 {
		write(buf, j.prefix)
		writeln(buf, msg...)
		unix.Write(printfd, buf.Bytes())
		buf.Reset()
	}
}

func (j Journal) Printf(format string, v ...interface{}) {
	if j.info {
		write(buf, j.prefix)
		fmt.Fprintf(buf, format, v...)
		if format == "" || format[len(format)-1] != '\n' {
			buf.WriteByte('\n')
		}
		unix.Write(printfd, buf.Bytes())
		buf.Reset()
	}
}

func (j Journal) Warn(msg ...string) {
	if warnfd != 0 {
		write(buf, j.prefix)
		writeln(buf, msg...)
		unix.Write(warnfd, buf.Bytes())
		buf.Reset()
	}
}

func (j Journal) Debug(msg ...string) {
	if debugfd != 0 {
		write(buf, j.prefix)
		writeln(buf, msg...)
		unix.Write(debugfd, buf.Bytes())
		buf.Reset()
	}
}

func (j Journal) Check(err error) bool {
	if err == nil {
		return false
	}
	if debugfd != 0 {
		write(buf, j.prefix)
		writeln(buf, err.Error())
		unix.Write(debugfd, buf.Bytes())
		buf.Reset()
	}
	return true
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
