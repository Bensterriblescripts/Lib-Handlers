package guid

import (
	"strconv"
	"strings"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

type Guid struct {
	String string
	Valid  bool
}

func New(value string) Guid {
	cleanguid := strings.TrimSpace(value)
	cleanguid = strings.ToLower(cleanguid)
	guid := Guid{
		cleanguid,
		false,
	}
	if len(cleanguid) != 36 {
		ErrorLog("Invalid GUID Length (" + strconv.Itoa(len(cleanguid)) + ") || " + cleanguid)
		return guid
	}
	if cleanguid[8] != '-' || cleanguid[13] != '-' || cleanguid[18] != '-' || cleanguid[23] != '-' {
		ErrorLog("Invalid GUID Format || " + cleanguid)
		return guid
	}

	guid.Valid = true
	return guid
}
func Empty() Guid {
	return Guid{
		"",
		false,
	}
}
func Matches(val Guid, otherval string) bool {
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
