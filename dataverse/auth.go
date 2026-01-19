package dataverse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Bensterriblescripts/Lib-Handlers/guid"
	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
	"github.com/Bensterriblescripts/Lib-Handlers/time"
)

type Token struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var VerboseLogging bool = true

var ClientID guid.Guid
var ClientSecret string
var TenantID guid.Guid
var Endpoint string

var CurrentAccessToken Token
var CurrentAccessExpires int

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
// - Optional (Default = true): dataverse.VerboseLogging (Bool, Verbose Logging)
//
// Returns the current access token, as well as storing it in dataverse.CurrentAccessToken
func Authenticate() Token {
	TraceLog("Retrieving access token...")

	newtoken := getAccessToken()
	if newtoken == (Token{}) {
		ErrorLog("Failed to get access token from Dataverse")
		return newtoken
	} else {
		TraceLog("Access token retrieved from Dataverse")
		if VerboseLogging {
			TraceLog("Access token expires in: " + strconv.Itoa(newtoken.ExpiresIn) + " seconds")
		}
		CurrentAccessToken = newtoken
		CurrentAccessExpires = int(time.GetUnixTime()) + newtoken.ExpiresIn
		return CurrentAccessToken
	}
}

// Checks if the access token is valid
//
// If empty or not valid, it will attempt to retrieve a new one.
//
// Returns true if valid or reauthenticated. False if not.
func IsAuthenticated() bool {
	if CurrentAccessToken == (Token{}) {
		TraceLog("Dataverse access token was not set")
		return false
	} else if CurrentAccessToken.AccessToken == "" {
		TraceLog("Dataverse access token was empty")
		return false
	} else if CurrentAccessExpires < int(time.GetUnixTime()) {
		TraceLog("Dataverse access token has expired")
		return false
	}

	if VerboseLogging {
		TraceLog("Dataverse access token passed validation")
	}
	return true
}

// Checks if the access token is valid
//
// If empty or not valid, it will attempt to retrieve a new one.
//
// Returns true if valid or reauthenticated. False if not.
func EnsureAuthenticated() bool {
	if !IsAuthenticated() {
		Authenticate()
		if !IsAuthenticated() {
			ErrorLog("Failed to reauthenticate with Dataverse")
			return false
		} else {
			TraceLog("Reauthenticated with Dataverse")
			return true
		}
	}
	return true
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

	if VerboseLogging {
		TraceLog("Sending request to: " + tokenurl)
		TraceLog("Data: " + data.Encode())
		TraceLog("----------")
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
		if VerboseLogging {
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
