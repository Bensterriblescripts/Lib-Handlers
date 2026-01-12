package dataverse

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

// Retrieve records by parameter/s
// E.g. dataverse.Request("contacts", "idnumber eq 1234567890", "fullname,lastname,emailaddress1") -> []byte JSON
func Retrieve(table, filter, returnValues string, order ...string) []byte {
	base := fmt.Sprintf("%s/api/data/v9.2/%s", Endpoint, table)

	if u, err := ErrorExists(url.Parse(base)); err {
		return nil
	} else {
		q := u.Query()
		q.Set("$filter", filter)                                  // Parameters
		q.Set("$select", returnValues)                            // Return values
		if len(order) == 1 && strings.TrimSpace(order[0]) != "" { // Order by
			q.Set("$orderby", order[0])
		} else if len(order) > 1 {
			q.Set("$orderby", strings.Join(order, ","))
		}

		u.RawQuery = q.Encode()
		return sendRequest(u.String(), "GET", nil)
	}
}

// Retrieve record by primary key
// E.g. dataverse.Retrieve("contacts", "1234567890", "fullname,lastname,emailaddress1") -> []byte JSON
func RetrieveByID(table string, id string, returnValues string) []byte {
	base := fmt.Sprintf("%s/api/data/v9.2/%s(%s)", Endpoint, table, id)

	if u, err := ErrorExists(url.Parse(base)); err {
		return nil
	} else {
		q := u.Query()
		q.Set("$select", returnValues) // Parameters
		u.RawQuery = q.Encode()
		return sendRequest(u.String(), "GET", nil)
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
// E.g. dataverse.Create("contacts", jsondata) -> []byte JSON
func Create(table string, data []byte) []byte {
	url := fmt.Sprintf("%s/api/data/v9.2/%s", Endpoint, table)
	return sendRequest(url, "POST", data)
}

func sendRequest(url string, method string, data []byte) []byte {
	if CurrentAccessToken == (Token{}) || CurrentAccessToken.AccessToken == "" {
		Authenticate()
		if CurrentAccessToken == (Token{}) || CurrentAccessToken.AccessToken == "" {
			ErrorLog("Failed to retrieve access token")
			return nil
		}
	}

	if req, err := ErrorExists(http.NewRequest(method, url, bytes.NewBuffer(data))); err { // Create request
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

			if DataverseDebug {
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
				if DataverseDebug {
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
