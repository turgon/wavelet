package wavelet

import (
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {
	var str string = "Sing, Goddess, of the wrath of Achilles"

	ab := alphabet(str)
	chars := strings.Split(str, "")

	wt := NewWaveletTree(ab, chars)

	var newstr string

	for x := range wt.Iter() {
		newstr += x
	}

	if str != newstr {
		t.Error("Couldn't rebuild input!", str, newstr)
	}
}

func TestMissingAlphaChars(t *testing.T) {
	var str string = "It was a pleasure to burn."

	ab := []string{"p", "l", "e", "a", "s", "u", "r", "e"}
	chars := strings.Split(str, "")

	wt := NewWaveletTree(ab, chars)

	var newstr string

	for x := range wt.Iter() {
		newstr += x
	}

	if "asapleasureur" != newstr {
		t.Error("Couldn't rebuild input!", str)
	}
}

func TestRank(t *testing.T) {
	var str string = "I’ll make my report as if I told a story, for I was taught as a child on my homeworld that Truth is a matter of the imagination."

	ab := alphabet(str)
	chars := strings.Split(str, "")

	wt := NewWaveletTree(ab, chars)

	var verify = func(x uint, q string, exp uint) {
		if r := wt.Rank(x, q); r != exp {
			t.Errorf("Rank reported wrong result! %v != %v", r, exp)
		}
	}

	verify(1, "I", 1)
	verify(2, "'", 1)
	verify(3, "l", 1)
	verify(4, "l", 2)

	verify(0, "l", 0)

	verify(uint(len(str)), "r", 7)
}

func BenchmarkNewWaveletTree(b *testing.B) {
	var str string = "As Gregor Samsa awoke one morning from uneasy dreams he found himself transformed in his bed into a monstrous vermin."

	ab := alphabet(str)
	chars := strings.Split(str, "")

	for i := 0; i < b.N; i++ {
		NewWaveletTree(ab, chars)
	}
}

func BenchmarkIter(b *testing.B) {
	var str string = "Call me Ishmael."

	ab := alphabet(str)
	chars := strings.Split(str, "")

	wt := NewWaveletTree(ab, chars)

	for i := 0; i < b.N; i++ {
		var newstr string
		for x := range wt.Iter() {
			newstr += x
		}
	}
}
