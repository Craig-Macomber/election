package merkle

import (
	"crypto"
	_ "crypto/sha512"
	_ "crypto/sha256"
	"testing"
)


func BenchmarkSha512(b *testing.B) {
    h := crypto.SHA512.New()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        data:=h.Sum(make([]byte, 0))
        h.Write(data)
        h.Write(data)
    }
}

func BenchmarkSha256(b *testing.B) {
    h := crypto.SHA256.New()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        data:=h.Sum(make([]byte, 0))
        h.Write(data)
        h.Write(data)
    }
}