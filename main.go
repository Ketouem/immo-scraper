package main

import (
  "flag"
  "fmt"
  "strings"
  "os"

  "github.com/gocolly/colly"
  "github.com/sirupsen/logrus"
)

var (
  leboncoinStartUrl string
)

func init() {
  flag.StringVar(&leboncoinStartUrl, "leboncoin-start-url", "", "www.leboncoin.fr start url")
  flag.Parse()

  if len(leboncoinStartUrl) == 0 {
    logrus.Fatal("leboncoin-start-url must be provided and not empty")
    os.Exit(0)
  }
}

func main() {
  c := colly.NewCollector(
    colly.AllowedDomains("www.leboncoin.fr"),
  )

  c.OnHTML("li[itemscope] a", func(e *colly.HTMLElement) {
    fmt.Println("Description:", e.ChildText("h2[class='item_title']"))
    fmt.Println("Prix:", e.ChildAttr("h3[class=item_price]", "content"))
    classifiedHref := e.Attr("href")
    queryStringStart := strings.Index(classifiedHref, "?")
    cleansedHref := ""
    if queryStringStart > -1 {
      cleansedHref = classifiedHref[0:queryStringStart]
    } else {
      cleansedHref = classifiedHref
    }
    fmt.Println("Lien:", e.Request.AbsoluteURL(cleansedHref))
    fmt.Println()
  })

  c.OnHTML("#listingAds div.pagination_links_container", func(e *colly.HTMLElement) {
    nextPageLink := e.ChildAttr("#next", "href")
    c.Visit(e.Request.AbsoluteURL(nextPageLink))
  })

  c.OnRequest(func(r *colly.Request) {
    fmt.Println("Visiting", r.URL)
  })

  c.Visit(leboncoinStartUrl)
}
