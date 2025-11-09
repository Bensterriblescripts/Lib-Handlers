package dataverse

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

func Request(table string, params string) []byte {
	token := AccessToken.AccessToken
	if token == "" {
		Authenticate()
		if AccessToken == (Token{}) {
			ErrorLog("Failed to get access token")
			return nil
		}
	}

	url := fmt.Sprintf("%s/api/data/v9.2/%s", Endpoint, table)
	if params != "" {
		url += "?" + params
	}

	if req, err := ErrorExists(http.NewRequest("GET", url, nil)); err { // Create request
		ErrorLog("Failed to create request: " + url)
		return nil
	} else {
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("OData-MaxVersion", "4.0")
		req.Header.Set("OData-Version", "4.0")

		client := &http.Client{}
		if resp, err := ErrorExists(client.Do(req)); err { // Send request
			ErrorLog("CRM Request Failed: " + url)
			return nil
		} else {
			defer resp.Body.Close()

			if NetworkDebug {
				TraceLog("Sending request...  " + url + " HTTP Status: " + strconv.Itoa(resp.StatusCode))
				TraceLog("----------")
				TraceLog("Response Headers:")
				for key, value := range resp.Header {
					TraceLog(fmt.Sprintf("%s: %s", key, value))
				}
				TraceLog("----------")
			}

			if body, err := ErrorExists(io.ReadAll(resp.Body)); err { // Read response
				ErrorLog("Failed to read response body")
				return nil
			} else {
				if NetworkDebug {
					TraceLog("Response body:")
					TraceLog(string(body))
					TraceLog("----------")
				}
				if resp.StatusCode != 200 {
					ErrorLog(fmt.Sprintf("HTTP Error %d: %s", resp.StatusCode, string(body)))
					return nil
				}
				return body
			}
		}
	}
}
func Create(table string, data []byte) []byte {
	token := AccessToken.AccessToken
	if token == "" {
		Authenticate()
		if AccessToken == (Token{}) {
			ErrorLog("Failed to get access token")
			return nil
		}
	}

	url := fmt.Sprintf("%s/api/data/v9.2/%s", Endpoint, table)

	if req, err := ErrorExists(http.NewRequest("POST", url, bytes.NewBuffer(data))); err { // Create request
		ErrorLog("Failed to create request: " + url)
		return nil
	} else {
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("OData-MaxVersion", "4.0")
		req.Header.Set("OData-Version", "4.0")

		client := &http.Client{}
		if resp, err := ErrorExists(client.Do(req)); err { // Send request
			ErrorLog("CRM Request Failed: " + url)
			return nil
		} else {
			defer resp.Body.Close()

			if NetworkDebug {
				TraceLog("Sending request...  " + url + " HTTP Status: " + strconv.Itoa(resp.StatusCode))
				TraceLog("----------")
				TraceLog("Response Headers:")
				for key, value := range resp.Header {
					TraceLog(fmt.Sprintf("%s: %s", key, value))
				}
				TraceLog("----------")
			}

			if body, err := ErrorExists(io.ReadAll(resp.Body)); err { // Read response
				ErrorLog("Failed to read response body")
				return nil
			} else {
				if NetworkDebug {
					TraceLog("Response body:")
					TraceLog(string(body))
					TraceLog("----------")
				}
				if resp.StatusCode != 200 {
					ErrorLog(fmt.Sprintf("HTTP Error %d: %s", resp.StatusCode, string(body)))
					return nil
				}
				return body
			}
		}
	}
}
