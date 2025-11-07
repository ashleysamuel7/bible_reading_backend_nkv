package dto

type ExplainRequest struct {
	Book       string `json:"book"`
	Chapter    int    `json:"chapter"`
	StartVerse int    `json:"start_verse"`
	EndVerse   int    `json:"end_verse"`
	Age        int    `json:"age"`
	Belief     int    `json:"belief"`
}

type OpenAIRequest struct {
	Model     string        `json:"model"`
	Messages  []ChatMessage `json:"messages"`
	MaxTokens int           `json:"max_tokens"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}
