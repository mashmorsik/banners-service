package models

import "time"

type Banner struct {
	ID        int   `json:"id"`
	TagIDs    []int `json:"tag_ids"`
	FeatureID int   `json:"feature_id"`
	IsActive  bool  `json:"is_active"`
	//Latest    bool      `json:"use_latest_revision"`
	Content   Content   `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

type Content struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}
