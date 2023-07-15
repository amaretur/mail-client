package scipher_test

import (
	"testing"

	"github.com/amaretur/spg/pkg/scipher"
)

func TestSKeyGen(t *testing.T) {

	cases := []int{2, 4, 8, 16, 32, 64, 126, 256}

	for i, size := range cases {
		key := scipher.SKeyGen(size)

		t.Logf("Test case #%d", i+1)

		if len(key) != size {
			t.Errorf("Incorrect result! Want: %d but have %d", size, len(key))
		}
	}
}
