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

func Authenticate() Token { // Should be run on a cache timer or on low demand requests
	TraceLog("Authenticating...")

	newtoken := GetAccessToken()
	if newtoken == (Token{}) {
		ErrorLog("Failed to get access token")
		return Token{}
	} else if newtoken != (Token{}) {
		if newtoken.AccessToken == "" {
			ErrorLog("Failed to get access token")
			return Token{}
		} else {
			TraceLog("Retrieved new access token")
			return newtoken
		}
	}
	return Token{}
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
