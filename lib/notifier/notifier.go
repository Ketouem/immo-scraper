package notifier

import (
  "bytes"
  "html/template"

  "github.com/Ketouem/immo-scraper/lib/scraper"
)

const TEMPLATE_FOLDER = "./templates"

func buildEmail(results []scraper.Result) (emailContent string, err error){
  tmpl, err := template.ParseFiles(TEMPLATE_FOLDER + "/new-results.tmpl")
  if err != nil {
    return emailContent, err
  }
  var buff bytes.Buffer
  if err := tmpl.Execute(&buff, results); err != nil {
      return emailContent, err
  }

  emailContent = buff.String()
  return emailContent, err
}
