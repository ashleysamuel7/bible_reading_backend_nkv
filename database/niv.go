package database

import (
	"bible_reading_backend_nkv/models"
	"golang.org/x/net/context"
)


func (c Client) GetAllVerse(ctx context.Context) ([]models.NIV , error){
	var verse []  models.NIV
	result:= c.DB.WithContext(ctx).Find(&verse)
	return verse, result.Error
}

func (c Client) GetAllVerseByChapter(ctx context.Context, bookId int, chapterId int) ([]models.NIV , error){
	var verse []  models.NIV
	result:= c.DB.WithContext(ctx).Where("chapter = ? AND book_id = ?", chapterId, bookId).Find(&verse)
	return verse, result.Error
}

func (c Client) GetAllBook(ctx context.Context) ([]BookDTO , error){
	var verse []  BookDTO
	result:= c.DB.WithContext(ctx).Model(&models.NIV{}).Distinct("book_id, book").Find(&verse)
	return verse, result.Error
}
func (c Client) GetAllChapter(ctx context.Context, bookId int) ( ChapterMaxDTO , error){
	var chap  ChapterMaxDTO


    query := "SELECT MAX(chapter) as maxChapter FROM niv WHERE book_id = ? "
    result := c.DB.WithContext(ctx).Raw(query, bookId).Scan(&chap)

	return chap, result.Error
}