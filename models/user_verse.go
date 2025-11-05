package models

import "time"

type UserFavoriteVerse struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    int       `gorm:"column:user_id;not null;index" json:"user_id"`
	BookID    int       `gorm:"column:book_id;not null;index:idx_verse" json:"book_id"`
	Chapter   int       `gorm:"column:chapter;not null;index:idx_verse" json:"chapter"`
	Verse     int       `gorm:"column:verse;not null;index:idx_verse" json:"verse"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName overrides the default pluralized table name
func (UserFavoriteVerse) TableName() string {
	return "user_favorite_verses"
}

type UserHighlightedVerse struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    int       `gorm:"column:user_id;not null;index" json:"user_id"`
	BookID    int       `gorm:"column:book_id;not null;index:idx_verse" json:"book_id"`
	Chapter   int       `gorm:"column:chapter;not null;index:idx_verse" json:"chapter"`
	Verse     int       `gorm:"column:verse;not null;index:idx_verse" json:"verse"`
	Note      string    `gorm:"column:note;type:text" json:"note"`
	Color     string    `gorm:"column:color;size:20;default:yellow" json:"color"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName overrides the default pluralized table name
func (UserHighlightedVerse) TableName() string {
	return "user_highlighted_verses"
}

type UserLastRead struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    int       `gorm:"column:user_id;uniqueIndex;not null" json:"user_id"`
	BookID    int       `gorm:"column:book_id;not null" json:"book_id"`
	BookName  string    `gorm:"column:book_name;not null;size:255" json:"book_name"`
	Chapter   int       `gorm:"column:chapter;not null" json:"chapter"`
	Verse     int       `gorm:"column:verse;not null" json:"verse"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName overrides the default pluralized table name
func (UserLastRead) TableName() string {
	return "user_last_read"
}

