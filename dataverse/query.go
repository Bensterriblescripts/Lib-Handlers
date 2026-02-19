package dataverse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

var AllowAnnotations bool = true
var MaxPageSize = 5000 // Set to 0 to use dataverse defaults

// Retrieve records using your own query.
//
// Automatically appends ? if missing.
//
// Example:
//
//	body := dataverse.Query("GET", "contacts", "$select=fullname")
func Query(method string, table string, query string) []byte {
	if table == "" {
		ErrorLog("Empty table found when trying to query dataverse")
		return nil
	}
	if query == "" {
		ErrorLog("Empty query found when trying to query dataverse")
		return nil
	}

	if method == "" {
		method = "GET"
		if VerboseLogging {
			TraceLog("No method found in query, using GET")
		}
	}

	if !strings.HasPrefix(query, "?") {
		query = "?" + query
	}

	queryurl := fmt.Sprintf("%s/api/data/v9.2/%s%s", Endpoint, table, query)
	return sendRequest(queryurl, method, nil)
}

// Retrieve records by parameter/s.
//
// Example:
//
//	body := dataverse.Retrieve("contacts", "idnumber eq 1234567890", "fullname,lastname,emailaddress1", "lastname asc")
func Retrieve(table, filter, returnValues string, order ...string) []byte {
	base := fmt.Sprintf("%s/api/data/v9.2/%s", Endpoint, table)

	if u, err := ErrorExists(url.Parse(base)); err {
		return nil
	} else {
		q := u.Query()
		if filter != "" {
			q.Set("$filter", filter) // Parameters
		}
		if returnValues != "" {
			q.Set("$select", returnValues) // Return values
		}
		if len(order) == 1 && strings.TrimSpace(order[0]) != "" { // Order by
			q.Set("$orderby", order[0])
		} else if len(order) > 1 {
			q.Set("$orderby", strings.Join(order, ","))
		}

		u.RawQuery = q.Encode()
		return sendRequest(u.String(), "GET", nil)
	}
}

// Retrieve the next page of results.
//
// Example:
//
//	body := dataverse.RetrieveNext(nextLink)
func RetrieveNext(nextLink string) []byte {
	return sendRequest(nextLink, "GET", nil)
}

// Retrieve record by primary key.
//
// Example:
//
//	body := dataverse.RetrieveByID("contacts", "1234567890", "fullname,lastname,emailaddress1")
func RetrieveByID(table string, primaryKey string, returnValues string) []byte {
	base := fmt.Sprintf("%s/api/data/v9.2/%s(%s)", Endpoint, table, primaryKey)

	if u, err := ErrorExists(url.Parse(base)); err {
		return nil
	} else {
		q := u.Query()
		if returnValues != "" {
			q.Set("$select", returnValues) // Parameters
		}
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
// Example:
//
//	ok := dataverse.Create("contacts", map[string]interface{}{"fullname": "Ada Lovelace"})
func Create(table string, createData map[string]interface{}) bool {
	if table == "" {
		ErrorLog("Empty table found when trying to create record")
		return false
	}
	if createData == nil {
		ErrorLog("Empty data found when trying to create record")
		return false
	}
	createurl := fmt.Sprintf("%s/api/data/v9.2/%s", Endpoint, table)
	if jsonData, err := ErrorExists(json.Marshal(createData)); err {
		ErrorLog("Failed to marshal create data: " + string(jsonData))
		return false
	} else {
		if response := sendRequest(createurl, "POST", jsonData); response == nil {
			ErrorLog("Failed to create record: " + createurl)
			return false
		}
	}
	return true
}

// Update a table record by ID
//
// Data should be a map[string (column name)]interface{}("new value")
//
// Example:
//
//	ok := dataverse.Update("contacts", "000000-000-0000-000000", map[string]interface{}{"mito_totarasync": nil})
func Update(table string, recordid string, updateData map[string]interface{}) bool {
	if table == "" || updateData == nil || recordid == "" {
		ErrorLog("Empty parameter during update")
		ErrorLog("Table: " + table)
		ErrorLog("Update Data: " + fmt.Sprintf("%v", updateData))
		ErrorLog("Record ID: " + recordid)
		return false
	}
	baseurl := fmt.Sprintf("%s/api/data/v9.2/%s(%s)", Endpoint, table, recordid)

	if u, err := ErrorExists(url.Parse(baseurl)); err {
		ErrorLog("Failed to parse URL: " + baseurl)
		return false
	} else {
		if jsonData, err := ErrorExists(json.Marshal(updateData)); err {
			ErrorLog("Failed to marshal update data: " + string(jsonData))
			return false
		} else {
			if response := sendRequest(u.String(), "PATCH", jsonData); response == nil {
				ErrorLog("Failed to update record: " + baseurl)
				return false
			}
		}
		return true
	}
}

// sendRequest issues a request to the Dataverse API.
//
// Example:
//
//	body := sendRequest(url, "GET", nil)
func sendRequest(url string, method string, data []byte) []byte {
	if !EnsureAuthenticated() {
		ErrorLog("Failed to validate access token")
		return nil
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

		prefer := []string{}
		if MaxPageSize > 0 {
			prefer = append(prefer, "odata.maxpagesize="+strconv.Itoa(MaxPageSize)) // Pagination, required if too many results (training packages)
		}
		if AllowAnnotations {
			prefer = append(prefer, `odata.include-annotations="*"`) // Required for entity record labels (e.g. guid vs display name)
		}
		if len(prefer) > 0 {
			req.Header.Set("Prefer", strings.Join(prefer, ","))
		}

		client := &http.Client{}
		if resp, err := ErrorExists(client.Do(req)); err { // Send request
			if resp.Body != nil {
				defer WrapErr(resp.Body.Close)
			} else {
				ErrorLog("Response body is empty")
			}
			ErrorLog("CRM Request Failed: " + url)
			return nil
		} else {
			defer WrapErr(resp.Body.Close)

			if VerboseLogging {
				TraceLog("Sending request...  " + url + " HTTP Status: " + strconv.Itoa(resp.StatusCode))
				TraceLog("----------")
				TraceLog("Response Headers:")
				for key, value := range resp.Header {
					TraceLog(fmt.Sprintf("%s: %s", key, value))
				}
				TraceLog("----------")
			}

			if body, err := ErrorExists(io.ReadAll(resp.Body)); err { // Read response
				ErrorLog("Failed to read response body: " + url)
				return nil
			} else {
				if VerboseLogging {
					TraceLog("Response body:")
					TraceLog(string(body))
					TraceLog("----------")
				}
				if resp.StatusCode < 200 || resp.StatusCode >= 300 {
					ErrorLog(fmt.Sprintf("HTTP Error %d for %s: %s", resp.StatusCode, url, string(body)))
					return nil
				}
				if len(body) == 0 {
					ErrorLog("Empty response body: " + url)
					return nil
				}
				if body[0] != '{' && body[0] != '[' {
					ErrorLog(fmt.Sprintf("Non-JSON response for %s: %.200s", url, string(body)))
					return nil
				}
				if bytes.Contains(body, []byte(`"error":{`)) {
					ErrorLog(fmt.Sprintf("Dataverse error for %s: %s", url, string(body)))
					return nil
				}
				return body
			}
		}
	}
}
