package main

import (
	"fmt"
	"time"
)

type Result struct {
	Link           string    `csv:"link"`
	Source         string    `csv:"source"`
	ExtractionDate time.Time `csv:"extraction_date"`
	Title          string    `csv:"title"`
	Price          int       `csv:"price"`
	LivingSpace    int       `csv:"living_space"`
	Rooms          int       `csv:"rooms"`
}

func (r *Result) Display() {
	fmt.Println("Titre:", r.Title)
	fmt.Println("Prix:", r.Price)
	fmt.Println("Lien:", r.Link)
	fmt.Println("Surface (m2):", r.LivingSpace)
	fmt.Println("Pièces:", r.Rooms)
	fmt.Println("Source:", r.Source)
	fmt.Println("Extrait le", r.ExtractionDate)
}
