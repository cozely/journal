package journal

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"

	"golang.org/x/sys/unix"
)

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
func (j *Journal) InfoTo(path string) *Journal {
	fd, err := unix.Open(path, unix.O_WRONLY, 0)
	if err != nil {
		fd = 0
	}
	j.infofd = fd
	return j
}

// WarnTo changes the file descriptor used to print warnings.
func (j *Journal) WarnTo(path string) *Journal {
	fd, err := unix.Open(path, unix.O_WRONLY, 0)
	if err != nil {
		fd = 0
	}
	j.warnfd = fd
	return j
}

// DebugTo changes the file descriptor used to print debugging entries.
func (j *Journal) DebugTo(path string) *Journal {
	fd, err := unix.Open(path, unix.O_WRONLY, 0)
	if err != nil {
		fd = 0
	}
	j.debugfd = fd
	return j
}

func (j *Journal) commitInfo() {
	unix.Write(j.infofd, j.buf.Bytes())
}

func (j *Journal) commitWarn() {
	unix.Write(j.warnfd, j.buf.Bytes())
}

func (j *Journal) commitDebug() {
	unix.Write(j.debugfd, j.buf.Bytes())
}
