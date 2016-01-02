package wavelet

import (
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

func TestRank(t *testing.T) {
	var str string = "I’ll make my report as if I told a story, for I was taught as a child on my homeworld that Truth is a matter of the imagination."

	ab := alphabetize(str)
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

func TestSelect(t *testing.T) {
	var str = "\"Where's Papa going with that axe?\" said Fern to her mother as they were setting the table for breakfast."

	ab := alphabetize(str)
	chars := strings.Split(str, "")

	wt := NewWaveletTree(ab, chars)

	var verify = func(x uint, q string, exp uint) {
		if r := wt.Select(x, q); r != exp {
			t.Errorf("Select reported wrong result! %v != %v", r, exp)
		}
	}

	verify(1, "\"", 1)
	verify(1, "W", 2)
	verify(1, "h", 3)
	verify(1, "e", 4)
	verify(1, "r", 5)
	verify(2, "e", 6)

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

