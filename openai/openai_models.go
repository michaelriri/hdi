package openai

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type CompletionResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text  string `json:"text"`
		Index int    `json:"index"`
		// Logprobs     string `json:"logprobs, omitempty"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}
}

type CompletionRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
	N           int     `json:"n"`
	Stream      bool    `json:"stream"`
	Logprobs    int     `json:"logprobs"`
	Stop        string  `json:"stop"`
}

type OpenAI struct {
	ApiKey string
}

func NewOpenAI(apiKey string) *OpenAI {

	return &OpenAI{
		ApiKey: apiKey,
	}
}

func (o *OpenAI) Completion(request *CompletionRequest) (*CompletionResponse, error) {
	client := &http.Client{}
	url := "https://api.openai.com/v1/completions"
	encodedStr, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(encodedStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var res CompletionResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
