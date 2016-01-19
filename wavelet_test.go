package wavelet

import (
	"math"
	"sort"
	"strings"
	"testing"
)

func TestRebuild(t *testing.T) {
	var str string = "Sing, Goddess, of the wrath of Achilles"

	ab := alphabetize(str)
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

func TestResize(t *testing.T) {

	ab := []string{"a"}
	chars := []string{"b", "b", "b"}

	wt := NewWaveletTree(ab, chars)
	wt.Rank(3, "b")
}

func TestRank(t *testing.T) {
	var str string = "I’ll make my report as if I told a story, for I was taught as a child on my homeworld that Truth is a matter of the imagination."

	ab := alphabetize(str)
	chars := strings.Split(str, "")

	wt := NewWaveletTree(ab, chars)

	var verify = func(x uint, q string, exp uint) {
		if r := wt.Rank(x, q); r != exp {
			t.Errorf("Rank reported wrong result! %v != %v", r, exp)
			t.Errorf("wt = %v", wt)
		}
	}

	verify(1, "I", 1)
	verify(2, "'", 1)
	verify(3, "l", 1)
	verify(4, "l", 2)

	verify(0, "l", 0)

	verify(uint(len(str)), "r", 7)

	ab = alphabetize("ab")
	chars = strings.Split("c000000", "")
	wt = NewWaveletTree(ab, chars)
	verify(7, "a", 0)
}

func TestSize(t *testing.T) {
	var str = "Far out in the uncharted backwaters of the unfashionable end of the western spiral arm of the Galaxy lies a small unregarded yellow sun."

	ab := alphabetize(str)
	chars := strings.Split(str, "")

	wt := NewWaveletTree(ab, chars)

	upper := uint(math.Ceil(math.Log2(float64(len(ab)))) * float64(len(str)))
	lower := uint(math.Floor(math.Log2(float64(len(ab)))) * float64(len(str)))

	size := wt.Size()

	if size > upper {
		t.Errorf("WaveletTree size greater than upper bound: %v > %v\n", size, upper)
	}
	if size < lower {
		t.Errorf("WaveletTree size less than lower bound: %v < %v\n", size, lower)
	}
}

func BenchmarkNewWaveletTree(b *testing.B) {
	var str string = "As Gregor Samsa awoke one morning from uneasy dreams he found himself transformed in his bed into a monstrous vermin."

	ab := alphabetize(str)
	chars := strings.Split(str, "")

	for i := 0; i < b.N; i++ {
		NewWaveletTree(ab, chars)
	}
}

func BenchmarkIter(b *testing.B) {
	var str string = "Call me Ishmael."

	ab := alphabetize(str)
	chars := strings.Split(str, "")

	wt := NewWaveletTree(ab, chars)

	for i := 0; i < b.N; i++ {
		var newstr string
		for x := range wt.Iter() {
			newstr += x
		}
	}
}

// alphabetize is a helper function to take a string and return
// all the distinct symbols it uses.
func alphabetize(s string) []string {
	var r []string
	var chars []string = strings.Split(s, "")

	sort.Strings(chars)

	r = append(r, chars[0])

	for _, x := range chars {
		if x != r[len(r)-1] {
			r = append(r, x)
		}
	}
	return r
}
