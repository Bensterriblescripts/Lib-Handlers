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

// Retrieve records by parameter/s
// E.g. dataverse.Request("contacts", "idnumber eq 1234567890", "fullname,lastname,emailaddress1", "lastname asc") -> []byte JSON
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

// Retrieve records using your own query
// Enter the entire url, this will not build for you
func RetrieveQuery(table string, query string) []byte {
	url := fmt.Sprintf("%s/api/data/v9.2/%s%s", Endpoint, table, query)
	return sendRequest(url, "GET", nil)
}

// Retrieve record by primary key
// E.g. dataverse.Retrieve("contacts", "1234567890", "fullname,lastname,emailaddress1") -> []byte JSON
func RetrieveByID(table string, primaryKey string, returnValues string) []byte {
	base := fmt.Sprintf("%s/api/data/v9.2/%s(%s)", Endpoint, table, primaryKey)

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
// E.g. dataverse.Create("contacts", jsondata) -> bool
func Create(table string, createData map[string]interface{}) bool {
	if table == "" {
		ErrorLog("Empty table found when trying to create record")
		return false
	}
	if createData == nil {
		ErrorLog("Empty data found when trying to create record")
		return false
	}
	url := fmt.Sprintf("%s/api/data/v9.2/%s", Endpoint, table)
	if jsonData, err := ErrorExists(json.Marshal(createData)); err {
		ErrorLog("Failed to marshal create data: " + string(jsonData))
		return false
	} else {
		if response := sendRequest(url, "POST", jsonData); response == nil {
			ErrorLog("Failed to create record: " + url)
			return false
		}
	}
	return true
}

// Update a table record by ID
//
// Data should be a map[string (column name)]interface{}("new value")
//
// E.g. Update("contacts", "000000-000-0000-000000", map[string]interface{}{"mito_totarasync": nil})
// Or create the map first | newData := map[string]interface{} | and append it
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
		if AllowAnnotations {
			req.Header.Set("Prefer", `odata.include-annotations="*"`)
			// req.Header.Set("Prefer", `odata.include-annotations="OData.Community.Display.V1.FormattedValue"`)
		}

		client := &http.Client{}
		if resp, err := ErrorExists(client.Do(req)); err { // Send request
			ErrorLog("CRM Request Failed: " + url)
			return nil
		} else {
			defer resp.Body.Close()

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
				ErrorLog("Failed to read response body")
				return nil
			} else {
				if VerboseLogging {
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
