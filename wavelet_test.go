package wavelet

import (
	"testing"
	"strings"
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
