package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindElRound(t *testing.T) {
	title := " Olympiacos-AS Monaco | Round 19 Highlights | 2023-24 Turkish Airlines EuroLeague "
	res, err := ExtractElRound(title)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, 19)
}
