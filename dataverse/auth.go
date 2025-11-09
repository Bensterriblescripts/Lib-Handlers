package dataverse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Bensterriblescripts/Lib-Handlers/guid"
	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

type Token struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

var ClientID guid.Guid
var ClientSecret string
var TenantID guid.Guid
var Endpoint string

var AccessToken Token

func Authenticate() { // Should be run on a cache timer or on low demand requests
	TraceLog("Authenticating...")

	newtoken := GetAccessToken()
	if newtoken == (Token{}) {
		ErrorLog("Failed to get access token")
		return
	}

	if newtoken != (Token{}) {
		TraceLog("Retrieved new access token")
		AccessToken = newtoken
	}
}
func GetAccessToken() Token {
	var token Token
	tokenurl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", TenantID.String)
	data := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {ClientID.String},
		"client_secret": {ClientSecret},
		"scope":         {Endpoint},
	}

	resp := PanicError(http.PostForm(tokenurl, data))
	defer resp.Body.Close()
	body := PanicError(io.ReadAll(resp.Body))
	if ErrExists(json.Unmarshal(body, &token)) {
		ErrorLog("Failed to deserialize response")
		return Token{}
	}
	if token.AccessToken == "" {
		ErrorLog("Failed to get access token")
		return Token{}
	} else {
		TraceLog("Retrieved Azure access token")
		return token
	}
}
func Request(table string, params string) []byte {
	debugrequests := true // Show full request/response body and headers
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

			if debugrequests {
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
				if debugrequests {
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
	debugrequests := true // Show full request/response body and headers
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

			if debugrequests {
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
				if debugrequests {
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
