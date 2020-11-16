package models

import "time"

// Repo is the response json expected from the Github API
type Repo struct {
	ID          int       `json:"id"`
	URL         string    `json:"html_url"`
	Description string    `json:"description"`
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Forks       int       `json:"forks"`
	CreatedAt   time.Time `json:"created_at"`
	Private     bool      `json:"private"`
}

// GistResponse is the response json expected after creating a gist
type GistResponse struct {
	URL         string `json:"html_url"`
	Description string `json:"description"`
}

// File contains the struct tpye of file
type File struct {
	Content string `json:"content"`
}

// Gist is the body used when creating a gist
type Gist struct {
	Description string          `json:"description"`
	Files       map[string]File `json:"files"`
	Public      bool            `json:"public"`
}
