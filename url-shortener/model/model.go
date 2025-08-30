package model

import "time"

type CreateLinkRequest struct {
	Link            string  `json:"link"`
	CustomShortLink *string `json:"custom_short_link"`
}

type LinkResponse struct {
	LongLink  string `json:"long_link"`
	ShortLink string `json:"short_link"`
}

type AnalyticsResponse struct {
	TotalRedirects int        `json:"total_redirects"`
	Redirects      []Redirect `json:"redirects"`
}

type StoreRedirectParams struct {
	UserAgent string
	LongLink  string
	ShortLink string
}

type Redirect struct {
	Id        int       `json:"id"`
	LongLink  string    `json:"long_link"`
	ShortLink string    `json:"short_link"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

type LinkPair struct {
	Short string
	Long  string
}
