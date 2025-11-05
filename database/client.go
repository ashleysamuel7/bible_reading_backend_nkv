package database

import (
	// "context"
	"bible_reading_backend_nkv/models"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Interface for our DB client
type DatabaseClient interface {
	Ready() bool
	GetAllVerse(ctx context.Context) ([]models.NIV, error)
	GetAllVerseByChapter(ctx context.Context, bookId string, chapterId int ) ([]models.NIV, error)
	GetAllBook(ctx context.Context) ([]BookDTO , error)
	GetAllChapter(ctx context.Context,bookId string) (ChapterMaxDTO , error)
	
	// User management methods
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, id int, updates map[string]interface{}) error
	DeleteUser(ctx context.Context, id int) error
	VerifyPassword(ctx context.Context, email, password string) (*models.User, error)
	
	// Favorite verses methods
	AddFavoriteVerse(ctx context.Context, userID, bookID, chapter, verse int) error
	GetFavoriteVerses(ctx context.Context, userID, limit, offset int) ([]models.UserFavoriteVerse, error)
	GetFavoriteVersesCount(ctx context.Context, userID int) (int64, error)
	RemoveFavoriteVerse(ctx context.Context, userID, bookID, chapter, verse int) error
	IsFavoriteVerse(ctx context.Context, userID, bookID, chapter, verse int) (bool, error)
	
	// Highlighted verses methods
	AddHighlightedVerse(ctx context.Context, userID, bookID, chapter, verse int, note, color string) error
	GetHighlightedVerses(ctx context.Context, userID, limit, offset int) ([]models.UserHighlightedVerse, error)
	GetHighlightedVersesCount(ctx context.Context, userID int) (int64, error)
	UpdateHighlightedVerse(ctx context.Context, userID, bookID, chapter, verse int, note, color string) error
	RemoveHighlightedVerse(ctx context.Context, userID, bookID, chapter, verse int) error
	
	// Last read methods
	UpdateLastRead(ctx context.Context, userID, bookID int, bookName string, chapter, verse int) error
	GetLastRead(ctx context.Context, userID int) (*models.UserLastRead, error)
	GetLastReadVerses(ctx context.Context, userID int) ([]models.UserLastRead, error)
}

// Client struct holding gorm DB instance
type Client struct {
	DB *gorm.DB
}

// NewDatabaseClient creates a new MySQL database client
func NewDatabaseClient() (DatabaseClient, error) {
	godotenv.Load()

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("MYSQL_DSN not set in environment")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{},
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		QueryFields: true,
	})
	if err != nil {
		return nil, err
	}

	client := &Client{DB: db}
	return client, nil
}

// Ready checks if DB connection is alive
func (c Client) Ready() bool {
	var ready string
	tx := c.DB.Raw("Select 1 as ready").Scan(&ready)
	if tx.Error != nil {
		return false
	}
	if ready == "1" {
		return true
	}
	return false
}
