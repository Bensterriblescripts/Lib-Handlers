package config

import (
	"bufio"
	"bytes"
	"maps"
	"os"
	"strings"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
	osapi "github.com/Bensterriblescripts/Lib-Handlers/osapi"
)

var ConfigPath string
var Current map[string]string
var Draft map[string]string

func ReadConfig() map[string]string {
	rawconfig := getConfig()

	out := make(map[string]string)
	s := bufio.NewScanner(bytes.NewReader(rawconfig))
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		i := strings.IndexRune(line, '=')
		if i < 0 {
			continue
		}
		k := strings.TrimSpace(line[:i])
		v := strings.TrimSpace(line[i+1:])
		if len(v) >= 2 && ((v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'')) {
			v = v[1 : len(v)-1]
		}
		out[k] = v
	}
	Current = out
	return out
}
func getConfig() []byte {
	if ConfigPath == "" {
		ConfigPath = "C:\\Local\\Config\\" + AppName + ".ini"
	}
	if !osapi.EnsurePath(ConfigPath) {
		ErrorLog("Failed to retrieve the config during directory creation")
		return []byte{}
	}

	var filecontent []byte
	var err bool
	if filecontent, err = ErrorExists(os.ReadFile(ConfigPath)); err {
		if _, err = ErrorExists(os.Create(ConfigPath)); err {
			ErrorLog("Failed to create the config file")
			return []byte{}
		} else {
			if filecontent, err = ErrorExists(os.ReadFile(ConfigPath)); err {
				return filecontent
			} else {
				if len(filecontent) > 0 {
					return filecontent
				} else {
					ErrorLog("Failed to retrieve the config during file reading")
					return []byte{}
				}
			}
		}
	} else {
		return filecontent
	}
}

func Write() { // Write config.Current to the file
	if !overwriteConfig() {
		ErrorLog("Failed to write the config file")
	} else {
		TraceLog("Wrote to the config file")
	}
}
func WriteSetting(label string, value string) { // Write single setting
	if label == "" {
		ErrorLog("Label cannot be an empty string")
		return
	}
	if value == "" {
		ErrorLog("The value passed into config.Write was empty")
		return
	}

	if Current == nil {
		_ = ReadConfig()
	}
	Current[label] = value
	if !overwriteConfig() {
		ErrorLog("Failed to write the config file")
	} else {
		TraceLog("Wrote to the config file")
	}
}
func WriteSettings(newConfig map[string]string) { // Write multiple settings
	if len(newConfig) == 0 {
		return
	}
	if Current == nil {
		_ = ReadConfig()
	}
	maps.Copy(Current, newConfig)

	if !overwriteConfig() {
		ErrorLog("Failed to write the config file")
	} else {
		TraceLog("Wrote to the config file")
	}
}
func overwriteConfig() bool {
	if !osapi.EnsurePath(ConfigPath) {
		ErrorLog("Failed to create the config directory")
		return false
	}
	if len(Current) == 0 {
		return false
	}

	var buffer bytes.Buffer
	for key, value := range Current {
		buffer.WriteString(key + "=" + value + "\n")
		if ErrExists(os.WriteFile(ConfigPath, buffer.Bytes(), 0644)) { // Truncates the file before writing
			ErrorLog("Failed to write the config file")
			return false
		}
	}
	return true
}
