package journal_test

import (
	"log"
	"testing"

	"github.com/cozely/journal"
)

var (
	null, _ = journal.Open("/dev/null")
	jnal    = journal.New().InfoTo(null).WarnTo(null).DebugTo(null)
)

func init() {
	null, _ := journal.Open("/dev/null")
	log.SetOutput(null)
	log.SetFlags(0)
}

// var jnal = journal.New().InfoTo(nil).WarnTo(nil).DebugTo(nil)

func BenchmarkInfo(b *testing.B) {
	for n := 0; n < b.N; n++ {
		jnal.Info("Foo%s%s", "Bar", "Baz")
	}
}

func BenchmarkInfoInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		jnal.Info("Foo %d", n)
	}
}

func BenchmarkWarn(b *testing.B) {
	for n := 0; n < b.N; n++ {
		jnal.Warn("Foo%s%s", "Bar", "Baz")
	}
}

func BenchmarkDebug(b *testing.B) {
	for n := 0; n < b.N; n++ {
		jnal.Debug("Foo%s%s", "Bar", "Baz")
	}
}

func BenchmarkStdLogPrint(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Print("Foo", "Bar", "Baz")
	}
}

func BenchmarkStdLogPrintf(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Printf("Foo%s%s", "Bar", "Baz")
	}
}

func BenchmarkStdLogPrintfInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Printf("Foo %d", n)
	}
}
