package models

import "time"

type Image struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	Filename  string    `json:"filename"`
	URL       string    `json:"url"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	CreatedAt time.Time `json:"createdAt"`
}

type PaginatedResponse struct {
	Images  []Image `json:"images"`
	Total   int     `json:"total"`
	Page    int     `json:"page"`
	PerPage int     `json:"perPage"`
	HasMore bool    `json:"hasMore"`
}
