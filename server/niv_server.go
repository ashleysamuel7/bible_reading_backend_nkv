package server

import (
	"bible_reading_backend_nkv/dto"
	"bible_reading_backend_nkv/server/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func (s *EchoServer) GetAllVerse(ctx echo.Context) error {
	versus, err := s.DB.GetAllVerse(ctx.Request().Context())
	if err != nil {
		log.Printf("Error getting all verses: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch verses"})
	}
	return ctx.JSON(http.StatusOK, versus)

}

func (s *EchoServer) GetAllVerseByChapter(ctx echo.Context) error {
	// Parse bookId
	bookIdStr := ctx.Param("bookId")
	bookId, err := strconv.Atoi(bookIdStr)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid book ID")
	}

	// Parse chapterId
	chapterStr := ctx.Param("chapterId")
	chapter, err := strconv.Atoi(chapterStr)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid chapter number")
	}

	versus, err := s.DB.GetAllVerseByChapter(ctx.Request().Context(), bookId, chapter)
	if err != nil {
		log.Printf("Error getting verses by chapter (bookId: %d, chapter: %d): %v", bookId, chapter, err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch verses"})
	}
	return ctx.JSON(http.StatusOK, versus)

}

func (s *EchoServer) GetAllBook(ctx echo.Context) error {

	versus, err := s.DB.GetAllBook(ctx.Request().Context())
	if err != nil {
		log.Printf("Error getting all books: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch books"})
	}
	return ctx.JSON(http.StatusOK, versus)

}
func (s *EchoServer) GetAllChapter(ctx echo.Context) error {
	// Parse bookId
	bookIdStr := ctx.Param("bookId")
	bookId, err := strconv.Atoi(bookIdStr)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid book ID")
	}

	versus, err := s.DB.GetAllChapter(ctx.Request().Context(), bookId)
	if err != nil {
		log.Printf("Error getting chapters for book (bookId: %d): %v", bookId, err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch chapters"})
	}
	return ctx.JSON(http.StatusOK, versus)

}
func (s *EchoServer) ExpainVerse(ctx echo.Context) error {

	var req dto.ExplainRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request parameters",
		})
	}

	// Try to get user data from token if available
	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			userID, err := utils.ValidateToken(parts[1])
			if err == nil {
				// Token is valid, fetch user data
				user, err := s.DB.GetUserByID(ctx.Request().Context(), userID)
				if err == nil && user != nil {
					// Use user's age and believer_category
					req.Age = user.Age
					req.Belief = user.BelieverCategory
				}
			}
		}
	}

	// Set default values if not provided (fallback if no token or user not found)
	if req.Age == 0 {
		req.Age = 25
	}
	if req.Belief == 0 {
		req.Belief = 3
	}

	// Check if OpenAI API key is set
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	apiKey = strings.Trim(apiKey, `"'`)
	apiKey = strings.TrimSpace(apiKey)

	if apiKey == "" {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "OpenAI API key not configured",
		})
	}

	if !strings.HasPrefix(apiKey, "sk-") {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Invalid API key format",
		})
	}

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
		MaxTokens: 500,
	}

	body, err := json.Marshal(openaiReq)
	if err != nil {
		log.Printf("Error marshaling OpenAI request: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to prepare request"})
	}

	reqHTTP, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create request"})
	}
	reqHTTP.Header.Set("Content-Type", "application/json")
	reqHTTP.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(reqHTTP)
	if err != nil {
		log.Printf("Error calling OpenAI API: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get explanation",
		})
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading OpenAI response: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read response"})
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("OpenAI API returned non-200 status: %d", resp.StatusCode)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get explanation",
		})
	}

	var aiResp dto.OpenAIResponse
	if err := json.Unmarshal(respBody, &aiResp); err != nil {
		log.Printf("Error parsing OpenAI response: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse response",
		})
	}

	if len(aiResp.Choices) == 0 {
		log.Printf("OpenAI API returned empty response")
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "No explanation available"})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"explanation": aiResp.Choices[0].Message.Content,
	})
}
