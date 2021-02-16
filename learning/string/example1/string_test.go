package example1

import (
	"github.com/stretchr/testify/assert"
	s "strings"
	"testing"
	"unicode"
)

func TestString(t *testing.T) {
	assert.Equal(t, "UPPER STR", s.ToUpper("uPpEr StR"))
	assert.Equal(t, "lower str", s.ToLower("LoWER StR"))
	assert.True(t, s.EqualFold("ArbitRaRyCASeSTr", "aRBITrarYcaSeStR"))
	assert.True(t, s.HasPrefix("Hello World!", "He"))
	assert.True(t, s.HasSuffix("Hello World!", "d!"))
	assert.Equal(t, 2, s.Index("Hello World!", "ll"))
	assert.Equal(t, -1, s.Index("Hello World!", "g"))
	assert.Equal(t, 3, s.Count("Hello World!", "l"))
	assert.Equal(t, 1, s.Count("Hello World!", "ll"))
	assert.Equal(t, "Hello World!", s.TrimSpace("  Hello World! \n "))
	assert.Equal(t, "Hello World! ", s.TrimLeft(" \n Hello World! ", " \n"))
	assert.Equal(t, "\n Hello World!", s.TrimRight("\n Hello World!\n \n", " \n"))
	assert.Equal(t, 1, s.Compare("fds", "abt"))
	assert.Equal(t, 0, s.Compare("GKT", "GKT"))
	assert.Equal(t, -1, s.Compare("Abg", "abg"))
	assert.Equal(t, []string{"a", "b", "cd", "f", "gh"}, s.Fields("a b cd\nf\tgh"))
	assert.Equal(t, []string{"a", "b", "cd", "f", "gh"}, s.Split("a|b|cd|f|gh", "|"))
	assert.Equal(t, "abcdefg", s.Replace("abczzzz", "zzzz", "defg", -1))
	assert.Equal(t, "abcdefgzzzz", s.Replace("abczzzzzzzz", "zzzz", "defg", 1))
	assert.Equal(t, "a,b,c", s.Join([]string{"a", "b", "c"}, ","))
	assert.Equal(t, []string{"abc_", "_", "def_", "g"}, s.SplitAfter("abc__def_g", "_"))
	assert.Equal(t, "abc", s.TrimFunc("abc44", func(r rune) bool {
		return unicode.IsNumber(r)
	}))
}
