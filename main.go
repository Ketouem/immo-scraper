package main

import (
	"flag"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/sirupsen/logrus"
)

var (
	leboncoinStartURL string
	parallelism       int
)

func init() {
	flag.IntVar(&parallelism, "parallelism", 1, "number of parallel crawlers")
	flag.StringVar(&leboncoinStartURL, "leboncoin-start-url", "", "www.leboncoin.fr start url")
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
	results := make([]Result, 0)
	if len(leboncoinStartURL) != 0 {
		logrus.Info("Fetching data from leboncoin")
		lbcLinks := gatherLeboncoinLinks(leboncoinStartURL)
		results = append(extractLeboncoinResults(lbcLinks))
	}
	dumpResultsToCsv(results)
}

func dumpResultsToCsv(results []Result) {
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
