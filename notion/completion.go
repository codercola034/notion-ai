package notion

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type getCompletionRequest struct {
	Id                string   `json:"id"`
	AiSessionId       string   `json:"aiSessionId"`
	Context           context  `json:"context"`
	Model             string   `json:"model"`
	Metadata          struct{} `json:"metadata"`
	SpaceId           string   `json:"spaceId"`
	IsSpacePermission bool     `json:"isSpacePermission"`
}

type getCompletionResult struct {
	httpRes   *http.Response
	reader    *bufio.Reader
	responses []getCompletionResponse
}

type getCompletionResponse struct {
	Type        string `json:"type"`
	InferenceId string `json:"inferenceId"`
	Completion  string `json:"completion"`
}

func (r *getCompletionResult) Output() (string, bool, error) {
	line, _, err := r.reader.ReadLine()
	if err != nil {
		if err == io.EOF {
			r.httpRes.Body.Close()
			var resultBuilder strings.Builder
			for _, v := range r.responses {
				resultBuilder.WriteString(v.Completion)
			}
			DefaultTranscripts.addHistory(transcript{Type: "assistant", Result: resultBuilder.String()})
			return resultBuilder.String(), true, nil
		}
		return "", false, err
	}

	var result getCompletionResponse
	err = json.Unmarshal(line, &result)
	if err != nil {
		return "", false, err
	}
	r.responses = append(r.responses, result)
	return result.Completion, false, nil
}

type context struct {
	Type       string       `json:"type"`
	Transcript []transcript `json:"transcript"`
}

type transcript struct {
	Type        string `json:"type"`
	PageTitle   string `json:"pageTitle,omitempty"`
	PageContent string `json:"pageContent,omitempty"`
	Prompt      string `json:"prompt,omitempty"`
	Result      string `json:"result,omitempty"`
}

func newGetCompletionRequest(query string) *getCompletionRequest {
	nt := transcript{Type: "human", Prompt: query}
	DefaultTranscripts.addHistory(nt)
	return &getCompletionRequest{
		Id:          "",
		AiSessionId: "",
		Context: context{
			Type:       "writerFollowup",
			Transcript: append(DefaultTranscripts, nt),
		},
		Model:             "default",
		Metadata:          struct{}{},
		SpaceId:           spaceId,
		IsSpacePermission: false,
	}
}

func (r getCompletionRequest) Data() []byte {
	b, _ := json.Marshal(r)
	return b
}

type Result interface {
	Output() (out string, ended bool, err error)
}

type Request interface {
	Data() []byte // return the http request data
}

func GetCompletion(query string) (Result, error) {
	httpReq, err := http.NewRequest("POST", CompletionUrl, bytes.NewReader(newGetCompletionRequest(query).Data()))
	if err != nil {
		return nil, err
	}
	tokenCookie := http.Cookie{Name: "token_v2", Value: tokenV2}
	httpReq.AddCookie(&tokenCookie)
	httpReq.Header.Add("Content-Type", "application/json")
	res, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(res.Body)
	return &getCompletionResult{
		httpRes:   res,
		reader:    reader,
		responses: []getCompletionResponse{},
	}, nil
}
