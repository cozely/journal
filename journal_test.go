package journal_test

import (
	"testing"

	"github.com/cozely/journal"
	"golang.org/x/sys/unix"
)

func init() {
	null, _ := unix.Open("/dev/null", unix.O_WRONLY, 0)
	journal.InfoTo(null)
	journal.WarnTo(null)
	journal.DebugTo(null)
}

var log = journal.New("benchmark: ", true)

func BenchmarkInfoInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Info("Foo %d", n)
	}
}

func BenchmarkInfo3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Info("Foo%s%s", "Bar", "Baz")
	}
}

func BenchmarkWarn3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Warn("Foo%s%s", "Bar", "Baz")
	}
}

func BenchmarkDebug3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Debug("Foo%s%s", "Bar", "Baz")
	}
}
