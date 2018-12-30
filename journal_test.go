package journal_test

import (
	"testing"

	"github.com/cozely/journal"
)

var jnal = journal.New().InfoTo("/dev/null").WarnTo("/dev/null").DebugTo("/dev/null")

func BenchmarkInfoInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		jnal.Info("Foo %d", n)
	}
}

func BenchmarkInfo3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		jnal.Info("Foo%s%s", "Bar", "Baz")
	}
}

func BenchmarkWarn3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		jnal.Warn("Foo%s%s", "Bar", "Baz")
	}
}

func BenchmarkDebug3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		jnal.Debug("Foo%s%s", "Bar", "Baz")
	}
}
