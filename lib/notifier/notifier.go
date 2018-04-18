package notifier

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/Ketouem/immo-scraper/lib/scraper"
)

// TemplateFolder : Where the tmpl files are stored
const TemplateFolder = "./templates"

func buildEmail(results []scraper.Result, templateName string) (emailContent string, err error) {
	tmpl, err := template.ParseFiles(TemplateFolder + fmt.Sprintf("/%s.tmpl", templateName))
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
