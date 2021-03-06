package scraper

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

// SOURCE is used to tag the results
const SOURCE = "leboncoin"

var (
	c = colly.NewCollector(colly.Async(true))
)

// SetupLeboncoin is used to perform additional tweak (e.g. colly's Collector)
func SetupLeboncoin(parallelism int) {
	c.Limit(&colly.LimitRule{DomainGlob: "*leboncoin.fr*", Parallelism: parallelism})
}

// GatherLeboncoinLinks returns a slice of ad links given a search URL
func GatherLeboncoinLinks(searchURL string, pageLimit int) (links []string) {
	logrus.Debug("Gathering leboncoin links")

	c.OnHTML("li[itemscope] a", func(e *colly.HTMLElement) {
		links = append(links, e.Attr("href"))
	})

	c.OnHTML("#listingAds div.pagination_links_container", func(e *colly.HTMLElement) {
		nextPageURL, _ := url.Parse(e.Request.AbsoluteURL(e.ChildAttr("#next", "href")))
		nextPageIndex, _ := strconv.Atoi(nextPageURL.Query().Get("o"))
		if pageLimit != -1 {
			if nextPageIndex >= pageLimit {
				logrus.Debug("Page limit reached, skipping link.")
			} else {
				c.Visit(nextPageURL.String())
			}
		} else {
			c.Visit(nextPageURL.String())
		}
	})

	c.Visit(searchURL)
	c.Wait()
	return
}

// ExtractLeboncoinResults parse the pages for each link inside the slice
func ExtractLeboncoinResults(links []string) (results []Result) {
	logrus.WithFields(logrus.Fields{
		"number": len(links),
	}).Debug("Extracting results from links")

	c.OnHTML("#container", func(e *colly.HTMLElement) {
		rawPrice := strings.Split(e.ChildText("div[data-qa-id=adview_price]"), "€")[0]
		price, _ := strconv.Atoi(strings.Replace(rawPrice, " ", "", -1))
		title := e.ChildText("div[data-qa-id='adview_title'] div:first-child h3")
		link := "https://" + e.Request.URL.Host + e.Request.URL.Path
		rawLivingSpace := e.ChildText("div[data-qa-id='criteria_item_square'] div div:nth-child(2)")
		livingSpace, _ := strconv.Atoi(rawLivingSpace[:len(rawLivingSpace)-4])
		rooms, _ := strconv.Atoi(e.ChildText("div[data-qa-id='criteria_item_rooms'] div div:nth-child(2)"))
		result := NewResult(link, SOURCE, title, price, livingSpace, rooms)
		results = append(results, *result)
	})

	for _, link := range links {
		c.Visit(link)
	}
	c.Wait()
	logrus.WithFields(logrus.Fields{
		"number": len(results),
	}).Debug("Results extracted from links")
	return
}
