package llm

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

type PromptMessage struct {
	Role    string `json:"role"`
	Content []struct {
		Type  string `json:"type"`
		Text  string `json:"text,omitempty"`
		Image string `json:"image_url,omitempty"`
		File  string `json:"file_id,omitempty"`
	} `json:"content"`
}
type ChatRequest struct {
	Model string          `json:"model"`
	Input []PromptMessage `json:"input"`
}
type ChatResponse struct {
	Output []struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text,omitempty"`
		} `json:"content"`
	} `json:"output"`
	Error string `json:"error,omitempty"`
}

func SendRequest(textprompts []PromptMessage) string {
	if OPENAIApiKey == "" {
		ErrorLog("LLM API key has not been not set")
		return ""
	}
	if len(textprompts) == 0 {
		ErrorLog("At least one text prompt is required")
		return ""
	}
	var err bool

	/* Create the request */
	reqBody := ChatRequest{
		Model: "gpt-5-2025-08-07",
	}
	if len(textprompts) > 0 {
		for _, prompt := range textprompts { // Build the text prompts
			reqBody.Input = append(reqBody.Input, prompt)
		}
	}
	var bodyBytes []byte // Serialise as json
	if bodyBytes, err = ErrorExists(json.Marshal(reqBody)); err {
		ErrorLog("Error marshalling request body" + string(bodyBytes))
		return ""
	}

	/* Create the request */
	var req *http.Request
	if req, err = ErrorExists(http.NewRequest("POST", "https://api.openai.com/v1/responses", bytes.NewReader(bodyBytes))); err {
		ErrorLog("Error creating request" + string(bodyBytes))
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+OPENAIApiKey)
	var client = &http.Client{Timeout: 5 * time.Minute} // Yes, it does sometimes take that long

	/* Send it */
	var resp *http.Response
	if resp, err = ErrorExists(client.Do(req)); err {
		ErrorLog("Error executing request " + string(bodyBytes))
		return ""
	} else if resp == nil {
		ErrorLog("No response from the server " + string(bodyBytes))
		return ""
	}
	defer WrapErr(resp.Body.Close)

	var respBytes []byte
	if respBytes, err = ErrorExists(io.ReadAll(resp.Body)); err { // Read the response body
		ErrorLog("Error reading response body" + string(respBytes))
		return ""
	} else if resp.StatusCode != http.StatusOK {
		ErrorLog("Error response status code: " + strconv.Itoa(resp.StatusCode) + " " + string(respBytes))
		return ""
	}

	var cr ChatResponse
	if ErrExists(json.Unmarshal(respBytes, &cr)) { // Deserialise and check for errors
		ErrorLog("Error unmarshalling response" + string(respBytes))
		return ""
	} else if cr.Error != "" {
		ErrorLog("Error in response " + cr.Error)
		ErrorLog(string(respBytes))
		return ""
	} else if len(cr.Output) == 1 {
		ErrorLog("No choices in response " + string(respBytes))
		return ""
	}

	return cr.Output[1].Content[0].Text
}

type FileResponse struct {
	ID    string `json:"id"`
	Error string `json:"error,omitempty"`
}

func UploadFile(path string) string {
	if strings.TrimSpace(path) == "" {
		ErrorLog("Path is empty")
		return ""
	}
	var err bool

	var file *os.File
	if file, err = ErrorExists(os.Open(path)); err {
		ErrorLog("Error opening file: " + path)
		return ""
	} else {
		defer WrapErr(file.Close)
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if ErrExists(writer.WriteField("purpose", "user_data")) {
		ErrorLog("Error writing user_data purpose field")
		return ""
	}
	var part io.Writer
	if part, err = ErrorExists(writer.CreateFormFile("file", filepath.Base(path))); err {
		ErrorLog("Error creating form file")
		return ""
	}
	if _, err = ErrorExists(io.Copy(part, file)); err {
		ErrorLog("Error copying file: " + path)
		return ""
	}
	if ErrExists(writer.Close()) {
		ErrorLog("Error closing writer")
		return ""
	}

	var req *http.Request
	if req, err = ErrorExists(http.NewRequest("POST", "https://api.openai.com/v1/files", &buf)); err {
		ErrorLog("Error creating request")
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+OPENAIApiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var client = &http.Client{Timeout: 60 * time.Second}
	var resp *http.Response
	if resp, err = ErrorExists(client.Do(req)); err {
		ErrorLog("Error executing request")
		return ""
	} else if resp == nil {
		ErrorLog("No response from the server")
		return ""
	}
	defer WrapErr(resp.Body.Close)

	var response FileResponse
	if ErrExists(json.NewDecoder(resp.Body).Decode(&response)) {
		ErrorLog("Error decoding response body")
		return ""
	} else {
		if response.Error != "" {
			ErrorLog("Error in response: " + response.Error)
			return ""
		}
		TraceLog("File uploaded: " + response.ID)
		return response.ID
	}
}
