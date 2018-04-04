package main

import (
  "flag"
  "fmt"
  "strconv"
  "strings"
  "os"

  "github.com/gocolly/colly"
  "github.com/sirupsen/logrus"
)

var (
  leboncoinStartUrl string
)

type Result struct {
  Title string
  Price int
  Link string
  // LivingSpace int
  // Rooms int
}

func (r *Result) Display() {
  fmt.Println("Titre:", r.Title)
  fmt.Println("Prix:", r.Price)
  fmt.Println("Lien:", r.Link)
  // fmt.Println("Surface (m2):", r.LivingSpace)
  // fmt.Println("Pièces:", r.Rooms)
}

func init() {
  flag.StringVar(&leboncoinStartUrl, "leboncoin-start-url", "", "www.leboncoin.fr start url")
  flag.Parse()

  formatter := &logrus.TextFormatter{
    FullTimestamp: true,
  }
  logrus.SetFormatter(formatter)

  if len(leboncoinStartUrl) == 0 {
    logrus.Fatal("leboncoin-start-url must be provided and not empty")
    os.Exit(0)
  }
}

func main() {
  links := gatherLeBonCoinLinks(leboncoinStartUrl)
  extractResults(links)
}

func extractResults(links []string) (results []Result) {
  logrus.WithFields(logrus.Fields{
    "number": len(links),
  }).Info("Extracting results from links")
  c := colly.NewCollector(
    colly.AllowedDomains("www.leboncoin.fr"),
  )

  c.OnHTML("#container", func(e *colly.HTMLElement) {
    rawPrice := strings.Split(e.ChildText("div[data-qa-id=adview_price]"), "€")[0]
    price, _ := strconv.Atoi(strings.Replace(rawPrice, " ", "", -1))
    title := e.ChildText("div[data-qa-id='adview_title']")
    link := "https://" + e.Request.URL.Host + e.Request.URL.Path
    result := Result {
      title,
      price,
      link,
    }
    results = append(results, result)
  })

  for _, link := range links {
    c.Visit(link)
  }
  logrus.WithFields(logrus.Fields{
    "number": len(results),
  }).Info("Results extracted from links ✅")
  return
}

func gatherLeBonCoinLinks(searchUrl string) (links []string){
  logrus.Info("Gathering leboncoin links")
  c := colly.NewCollector(
    colly.AllowedDomains("www.leboncoin.fr"),
  )

  c.OnHTML("li[itemscope] a", func(e *colly.HTMLElement) {
    links = append(links, e.Attr("href"))
  })

  c.OnHTML("#listingAds div.pagination_links_container", func(e *colly.HTMLElement) {
    nextPageLink := e.ChildAttr("#next", "href")
    c.Visit(e.Request.AbsoluteURL(nextPageLink))
  })

  c.Visit(searchUrl)
  return
}
