package notifier

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Ketouem/immo-scraper/lib/scraper"
)

func TestSetupLeboncoin(t *testing.T) {
	Setup("../../templates")
}

func TestBuildEmail_nominal(t *testing.T) {
	results := []scraper.Result{
		scraper.Result{
			Link:           "https://www.test.com/1",
			Source:         "leboncoin",
			ExtractionDate: time.Now(),
			Title:          "My First Test Title",
			Price:          800000,
			LivingSpace:    90,
			Rooms:          4,
		},
		scraper.Result{
			Link:           "https://www.test.com/2",
			Source:         "leboncoin",
			ExtractionDate: time.Now(),
			Title:          "My Second Test Title",
			Price:          900000,
			LivingSpace:    100,
			Rooms:          5,
		},
	}
	emailContent, err := buildEmail(results, "new-results")
	if err != nil {
		t.Error(err)
	}
	for _, result := range results {
		assert.Contains(t, emailContent, result.Link)
		assert.Contains(t, emailContent, result.Title)
		assert.Contains(t, emailContent, strconv.Itoa(result.Price))
		assert.Contains(t, emailContent, strconv.Itoa(result.LivingSpace))
		assert.Contains(t, emailContent, strconv.Itoa(result.Rooms))
	}
}

func TestBuildEmail_templateError(t *testing.T) {
	results := []scraper.Result{
		scraper.Result{
			Link:           "https://www.test.com/1",
			Source:         "leboncoin",
			ExtractionDate: time.Now(),
			Title:          "My First Test Title",
			Price:          800000,
			LivingSpace:    90,
			Rooms:          4,
		},
	}
	emailContent, err := buildEmail(results, "invalid-template-name")
	assert.Empty(t, emailContent)
	assert.Error(t, err)
	emailContent, err = buildEmail(results, "invalid-template")
	assert.Empty(t, emailContent)
	assert.Error(t, err)
}
