package model

import "time"

type GitCommit struct {
	Hash      string    `json:"hash"`
	Author    string    `json:"author"`
	Directory string    `json:"directory"`
	Date      time.Time `json:"date"`
	Message   string    `json:"message"`
}
