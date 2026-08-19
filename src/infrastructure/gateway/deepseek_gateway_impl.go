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

const deepSeekTimeout = 180 * time.Second

type DeepSeekGatewayImpl struct {
	client  *http.Client
	baseURL string
	apiKey  string
	model   string
}

// @inject
func NewDeepSeekGateway(env *config.Config) gateway.DeepSeekGateway {
	model := strings.TrimSpace(env.DeepSeekModel)
	if model == "" {
		model = "deepseek-v4-flash"
	}
	return &DeepSeekGatewayImpl{
		client:  &http.Client{Timeout: deepSeekTimeout},
		baseURL: strings.TrimRight(env.DeepSeekBaseURL, "/"),
		apiKey:  env.DeepSeekAPIKey,
		model:   model,
	}
}

type deepSeekChatRequest struct {
	Model          string                 `json:"model"`
	Messages       []deepSeekChatMessage  `json:"messages"`
	ResponseFormat deepSeekResponseFormat `json:"response_format"`
	Thinking       deepSeekThinking       `json:"thinking"`
}

type deepSeekChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type deepSeekResponseFormat struct {
	Type string `json:"type"`
}

type deepSeekThinking struct {
	Type string `json:"type"`
}

type deepSeekChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (this *DeepSeekGatewayImpl) ExtractArticleMetadata(articleText string, studyFields []entity.StudyFieldCatalogItem, langs []string) (*entity.ArticleMetadataExtraction, error) {
	catalog, err := json.Marshal(studyFields)
	if err != nil {
		return nil, response.ExternalServiceError
	}

	payload, err := json.Marshal(deepSeekChatRequest{
		Model: this.model,
		Messages: []deepSeekChatMessage{
			{Role: "system", Content: articleMetadataSystemPrompt(langs)},
			{Role: "user", Content: articleMetadataUserPrompt(articleText, string(catalog), langs)},
		},
		ResponseFormat: deepSeekResponseFormat{Type: "json_object"},
		Thinking:       deepSeekThinking{Type: "disabled"},
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
		return nil, response.ExternalServiceError
	}

	var parsed deepSeekChatResponse
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

func stripJSONContent(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```json")
		s = strings.TrimPrefix(s, "```JSON")
		s = strings.TrimPrefix(s, "```")
		if i := strings.LastIndex(s, "```"); i >= 0 {
			s = s[:i]
		}
		s = strings.TrimSpace(s)
	}
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}

func articleMetadataSystemPrompt(langs []string) string {
	keys := langObjectExample(langs)
	tagKeys := langArrayExample(langs)
	return `You extract scholarly article metadata from PDF text.
Return ONLY a JSON object with this exact shape:
{
  "article_name": ` + keys + `,
  "article_language": "",
  "study_field_ids": [],
  "doi": "",
  "tags": ` + tagKeys + `,
  "references": [],
  "annotation": ` + keys + `
}
Rules:
- Use only the provided article text. Do not invent facts that are not present.
- article_name, tags, and annotation MUST include every listed language key. If the source has only one language, translate the values into the other keys.
- article_language must be the provided article language code.
- study_field_ids must be chosen only from the provided catalog ids. If none match, use [].
- doi is a full DOI URL if possible, otherwise the raw DOI, otherwise "".
- references is an array of bibliography strings. If none, [].
- Missing values must be empty strings or empty arrays, never omit keys.
- Output JSON only.`
}

func articleMetadataUserPrompt(articleText, catalogJSON string, langs []string) string {
	return "Required language keys: " + strings.Join(langs, ", ") +
		"\nStudy field catalog (JSON):\n" + catalogJSON +
		"\n\nArticle text:\n" + articleText
}

func langObjectExample(langs []string) string {
	if len(langs) == 0 {
		langs = []string{"uz", "ru", "en"}
	}
	parts := make([]string, 0, len(langs))
	for _, lang := range langs {
		parts = append(parts, `"`+lang+`": ""`)
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

func langArrayExample(langs []string) string {
	if len(langs) == 0 {
		langs = []string{"uz", "ru", "en"}
	}
	parts := make([]string, 0, len(langs))
	for _, lang := range langs {
		parts = append(parts, `"`+lang+`": []`)
	}
	return "{" + strings.Join(parts, ", ") + "}"
}
