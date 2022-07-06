package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrimFirstRune(t *testing.T) {
	assert.Equal(t, "bc", trimFirstRune("abc"))
}

func TestRandStringRunes(t *testing.T) {
	assert.NotSame(t, RandStringRunes(3), RandStringRunes(3))
	assert.Len(t, RandStringRunes(3), 3)
}
