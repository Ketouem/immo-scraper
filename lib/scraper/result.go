package scraper

import (
	"fmt"
	"time"
)

// Result represents a parsed real-state ad
type Result struct {
	Link           string    `csv:"link" json:"link"`
	Source         string    `csv:"source"`
	ExtractionDate time.Time `csv:"extraction_date"`
	Title          string    `csv:"title"`
	Price          int       `csv:"price"`
	LivingSpace    int       `csv:"living_space"`
	Rooms          int       `csv:"rooms"`
	Notified       bool      `csv:"notified" json:"notified"`
}

// NewResult Result constructor
func NewResult(link string, source string, title string, price int, livingSpace int, rooms int) (result *Result) {
	result = &Result{
		Link:           link,
		Source:         source,
		ExtractionDate: time.Now().UTC(),
		Title:          title,
		Price:          price,
		LivingSpace:    livingSpace,
		Rooms:          rooms,
		Notified:       false,
	}
	return
}

// Display to screen the content of a Result
func (r *Result) Display() {
	fmt.Println("Titre:", r.Title)
	fmt.Println("Prix:", r.Price)
	fmt.Println("Lien:", r.Link)
	fmt.Println("Surface (m2):", r.LivingSpace)
	fmt.Println("Pièces:", r.Rooms)
	fmt.Println("Source:", r.Source)
	fmt.Println("Extrait le", r.ExtractionDate)
}
