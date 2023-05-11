package entity

import (
	"gorm.io/gorm"
)

type Template struct {
	gorm.Model
	Type              string `json:"type"`
	Description       string `json:"description"`
	Message           string `json:"message"`
	EmailSubject      string `json:"email_subject"`
	EmailMarkup       string `json:"markup"`
	EmailRenderedHTML string `json:"email_rendered_html"`
	EventId           uint   `json:"event_id"`
	Active            bool   `json:"active"`
	Title             string `json:"title"`
	Link              string `json:"link"`
	File              string `json:"file"`
	Language          string `json:"language"`
}
