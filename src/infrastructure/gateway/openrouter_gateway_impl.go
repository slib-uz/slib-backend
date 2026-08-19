package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/infrastructure/config"
)

const (
	openRouterTimeout      = 180 * time.Second
	openRouterExtractModel = "deepseek/deepseek-v4-flash-0731"
)

type OpenRouterGatewayImpl struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

// @inject
func NewOpenRouterGateway(env *config.Config) gateway.OpenRouterGateway {
	return &OpenRouterGatewayImpl{
		client:  &http.Client{Timeout: openRouterTimeout},
		baseURL: strings.TrimRight(env.OpenRouterBaseURL, "/"),
		apiKey:  env.OpenRouterAPIKey,
	}
}

type openRouterChatRequest struct {
	Model     string                  `json:"model"`
	Messages  []openRouterChatMessage `json:"messages"`
	Reasoning openRouterReasoning     `json:"reasoning"`
}

type openRouterChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openRouterReasoning struct {
	Enabled bool `json:"enabled"`
}

type openRouterChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (this *OpenRouterGatewayImpl) ExtractArticleMetadata(articleText string, studyFields []entity.StudyFieldCatalogItem, langs []string) (*entity.ArticleMetadataExtraction, error) {
	catalog, err := json.Marshal(studyFields)
	if err != nil {
		return nil, response.ExternalServiceError
	}

	payload, err := json.Marshal(openRouterChatRequest{
		Model: openRouterExtractModel,
		Messages: []openRouterChatMessage{
			{Role: "system", Content: articleMetadataSystemPrompt(langs)},
			{Role: "user", Content: articleMetadataUserPrompt(articleText, string(catalog), langs)},
		},
		Reasoning: openRouterReasoning{Enabled: true},
	})
	if err != nil {
		return nil, response.ExternalServiceError
	}

	req, err := http.NewRequest(http.MethodPost, this.baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {

		return nil, response.ExternalServiceError
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+this.apiKey)

	resp, err := this.client.Do(req)
	if err != nil {
		return nil, response.ExternalServiceError
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, response.ExternalServiceError
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("OpenRouter Error [Status: %d]: %s\n", resp.StatusCode, string(body))
		return nil, response.ExternalServiceError
	}

	var parsed openRouterChatResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, response.ExternalServiceError
	}
	if len(parsed.Choices) == 0 {
		return nil, response.ExternalServiceError
	}

	raw := stripJSONContent(parsed.Choices[0].Message.Content)
	result := entity.EmptyArticleMetadataExtraction(langs)
	if err := json.Unmarshal([]byte(raw), result); err != nil {
		return nil, fmt.Errorf("%w: %v", response.ExternalServiceError, err)
	}
	result.Normalize(langs)
	return result, nil
}
