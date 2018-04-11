package notifier

import (
  "strconv"
  "testing"
  "time"

  "github.com/stretchr/testify/assert"

  "github.com/Ketouem/immo-scraper/lib/scraper"
)

func TestBuildEmail(t *testing.T) {
  results := []scraper.Result{
    scraper.Result{
      "https://www.test.com/1",
      "leboncoin",
      time.Now(),
      "My First Test Title",
      800000,
      90,
      4,
    },
    scraper.Result{
      "https://www.test.com/2",
      "leboncoin",
      time.Now(),
      "My Second Test Title",
      900000,
      100,
      5,
    },
  }
  emailContent, err := buildEmail(results)
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
