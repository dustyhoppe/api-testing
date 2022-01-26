package models

type MovieSearchResult struct {
	Id            int     `json:"id"`
	OriginalTitle string  `json:"original_title"`
	Popularity    float32 `json:"popularity"`
}
