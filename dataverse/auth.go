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

var NetworkDebug bool = false

var ClientID guid.Guid
var ClientSecret string
var TenantID guid.Guid
var Endpoint string

var CurrentAccessToken Token

func Authenticate() Token { // Should be run on a cache timer or on low demand requests
	TraceLog("Authenticating...")

	newtoken := GetAccessToken()
	if newtoken == (Token{}) {
		ErrorLog("Failed to get access token")
		return newtoken
	} else if newtoken != (Token{}) {
		TraceLog("Retrieved new access token")
		CurrentAccessToken = newtoken
		return newtoken
	}
	return CurrentAccessToken
}
func GetAccessToken() Token {
	var token Token
	tokenurl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", TenantID.String)
	data := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {ClientID.String},
		"client_secret": {ClientSecret},
		"scope":         {Endpoint + "/.default"},
	}

	if NetworkDebug {
		TraceLog("Sending request to: " + tokenurl)
		TraceLog("Data: " + data.Encode())
	}
	resp := PanicError(http.PostForm(tokenurl, data))
	if NetworkDebug {
		TraceLog("Response: " + resp.Status)
		TraceLog("Response Headers:")
		for key, value := range resp.Header {
			TraceLog(fmt.Sprintf("%s: %s", key, value))
		}
	}

	defer resp.Body.Close()
	body := PanicError(io.ReadAll(resp.Body))
	if ErrExists(json.Unmarshal(body, &token)) {
		ErrorLog("Failed to deserialize response " + string(body))
		return Token{}
	} else {
		TraceLog("Retrieved Azure access token")
		return token
	}
}
