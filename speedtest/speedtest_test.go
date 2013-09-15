package speedtest

import "testing"
import (
	"io"
	"io/ioutil"
	"math/rand"
)

func BenchmarkRandomDataMaker(b *testing.B) {
	randomSrc := randomDataMaker{rand.NewSource(1028890720402726901)}
	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(i))
		_, err := io.CopyN(ioutil.Discard, &randomSrc, int64(i))
		if err != nil {
			b.Fatalf("Error copying at %v: %v", i, err)
		}
	}
}
