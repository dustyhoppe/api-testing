package models

type SearchMoviesResponse struct {
	Page    int                 `json:"page"`
	Results []MovieSearchResult `json:"results"`
}

type CreateGuestSessionResponse struct {
	Success        bool   `json:"success"`
	GuestSessionId string `json:"guest_session_id"`
	ExpiresAt      string `json:"expires_at"`
}