package database

type BookDTO struct {
    BookID int    `json:"book_id"`
    Book   string `json:"book"`
}

type ChapterMaxDTO struct {
    MaxChapter int64    `gorm:"column:maxChapter"`
}