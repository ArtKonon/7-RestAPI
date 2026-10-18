// Package model holds the domain type and the request structs for your resource.
package model

type Book struct {
	ID            string `json:"id"`
	Title         string `json:"title" binding:"required"`
	ISBN          string `json:"isbn"`
	Author        string `json:"author"`         // тимчасово текстове поле
	Category      string `json:"category"`       // тимчасово текстове поле
	PublishedYear int    `json:"published_year"`
	Description   string `json:"description"`
} 