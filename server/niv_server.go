package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"bible_reading_backend_nkv/dto"
	"github.com/labstack/echo/v4"
)

func (s *EchoServer) GetAllVerse(ctx echo.Context) (error) {
	versus, err := s.DB.GetAllVerse(ctx.Request().Context())
	if(err!=nil){
		log.Fatalf("server shutdown occured %s", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, versus)
	
}

func (s *EchoServer) GetAllVerseByChapter(ctx echo.Context) (error) {
	
	chapterStr := ctx.Param("chapter")

	chapter, err := strconv.Atoi(chapterStr) 
	fmt.Printf("chapter s %d", chapter)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid chapter number")
	}
	
	versus, err := s.DB.GetAllVerseByChapter(ctx.Request().Context(), ctx.Param("book"), chapter)
	if(err!=nil){
		log.Fatalf("server shutdown occured %s", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, versus)
	
}

func (s *EchoServer) GetAllBook(ctx echo.Context) (error) {
	
	
	versus, err := s.DB.GetAllBook(ctx.Request().Context())
	if(err!=nil){
		log.Fatalf("server shutdown occured %s", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, versus)
	
}
func (s *EchoServer) GetAllChapter(ctx echo.Context) (error) {
	
	versus, err := s.DB.GetAllChapter(ctx.Request().Context(), ctx.Param("book"))
	if(err!=nil){
		log.Fatalf("server shutdown occured %s", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, versus)
	
}
func (s *EchoServer) ExpainVerse(ctx echo.Context) (error) {
	
	
	var req dto.ExplainRequest

    if err := ctx.Bind(&req); err != nil {
        // log full error server-side too
        fmt.Printf("Bind error: %#v\n", err)
        return ctx.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

	// fmt.Print(req)
	// Construct the OpenAI prompt
	promptIntro := fmt.Sprintf(
		"Context: Book %s, Chapter %d, Verses %d-%d, Age %d, Belief %d/5. "+
			"Use age and belief only to adjust tone and depth. "+
			"Do not mention them in the response. "+
			"Give a clear summary and explain the verses in a simple, relevant way.",
		req.Book, req.Chapter, req.StartVerse, req.EndVerse, req.Age, req.Belief, 
	)

	openaiReq := dto.OpenAIRequest{
		Model: "gpt-4o-mini",
		Messages: []dto.ChatMessage{
			{Role: "system", Content: "You are a helpful assistant that explains Bible verses clearly and simply."},
			{Role: "user", Content: promptIntro},
		},
	}

	body, _ := json.Marshal(openaiReq)

	reqHTTP, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create OpenAI request"})
	}
	reqHTTP.Header.Set("Content-Type", "application/json")
	reqHTTP.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(reqHTTP)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "OpenAI API call failed"})
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var aiResp dto.OpenAIResponse
	if err := json.Unmarshal(respBody, &aiResp); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse OpenAI response"})
	}

	if len(aiResp.Choices) == 0 {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Empty response from OpenAI"})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"explanation": aiResp.Choices[0].Message.Content,
	})
}