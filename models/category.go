package models

type Category struct {
	ID          string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Published   bool   `json:"published"`
}
