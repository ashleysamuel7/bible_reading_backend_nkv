package server

import (
	"bible_reading_backend_nkv/database"
	"bible_reading_backend_nkv/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


type Server interface{
	Start() error
	Readiness(ctx echo.Context) error
	Liveness(ctx echo.Context) error
	GetAllVerse(ctx echo.Context) error
	GetAllVerseByChapter(ctx echo.Context) error
	GetAllChapter(ctx echo.Context) error
	ExpainVerse(ctx echo.Context) error
}


type EchoServer struct{
	echo *echo.Echo
	DB database.DatabaseClient
}

func NewEchoServer(db database.DatabaseClient) Server{
	e := echo.New()
	// ✅ CORS configuration
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
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

	nivServerGroup := s.echo.Group("/api/niv")
	nivServerGroup.GET("/verses", s.GetAllVerse)
	nivServerGroup.GET("/:book/:chapter/verses", s.GetAllVerseByChapter)
	nivServerGroup.GET("/books", s.GetAllBook)
	nivServerGroup.GET("/chapters/:book", s.GetAllChapter)
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