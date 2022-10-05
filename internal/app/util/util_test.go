package util

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveFirstSymbol(t *testing.T) {
	assert.Equal(t, "bc", RemoveFirstSymbol("abc"))
}

func TestRandomString(t *testing.T) {
	assert.NotSame(t, RandomString(3), RandomString(3))
	assert.Len(t, RandomString(3), 3)
}

func TestGrabHashFromUrl(t *testing.T) {
	s := `http://localhost:8080/VEurjx`
	hash := GrabHashFromURL(s)
	log.Print("hash hash hash ", hash)
}
