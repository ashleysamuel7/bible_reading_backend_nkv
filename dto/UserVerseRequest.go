package dto

import "time"

type AddFavoriteVerseRequest struct {
	BookID  int `json:"book_id" validate:"required,min=1"`
	Chapter int `json:"chapter" validate:"required,min=1"`
	Verse   int `json:"verse" validate:"required,min=1"`
}

type AddHighlightRequest struct {
	BookID  int    `json:"book_id" validate:"required,min=1"`
	Chapter int    `json:"chapter" validate:"required,min=1"`
	Verse   int    `json:"verse" validate:"required,min=1"`
	Note    string `json:"note,omitempty"`
	Color   string `json:"color,omitempty" validate:"omitempty,max=20"`
}

type UpdateLastReadRequest struct {
	BookID   int    `json:"book_id" validate:"required,min=1"`
	BookName string `json:"book_name" validate:"required,min=1,max=255"`
	Chapter  int    `json:"chapter" validate:"required,min=1"`
	Verse    int    `json:"verse" validate:"required,min=1"`
}

type VerseReferenceResponse struct {
	BookID   int    `json:"book_id"`
	BookName string `json:"book_name"`
	Chapter  int    `json:"chapter"`
	Verse    int    `json:"verse"`
	Text     string `json:"text"`
}

type FavoriteVerseResponse struct {
	ID        int                     `json:"id"`
	UserID    int                     `json:"user_id"`
	BookID    int                     `json:"book_id"`
	BookName  string                  `json:"book_name"`
	Chapter   int                     `json:"chapter"`
	Verse     int                     `json:"verse"`
	Text      string                  `json:"text"`
	CreatedAt time.Time               `json:"created_at"`
}

type HighlightedVerseResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	BookID    int       `json:"book_id"`
	BookName  string    `json:"book_name"`
	Chapter   int       `json:"chapter"`
	Verse     int       `json:"verse"`
	Text      string    `json:"text"`
	Note      string    `json:"note"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LastReadResponse struct {
	UserID    int       `json:"user_id"`
	BookID    int       `json:"book_id"`
	BookName  string    `json:"book_name"`
	Chapter   int       `json:"chapter"`
	Verse     int       `json:"verse"`
	Text      string    `json:"text"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

