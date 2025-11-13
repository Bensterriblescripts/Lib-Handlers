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
	Text  struct {
		Format struct {
			Type   string         `json:"type"`
			Name   string         `json:"name"`
			Schema map[string]any `json:"schema"`
			Strict bool           `json:"strict"`
		} `json:"format"`
	} `json:"text"`
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

func HandleNewFile(path string) []byte {
	if path == "" {
		ErrorLog("Path is empty")
		return nil
	}
	// fileid := UploadFile(path)
	// if fileid == "" {
	// 	ErrorLog("Failed to upload file")
	// 	return false
	// } else {
	// 	TraceLog("Uploaded new file")
	// }

	response := GetContents("file-2fRFoyxgnASo4kHc9Zxb5d")
	if response == nil {
		ErrorLog("Failed to get contents")
		return nil
	} else {
		TraceLog("Contents: " + string(response))
	}

	response = StructureContents(response)
	if response == nil {
		ErrorLog("Failed to structure contents")
		return nil
	} else {
		TraceLog("Structured contents: " + string(response))
	}

	return response
}

func HandleFileUpload(path string) string {
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

	TraceLog("Uploading file: " + path)
	var req *http.Request
	if req, err = ErrorExists(http.NewRequest("POST", "https://api.openai.com/v1/files", &buf)); err {
		ErrorLog("Error creating request")
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
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
func GetContents(fileid string) []byte {
	response := SendRequest([]PromptMessage{
		{
			Role: "user",
			Content: []struct {
				Type  string `json:"type"`
				Text  string `json:"text,omitempty"`
				Image string `json:"image_url,omitempty"`
				File  string `json:"file_id,omitempty"`
			}{
				{
					Type: "input_text",
					Text: `From the document, create a table of contents in the format of Section, Subsection, Page Range. The page range is the start and end page of the section or subsection. Return the response as a JSON object.
					Example:
					{"toc":[{"section":"Purpose of this standard", range: "5-9", "subsections":[]},{"section":"Scope and application of this standard","subsections":[["Scope of this standard","5-6"],["Application of this standard","6-9"]]}]`,
				},
				{
					Type: "input_file",
					File: fileid,
				},
			},
		},
	}, "gpt-5-2025-08-07")
	if response == nil {
		ErrorLog("No response from the LLM")
		return nil
	}
	return response
}
func StructureContents(contents []byte) []byte {
	response := SendRequest([]PromptMessage{
		{
			Role: "user",
			Content: []struct {
				Type  string `json:"type"`
				Text  string `json:"text,omitempty"`
				Image string `json:"image_url,omitempty"`
				File  string `json:"file_id,omitempty"`
			}{
				{
					Type: "input_text",
					Text: `Double check the json object doesn't contain obvious sequential duplicates, missing page ranges or the words "Section" or "Subsection" in the section or subsection name. There is allowed to be missing subsections. If there are any issues, remove the duplicate or add the missing page range. Return the corrected json object.
					Format example:
					{"toc":[{"section":"Purpose of this standard", range: "5-9", "subsections":[]},{"section":"Scope and application of this standard","subsections":[["Scope of this standard","5-6"],["Application of this standard","6-9"]]}]`,
				},
				{
					Type: "input_text",
					Text: string(contents),
				},
			},
		},
	}, "gpt-5-mini-2025-08-07")
	if response == nil {
		ErrorLog("No response from the LLM")
		return nil
	}
	return response
}

func SendRequest(textprompts []PromptMessage, model string) []byte {
	if OpenAIApiKey == "" {
		Panic("OpenAI API key has not been set")
		return nil
	}
	if len(textprompts) == 0 {
		ErrorLog("At least one text prompt is required")
		return nil
	}
	var err bool

	/* Create the request */
	reqBody := ChatRequest{
		Model: model,
	}
	if len(textprompts) > 0 {
		for _, prompt := range textprompts { // Build the text prompts
			reqBody.Input = append(reqBody.Input, prompt)
		}
	}
	reqBody.Text.Format.Type = "json_schema"
	reqBody.Text.Format.Name = "lvvta_document_extraction"
	reqBody.Text.Format.Strict = true
	reqBody.Text.Format.Schema = map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"required":             []string{"toc"},
		"properties": map[string]any{
			"toc": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"section": map[string]any{
							"type": "string",
						},
						"range": map[string]any{
							"type": "string",
						},
						"subsections": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type":  "array",
								"items": map[string]any{"type": "string"},
							},
						},
					},
					"required":             []string{"section", "range", "subsections"},
					"additionalProperties": false,
				},
			},
		},
	}
	var bodyBytes []byte // Serialise as json
	if bodyBytes, err = ErrorExists(json.Marshal(reqBody)); err {
		ErrorLog("Error marshalling request body" + string(bodyBytes))
		return nil
	}

	/* Create the request */
	var req *http.Request
	if req, err = ErrorExists(http.NewRequest("POST", "https://api.openai.com/v1/responses", bytes.NewReader(bodyBytes))); err {
		ErrorLog("Error creating request" + string(bodyBytes))
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+OpenAIApiKey)
	var client = &http.Client{Timeout: 10 * time.Minute} // Yes, it does sometimes take that long

	/* Send it */
	TraceLog("Sending request: " + string(bodyBytes))
	var resp *http.Response
	if resp, err = ErrorExists(client.Do(req)); err {
		ErrorLog("Error executing request " + string(bodyBytes))
		return nil
	} else if resp == nil {
		ErrorLog("No response from the server " + string(bodyBytes))
		return nil
	}
	defer WrapErr(resp.Body.Close)

	var respBytes []byte
	if respBytes, err = ErrorExists(io.ReadAll(resp.Body)); err { // Read the response body
		ErrorLog("Error reading response body" + string(respBytes))
		return nil
	} else if resp.StatusCode != http.StatusOK {
		ErrorLog("Error response status code: " + strconv.Itoa(resp.StatusCode) + " " + string(respBytes))
		return nil
	}

	TraceLog("Response: " + string(respBytes))
	return respBytes
}
