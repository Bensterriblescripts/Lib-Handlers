package dataverse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Bensterriblescripts/Lib-Handlers/guid"
	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

type Token struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

var DataverseDebug bool = false

var ClientID guid.Guid
var ClientSecret string
var TenantID guid.Guid
var Endpoint string

var CurrentAccessToken Token

// Authenticate with the Azure and retrieve a new access token
//
// Requires the following variables to be set:
//
// - dataverse.ClientID (guid.Guid)
//
// - dataverse.ClientSecret (string)
//
// - dataverse.TenantID (guid.Guid)
//
// - dataverse.Endpoint (string)
//
// - Optional (Default = false): dataverse.DataverseDebug (Bool, Verbose Logging)
//
// Returns the current access token, as well as storing it in dataverse.CurrentAccessToken
func Authenticate() Token {
	TraceLog("Authenticating...")

	newtoken := getAccessToken()
	if newtoken == (Token{}) {
		ErrorLog("Failed to get access token")
		return newtoken
	} else {
		TraceLog("Retrieved new access token")
		CurrentAccessToken = newtoken
		return CurrentAccessToken
	}
}
func getAccessToken() Token {
	var token Token
	tokenurl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", TenantID.String)
	data := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {ClientID.String},
		"client_secret": {ClientSecret},
		"scope":         {Endpoint + "/.default"},
	}

	if DataverseDebug {
		TraceLog("Sending request to: " + tokenurl)
		TraceLog("Data: " + data.Encode())
	}
	if resp, err := ErrorExists(http.PostForm(tokenurl, data)); err {
		if resp.Body != nil {
			defer WrapErr(resp.Body.Close)
		} else {
			ErrorLog("Response body is empty")
		}
		return Token{}
	} else {
		if resp.Body != nil {
			defer WrapErr(resp.Body.Close)
		} else {
			ErrorLog("Response body is empty")
			return Token{}
		}
		if DataverseDebug {
			TraceLog("Response: " + resp.Status)
			TraceLog("Response Headers:")
			for key, value := range resp.Header {
				TraceLog(fmt.Sprintf("%s: %s", key, value))
			}
		}

		if body, err := ErrorExists(io.ReadAll(resp.Body)); err {
			ErrorLog("Failed to read response body")
			return Token{}
		} else {
			if ErrExists(json.Unmarshal(body, &token)) {
				ErrorLog("Failed to deserialize response " + string(body))
				return Token{}
			}
			TraceLog("Retrieved Azure access token")
			return token
		}
	}
}
