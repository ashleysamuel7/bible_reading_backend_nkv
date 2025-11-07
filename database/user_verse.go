package database

import (
	"bible_reading_backend_nkv/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

// Favorite Verses Methods

func (c Client) AddFavoriteVerse(ctx context.Context, userID, bookID, chapter, verse int) error {
	favorite := models.UserFavoriteVerse{
		UserID:  userID,
		BookID:  bookID,
		Chapter: chapter,
		Verse:   verse,
	}

	result := c.DB.WithContext(ctx).Create(&favorite)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return errors.New("verse already in favorites")
		}
		return result.Error
	}
	return nil
}

func (c Client) GetFavoriteVerses(ctx context.Context, userID, limit, offset int) ([]models.UserFavoriteVerse, error) {
	var favorites []models.UserFavoriteVerse
	result := c.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&favorites)
	return favorites, result.Error
}

func (c Client) GetFavoriteVersesCount(ctx context.Context, userID int) (int64, error) {
	var count int64
	result := c.DB.WithContext(ctx).
		Model(&models.UserFavoriteVerse{}).
		Where("user_id = ?", userID).
		Count(&count)
	return count, result.Error
}

func (c Client) RemoveFavoriteVerse(ctx context.Context, userID, bookID, chapter, verse int) error {
	result := c.DB.WithContext(ctx).
		Where("user_id = ? AND book_id = ? AND chapter = ? AND verse = ?",
			userID, bookID, chapter, verse).
		Delete(&models.UserFavoriteVerse{})
	return result.Error
}

func (c Client) IsFavoriteVerse(ctx context.Context, userID, bookID, chapter, verse int) (bool, error) {
	var count int64
	result := c.DB.WithContext(ctx).
		Model(&models.UserFavoriteVerse{}).
		Where("user_id = ? AND book_id = ? AND chapter = ? AND verse = ?",
			userID, bookID, chapter, verse).
		Count(&count)
	return count > 0, result.Error
}

// Highlighted Verses Methods

func (c Client) AddHighlightedVerse(ctx context.Context, userID, bookID, chapter, verse int, note, color string) error {
	if color == "" {
		color = "yellow"
	}

	highlight := models.UserHighlightedVerse{
		UserID:  userID,
		BookID:  bookID,
		Chapter: chapter,
		Verse:   verse,
		Note:    note,
		Color:   color,
	}

	result := c.DB.WithContext(ctx).Create(&highlight)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return errors.New("verse already highlighted")
		}
		return result.Error
	}
	return nil
}

func (c Client) GetHighlightedVerses(ctx context.Context, userID, limit, offset int) ([]models.UserHighlightedVerse, error) {
	var highlights []models.UserHighlightedVerse
	result := c.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&highlights)
	return highlights, result.Error
}

func (c Client) GetHighlightedVersesCount(ctx context.Context, userID int) (int64, error) {
	var count int64
	result := c.DB.WithContext(ctx).
		Model(&models.UserHighlightedVerse{}).
		Where("user_id = ?", userID).
		Count(&count)
	return count, result.Error
}

func (c Client) UpdateHighlightedVerse(ctx context.Context, userID, bookID, chapter, verse int, note, color string) error {
	updates := map[string]interface{}{}
	if note != "" {
		updates["note"] = note
	}
	if color != "" {
		updates["color"] = color
	}

	if len(updates) == 0 {
		return nil
	}

	result := c.DB.WithContext(ctx).
		Model(&models.UserHighlightedVerse{}).
		Where("user_id = ? AND book_id = ? AND chapter = ? AND verse = ?",
			userID, bookID, chapter, verse).
		Updates(updates)
	
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("highlight not found")
	}
	return nil
}

func (c Client) RemoveHighlightedVerse(ctx context.Context, userID, bookID, chapter, verse int) error {
	result := c.DB.WithContext(ctx).
		Where("user_id = ? AND book_id = ? AND chapter = ? AND verse = ?",
			userID, bookID, chapter, verse).
		Delete(&models.UserHighlightedVerse{})
	return result.Error
}

// Last Read Methods

func (c Client) UpdateLastRead(ctx context.Context, userID, bookID int, bookName string, chapter, verse int) error {
	lastRead := models.UserLastRead{
		UserID:   userID,
		BookID:   bookID,
		BookName: bookName,
		Chapter:  chapter,
		Verse:    verse,
	}

	// Use Clauses to handle upsert
	result := c.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Assign(models.UserLastRead{
			BookID:   bookID,
			BookName: bookName,
			Chapter:  chapter,
			Verse:    verse,
		}).
		FirstOrCreate(&lastRead, models.UserLastRead{UserID: userID})
	return result.Error
}

func (c Client) GetLastRead(ctx context.Context, userID int) (*models.UserLastRead, error) {
	var lastRead models.UserLastRead
	result := c.DB.WithContext(ctx).Where("user_id = ?", userID).First(&lastRead)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &lastRead, nil
}

func (c Client) GetLastReadVerses(ctx context.Context, userID int) ([]models.UserLastRead, error) {
	var lastReads []models.UserLastRead
	result := c.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&lastReads)
	return lastReads, result.Error
}

