package main

import (
	"bytes"
	"fmt"
	"log"
	"testing"
)

func BenchmarkWithoutPool(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := new(bytes.Buffer)

			fmt.Fprintf(buf, "processing data for worker id: %d", 12345)

			_ = buf.String()
		}
	})
}

func BenchmarkWithPool(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf, ok := bufferPool.Get().(*bytes.Buffer)
			if !ok {
				log.Panic("invalid type")
				return
			}

			fmt.Fprintf(buf, "processing data for worker id: %d", 12345)

			_ = buf.String()

			buf.Reset()
			bufferPool.Put(buf)
		}
	})
}
