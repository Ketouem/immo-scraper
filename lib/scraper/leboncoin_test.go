package scraper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupLeboncoin(t *testing.T) {
	SetupLeboncoin(1)
}

func TestGatherLeboncoinLinks(t *testing.T) {
	links := GatherLeboncoinLinks("https://www.leboncoin.fr/annonces/offres/?th=1&q=appartement", 1)
	assert.NotEmpty(t, links)
}
