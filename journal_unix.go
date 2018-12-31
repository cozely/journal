package journal

import (
	"golang.org/x/sys/unix"
)

////////////////////////////////////////////////////////////////////////////////

type File int

var (
	Stdout = File(unix.Stdout)
	Stderr = File(unix.Stderr)
)

func Open(path string) (File, error) {
	fd, err := unix.Open(path, unix.O_WRONLY, 0)
	return File(fd), err
}

func (f File) Close(path string) error {
	return unix.Close(int(f))
}

func (f File) Write(p []byte) (n int, err error) {
	return unix.Write(int(f), p)
}
