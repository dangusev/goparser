package models

// Item from site
type Item struct {
	Id         int64
	Title, Url string
}

// Query from site
type Query struct {
	Text string
}
