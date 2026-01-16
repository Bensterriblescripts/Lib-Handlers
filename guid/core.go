package guid

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

type Guid struct {
	String string
	Valid  bool
}

// # Create a new Guid from a String
//
// primarily used to clean a string into a valid Guid or store a new Guid
//
// passing an empty string will create a 0-value guid.String with a false guid.Valid value
//
// if the guid passed in is invalid, the resulting guid.Valid will be false
func New(value ...string) Guid {
	guid := Guid{
		"00000000-0000-0000-0000-000000000000",
		false,
	}

	if len(value) == 0 || value[0] == "" {
		return guid
	} else if len(value) == 2 {
		ErrorLog("Too many arguments passed into guid.New")
		return guid
	} else {
		guid.String = strings.TrimSpace(value[0])
		guid.String = strings.ToLower(guid.String)
	}

	if len(guid.String) != 36 {
		ErrorLog("Invalid GUID Length (" + strconv.Itoa(len(guid.String)) + ") || " + guid.String)
		return guid
	}
	if guid.String[8] != '-' || guid.String[13] != '-' || guid.String[18] != '-' || guid.String[23] != '-' {
		ErrorLog("Invalid GUID Format || " + guid.String)
		return guid
	}

	guid.Valid = true
	return guid
}

// # Check if a Guid matches a String
//
// Note: guid.String can always be compared against another guid.String with == due to the formatting enforced in guid.New()
func MatchesString(val Guid, otherval string) bool {
	if !val.Valid {
		ErrorLog("First parameter passed into guid.MatchesString is invalid " + val.String)
		return false
	}
	if otherval == "" {
		ErrorLog("Second parameter passed into guid.MatchesString is empty")
		return false
	}

	return strings.EqualFold(val.String, otherval)
}

func (g *Guid) UnmarshalJSON(b []byte) error { // Automatically called when unmarshalling JSON
	if bytes.Equal(b, []byte("null")) {
		*g = Guid{}
		return nil
	}

	var s string
	if ErrExists(json.Unmarshal(b, &s)) {
		return nil
	}

	if s == "" {
		*g = Guid{}
		return nil
	}

	*g = New(s)
	return nil
}
