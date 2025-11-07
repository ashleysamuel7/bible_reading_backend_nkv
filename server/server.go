package server

import (
	"bible_reading_backend_nkv/database"
	"bible_reading_backend_nkv/models"
	"bible_reading_backend_nkv/server/middleware"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)


type Server interface{
	Start() error
	Readiness(ctx echo.Context) error
	Liveness(ctx echo.Context) error
	GetAllVerse(ctx echo.Context) error
	GetAllVerseByChapter(ctx echo.Context) error
	GetAllChapter(ctx echo.Context) error
	ExpainVerse(ctx echo.Context) error
	
	// Authentication methods
	Register(ctx echo.Context) error
	Login(ctx echo.Context) error
	
	// User management methods
	GetCurrentUser(ctx echo.Context) error
	UpdateCurrentUser(ctx echo.Context) error
	DeleteCurrentUser(ctx echo.Context) error
	
	// Verse tracking methods
	AddFavoriteVerse(ctx echo.Context) error
	GetFavoriteVerses(ctx echo.Context) error
	RemoveFavoriteVerse(ctx echo.Context) error
	AddHighlightedVerse(ctx echo.Context) error
	GetHighlightedVerses(ctx echo.Context) error
	UpdateHighlightedVerse(ctx echo.Context) error
	RemoveHighlightedVerse(ctx echo.Context) error
	UpdateLastRead(ctx echo.Context) error
	GetLastRead(ctx echo.Context) error
	GetLastReadVerses(ctx echo.Context) error
}


type EchoServer struct{
	echo *echo.Echo
	DB database.DatabaseClient
}

// GetEcho returns the echo instance for testing purposes
func (s *EchoServer) GetEcho() *echo.Echo {
	return s.echo
}

func NewEchoServer(db database.DatabaseClient) Server{
	e := echo.New()
	// ✅ CORS configuration
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{
			"https://ashley-samuel.in",
			"http://13.203.234.131:3000",
			"http://localhost:3000",  // React dev server
		},
		AllowMethods: []string{
			echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization,
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
	}))

	server:= &EchoServer{
		echo: e, 
		DB: db,
	}

	server.registerRoutes()
	return server
}


func (s *EchoServer) registerRoutes(){
	s.echo.GET("/readiness", s.Readiness)
	s.echo.GET("/liveness", s.Liveness)

	// Authentication endpoints (public)
	s.echo.POST("/api/register/", s.Register)
	s.echo.POST("/api/login/", s.Login)

	// User-related protected routes with JWT middleware
	protected := s.echo.Group("/api", middleware.JWTAuth())

	// User profile endpoints
	userGroup := protected.Group("/users")
	userGroup.GET("/me", s.GetCurrentUser)
	userGroup.PUT("/me", s.UpdateCurrentUser)
	userGroup.DELETE("/me", s.DeleteCurrentUser)

	// Verse tracking endpoints
	userGroup.POST("/me/favorites", s.AddFavoriteVerse)
	userGroup.GET("/me/favorites", s.GetFavoriteVerses)
	userGroup.DELETE("/me/favorites/:book_id/:chapter/:verse", s.RemoveFavoriteVerse)
	
	userGroup.POST("/me/highlights", s.AddHighlightedVerse)
	userGroup.GET("/me/highlights", s.GetHighlightedVerses)
	userGroup.PUT("/me/highlights/:book_id/:chapter/:verse", s.UpdateHighlightedVerse)
	userGroup.DELETE("/me/highlights/:book_id/:chapter/:verse", s.RemoveHighlightedVerse)
	
	userGroup.POST("/me/last-read", s.UpdateLastRead)
	userGroup.GET("/me/last-read", s.GetLastRead)
	protected.GET("/last-read-verses/", s.GetLastReadVerses)

	// NIV endpoints (public, but explain can use token if provided)
	nivServerGroup := s.echo.Group("/api/niv")
	nivServerGroup.GET("/verses", s.GetAllVerse)
	nivServerGroup.GET("/:bookId/:chapterId/verses", s.GetAllVerseByChapter)
	nivServerGroup.GET("/books", s.GetAllBook)
	nivServerGroup.GET("/chapters/:bookId", s.GetAllChapter)
	nivServerGroup.POST("/explain", s.ExpainVerse)

}


func (s *EchoServer) Start() error{
	if err:= s.echo.Start(":8000"); err != nil && err!= http.ErrServerClosed{
		log.Fatalf("server shutdown occured %s", err)
		return err
	}
	return nil
}




func (s *EchoServer) Readiness(ctx echo.Context) error{
	ready:=s.DB.Ready()
	if ready{
		return ctx.JSON(http.StatusOK, models.Health{Status: "OK"})
	}
	return ctx.JSON(http.StatusInternalServerError, models.Health{Status: "Failure"})
	
}

func (s *EchoServer) Liveness(ctx echo.Context) error{

		return ctx.JSON(http.StatusOK, models.Health{Status: "OK"})
	
}