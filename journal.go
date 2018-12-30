package journal

import (
	"bytes"
	"io"
	"os"
)

////////////////////////////////////////////////////////////////////////////////

var (
	outbuf = bytes.NewBuffer(make([]byte, 0, 256))
	errbuf = bytes.NewBuffer(make([]byte, 0, 256))
)

////////////////////////////////////////////////////////////////////////////////

type Journal string

func (j Journal) Info(msg ...string) {
	if j != "" {
		write(outbuf, string(j))
		writeln(outbuf, msg...)
		os.Stdout.Write(outbuf.Bytes())
		outbuf.Reset()
	}
}

func (j Journal) Warn(msg ...string) {
	if j != "" {
		write(errbuf, string(j))
		writeln(errbuf, msg...)
		os.Stderr.Write(errbuf.Bytes())
		errbuf.Reset()
	}
}

func (j Journal) Debug(msg ...string) {
	if j != "" {
		write(errbuf, string(j))
		writeln(errbuf, msg...)
		os.Stderr.Write(errbuf.Bytes())
		errbuf.Reset()
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
