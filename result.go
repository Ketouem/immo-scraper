package main

import (
	"fmt"
)

type Result struct {
  Source string `csv:"source"`
	Title string `csv:"title"`
	Price int    `csv:"price"`
	Link  string `csv:"link"`
	// LivingSpace int
	// Rooms int
}

func (r *Result) Display() {
	fmt.Println("Titre:", r.Title)
	fmt.Println("Prix:", r.Price)
	fmt.Println("Lien:", r.Link)
	// fmt.Println("Surface (m2):", r.LivingSpace)
	// fmt.Println("Pièces:", r.Rooms)
  fmt.Println("Source:", r.Source)
}
