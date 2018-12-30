package journal

import (
	"bytes"
	"io"
	"os"
)

////////////////////////////////////////////////////////////////////////////////

type Journal struct {
	prefix string
	buf *bytes.Buffer
	info, warn, debug io.Writer
}

func New(prefix string, info, warn, debug io.Writer) *Journal {
	return &Journal{
		prefix: prefix,
		buf: bytes.NewBuffer(make([]byte, 0, 256)),
		info: info,
		warn: warn,
		debug: debug,
	}
}

func (j Journal) Info(msg ...string) {
	if j.info != nil {
		write(j.buf, j.prefix)
		writeln(j.buf, msg...)
		os.Stdout.Write(j.buf.Bytes())
		j.buf.Reset()
	}
}

func (j Journal) Warn(msg ...string) {
	if j.warn != nil {
		write(j.buf, j.prefix)
		writeln(j.buf, msg...)
		os.Stderr.Write(j.buf.Bytes())
		j.buf.Reset()
	}
}

func (j Journal) Debug(msg ...string) {
	if j.debug != nil {
		write(j.buf, j.prefix)
		writeln(j.buf, msg...)
		os.Stderr.Write(j.buf.Bytes())
		j.buf.Reset()
	}
}

func (j Journal) Err(err error) {
	if j.debug != nil {
		write(j.buf, j.prefix)
		writeln(j.buf, err.Error())
		os.Stderr.Write(j.buf.Bytes())
		j.buf.Reset()
	}
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
