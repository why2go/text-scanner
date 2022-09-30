package actrie

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestACTrie(t *testing.T) {
	at := NewACTrie()

	at.Put([]rune("a"))
	at.Put([]rune("ab"))
	at.Put([]rune("bab"))
	at.Put([]rune("bc"))
	at.Put([]rune("bca"))
	at.Put([]rune("c"))
	at.Put([]rune("caa"))

	at.ConstructFailureLinks()

	assert.True(t, at.Contains([]rune("bab")), "bab")
	assert.True(t, at.Contains([]rune("caa")), "caa")
	assert.True(t, at.Contains([]rune("ab")), "ab")

	fmt.Printf("at.FindMatches([]rune(\"abccab\")): %v\n", at.FindMatches([]rune("abccab")))
	fmt.Printf("at.FindMatches([]rune(\"dabcdceab\")): %v\n", at.FindMatches([]rune("dabcdceab")))
}
