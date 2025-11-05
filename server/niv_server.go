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
		log.Fatalf("server shutdown occured %s", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, versus)

}

func (s *EchoServer) GetAllVerseByChapter(ctx echo.Context) error {

	chapterStr := ctx.Param("chapter")

	chapter, err := strconv.Atoi(chapterStr)
	fmt.Printf("chapter s %d", chapter)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid chapter number")
	}

	versus, err := s.DB.GetAllVerseByChapter(ctx.Request().Context(), ctx.Param("book"), chapter)
	if err != nil {
		log.Fatalf("server shutdown occured %s", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, versus)

}

func (s *EchoServer) GetAllBook(ctx echo.Context) error {

	versus, err := s.DB.GetAllBook(ctx.Request().Context())
	if err != nil {
		log.Fatalf("server shutdown occured %s", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, versus)

}
func (s *EchoServer) GetAllChapter(ctx echo.Context) error {

	versus, err := s.DB.GetAllChapter(ctx.Request().Context(), ctx.Param("book"))
	if err != nil {
		log.Fatalf("server shutdown occured %s", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, versus)

}
func (s *EchoServer) ExpainVerse(ctx echo.Context) error {

	var req dto.ExplainRequest

	if err := ctx.Bind(&req); err != nil {
		log.Printf("Bind error: %v", err)
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
		log.Printf("ERROR: OPENAI_API_KEY not set")
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "OpenAI API key not configured",
		})
	}

	if !strings.HasPrefix(apiKey, "sk-") {
		log.Printf("ERROR: Invalid API key format")
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
		log.Printf("ERROR: Failed to marshal request: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to prepare request"})
	}

	reqHTTP, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("ERROR: Failed to create request: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create request"})
	}
	reqHTTP.Header.Set("Content-Type", "application/json")
	reqHTTP.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(reqHTTP)
	if err != nil {
		log.Printf("ERROR: API call failed: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get explanation",
		})
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read response: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read response"})
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR: API returned status %d", resp.StatusCode)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get explanation",
		})
	}

	var aiResp dto.OpenAIResponse
	if err := json.Unmarshal(respBody, &aiResp); err != nil {
		log.Printf("ERROR: Failed to parse response: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse response",
		})
	}

	if len(aiResp.Choices) == 0 {
		log.Printf("ERROR: Empty response from API")
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "No explanation available"})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"explanation": aiResp.Choices[0].Message.Content,
	})
}
