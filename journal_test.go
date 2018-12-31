package journal_test

import (
	"log"
	"testing"

	"github.com/cozely/journal"
)

var (
	null, _ = journal.Open("/dev/null")
	info    = journal.New(null, "", journal.LstdFlags)
	nilj    = journal.New(nil, "")
)

func init() {
	log.SetOutput(null)
	log.SetFlags(0)
	journal.SetOutput(null)
}

func BenchmarkPrint(b *testing.B) {
	for n := 0; n < b.N; n++ {
		info.Print("Foo", "Bar", "Baz")
	}
}

func BenchmarkPrintInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		info.Print("Foo", n)
	}
}

func BenchmarkPrintln(b *testing.B) {
	for n := 0; n < b.N; n++ {
		info.Println("Foo", "Bar", "Baz")
	}
}

func BenchmarkPrintlnInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		info.Println("Foo", n)
	}
}

func BenchmarkPrintf(b *testing.B) {
	for n := 0; n < b.N; n++ {
		info.Printf("Foo%s%s", "Bar", "Baz")
	}
}

func BenchmarkPrintfInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		info.Printf("Foo %d", n)
	}
}

func BenchmarkNilPrintf(b *testing.B) {
	for n := 0; n < b.N; n++ {
		nilj.Printf("Foo%s%s", "Bar", "Baz")
	}
}

func BenchmarkNilPrintfInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		nilj.Printf("Foo %d", n)
	}
}

func BenchmarkStdPrintf(b *testing.B) {
	for n := 0; n < b.N; n++ {
		journal.Printf("Foo%s%s", "Bar", "Baz")
	}
}

func BenchmarkStdPrintfInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		journal.Printf("Foo %d", n)
	}
}

func BenchmarkLogPrintf(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Printf("Foo%s%s", "Bar", "Baz")
	}
}

func BenchmarkLogPrintfInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Printf("Foo %d", n)
	}
}
