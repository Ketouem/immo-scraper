package notifier

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/Ketouem/immo-scraper/lib/scraper"
	"github.com/sirupsen/logrus"
)

var templateFolder string

func Setup(folder string) {
	templateFolder = folder
}

func buildEmail(results []scraper.Result, templateName string) (emailContent string, err error) {
	tmpl, err := template.ParseFiles(templateFolder + fmt.Sprintf("/%s.tmpl", templateName))
	if err != nil {
		return emailContent, err
	}
	var buff bytes.Buffer
	if err = tmpl.Execute(&buff, results); err != nil {
		return emailContent, err
	}

	emailContent = buff.String()
	return emailContent, err
}

func SendEmail(results []scraper.Result) {
	email, err := buildEmail(results, "new-results")
	if err != nil {
		panic(err)
	}
	logrus.Debug("Email content: " + email)
}
