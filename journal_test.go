package journal_test

import (
	"testing"

	"github.com/cozely/journal"
	"golang.org/x/sys/unix"
)

func init() {
	null, _ := unix.Open("/dev/null", unix.O_WRONLY, 0)
	journal.PrintTo(null)
	journal.WarnTo(null)
	journal.DebugTo(null)
}
var log = journal.New("benchmark: ", true)

func BenchmarkPut(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Put("Foo")
	}
}

func BenchmarkPrint(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Print("Foo")
	}
}

func BenchmarkPrint3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Print("Foo", "Bar", "Baz")
	}
}

func BenchmarkPrintf(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Printf("Foo %d", n)
	}
}

func BenchmarkPrintf3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Printf("Foo%s%s", "Bar", "Baz")
	}
}

func BenchmarkWarn(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Warn("Foo", "Bar", "Baz")
	}
}

func BenchmarkDebug(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Warn("Foo", "Bar", "Baz")
	}
}
