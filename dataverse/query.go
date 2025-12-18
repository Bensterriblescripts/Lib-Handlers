package dataverse

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

// Requires an existing access token in dataverse.CurrentAccessToken
//
// - Use dataverse.Authenticate() to retrieve a new token
//
// Usage: dataverse.Request("contacts", "contactid eq 1234567890")
func Request(table string, params string, returnvalues string) []byte {
	url := fmt.Sprintf("%s/api/data/v9.2/%s", Endpoint, table)
	if returnvalues != "" && params == "" {
		url += "?$select=" + returnvalues
	} else if returnvalues != "" && params != "" {
		url += "?" + params + "&$select=" + returnvalues
	} else if returnvalues == "" && params != "" {
		url += "?" + params
	}

	if req, err := ErrorExists(http.NewRequest("GET", url, nil)); err { // Create request
		ErrorLog("Failed to create request: " + url)
		return nil
	} else {
		req.Header.Set("Authorization", "Bearer "+CurrentAccessToken.AccessToken)
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

// Requires an existing access token in dataverse.CurrentAccessToken
//
// - Use dataverse.Authenticate() to retrieve a new token
//
// Usage: dataverse.Retrieve("contacts", "1234567890", "fullname,lastname,emailaddress1")
func Retrieve(table string, id string, returnvalues string) []byte {
	url := fmt.Sprintf("%s/api/data/v9.2/%s(%s)", Endpoint, table, id)
	if returnvalues != "" {
		url += "?$select=" + returnvalues
	}

	if req, err := ErrorExists(http.NewRequest("GET", url, nil)); err { // Create request
		ErrorLog("Failed to create request: " + url)
		return nil
	} else {
		req.Header.Set("Authorization", "Bearer "+CurrentAccessToken.AccessToken)
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

// Add a record into an existing table.
//
// Data should be a marshalled JSON object. E.g.
//
//	type Change struct {
//			Datestring   string `json:"cr244_datestring"`
//			Author       string `json:"cr244_author"`
//	}
//
// dataverse.Create("contacts", data)
func Create(table string, data []byte) []byte {
	url := fmt.Sprintf("%s/api/data/v9.2/%s", Endpoint, table)

	if req, err := ErrorExists(http.NewRequest("POST", url, bytes.NewBuffer(data))); err { // Create request
		ErrorLog("Failed to create request: " + url)
		return nil
	} else {
		req.Header.Set("Authorization", "Bearer "+CurrentAccessToken.AccessToken)
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
				if (resp.StatusCode-200) > 99 && (resp.StatusCode-200) < 0 { // Dataverse sends back different 200 status codes for some unholy reason
					ErrorLog(fmt.Sprintf("HTTP Error %d: %s", resp.StatusCode, string(body)))
					return nil
				}
				return body
			}
		}
	}
}
