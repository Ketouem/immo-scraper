package main

import (
	"flag"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/sirupsen/logrus"

	"github.com/Ketouem/immo-scraper/lib/db"
	"github.com/Ketouem/immo-scraper/lib/scraper"
)

var (
	dumpCSV             bool
	dynamodbEndpointURL string
	leboncoinStartURL   string
	pageLimit           int
	parallelism         int
	verbose             bool
)

func init() {
	flag.IntVar(&parallelism, "parallelism", 1, "number of parallel crawlers")
	flag.IntVar(&pageLimit, "page-limit", -1, "number of results pages to crawl through")
	flag.StringVar(&leboncoinStartURL, "leboncoin-start-url", "", "www.leboncoin.fr start url")
	flag.StringVar(&dynamodbEndpointURL, "dynamodb-endpoint-url", "", "Endpoint for local development")
	flag.BoolVar(&dumpCSV, "export-csv", false, "Exports results to csv file")
	flag.BoolVar(&verbose, "verbose", false, "Display debug logs")
	flag.Parse()

	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}
	logrus.SetFormatter(formatter)

	if len(leboncoinStartURL) == 0 {
		logrus.Fatal("leboncoin-start-url must be provided and not empty")
		os.Exit(0)
	}
}

func main() {
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	results := make([]scraper.Result, 0)
	if len(leboncoinStartURL) != 0 {
		logrus.Info("leboncoin : Fetching data")
		scraper.SetupLeboncoin(parallelism)
		lbcLinks := scraper.GatherLeboncoinLinks(leboncoinStartURL, 2)
		results = append(scraper.ExtractLeboncoinResults(lbcLinks))
		logrus.WithField("results", len(results)).Info("leboncoin: Results fetched")
	}

	db.Setup(dynamodbEndpointURL)
	databaseHandler, _ := db.Get()
	err := db.Provision(databaseHandler)
	if err != nil {
		panic(err)
	}

	logrus.Info("Persisting results")
	err = db.PutResults(databaseHandler, results)
	if err != nil {
		panic(err)
	}

	if dumpCSV {
		dumpResultsToCsv(results)
	}
}

func dumpResultsToCsv(results []scraper.Result) {
	resultsFile, err := os.OpenFile("results.csv", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer resultsFile.Close()

	logrus.Info("Dumping results to CSV")
	err = gocsv.MarshalFile(&results, resultsFile)
	if err != nil {
		panic(err)
	}
}
