package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

// SOURCE is used to tag the results
const SOURCE = "leboncoin"

var (
	c *colly.Collector = colly.NewCollector(
		colly.Async(true),
	)
)

func init() {
	c.Limit(&colly.LimitRule{DomainGlob: "*leboncoin.fr*", Parallelism: parallelism})
}

func gatherLeboncoinLinks(searchURL string) (links []string) {
	logrus.Info("Gathering leboncoin links")

	c.OnHTML("li[itemscope] a", func(e *colly.HTMLElement) {
		links = append(links, e.Attr("href"))
	})

	c.OnHTML("#listingAds div.pagination_links_container", func(e *colly.HTMLElement) {
		nextPageLink := e.ChildAttr("#next", "href")
		c.Visit(e.Request.AbsoluteURL(nextPageLink))
	})

	c.Visit(searchURL)
	c.Wait()
	return
}

func extractLeboncoinResults(links []string) (results []Result) {
	logrus.WithFields(logrus.Fields{
		"number": len(links),
	}).Info("Extracting results from links")

	c.OnHTML("#container", func(e *colly.HTMLElement) {
		rawPrice := strings.Split(e.ChildText("div[data-qa-id=adview_price]"), "€")[0]
		price, _ := strconv.Atoi(strings.Replace(rawPrice, " ", "", -1))
		title := e.ChildText("div[data-qa-id='adview_title'] div:first-child h3")
		link := "https://" + e.Request.URL.Host + e.Request.URL.Path
		rawLivingSpace := e.ChildText("div[data-qa-id='criteria_item_square'] div div:nth-child(2)")
		livingSpace, _ := strconv.Atoi(rawLivingSpace[:len(rawLivingSpace)-4])
		rooms, _ := strconv.Atoi(e.ChildText("div[data-qa-id='criteria_item_rooms'] div div:nth-child(2)"))
		result := Result{
			link,
			SOURCE,
			time.Now().UTC(),
			title,
			price,
			livingSpace,
			rooms,
		}
		results = append(results, result)
	})

	for _, link := range links {
		c.Visit(link)
	}
	c.Wait()
	logrus.WithFields(logrus.Fields{
		"number": len(results),
	}).Info("Results extracted from links")
	return
}
