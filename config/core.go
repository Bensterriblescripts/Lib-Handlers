package config

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
	osapi "github.com/Bensterriblescripts/Lib-Handlers/os"
)

var CurrentConfiguration map[string]string

func GetConfig() []byte {
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

func ReadConfig(config []byte) map[string]string {
	out := make(map[string]string)
	s := bufio.NewScanner(bytes.NewReader(config))
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
	return out
}
func OverwriteConfig(newconfigmap map[string]string) {
	if !osapi.EnsurePath(ConfigPath) {
		ErrorLog("Failed to create the config directory")
		return
	}

	currentconfig := GetConfig()
	if currentconfig == nil {
		if !WriteConfig(ConfigPath, newconfigmap) {
			ErrorLog("Failed to write the config file")
			return
		} else {
			TraceLog("Wrote to the config file")
		}
	}
	currentconfigmap := ReadConfig(currentconfig)
	for key, _ := range currentconfigmap {
		if _, exists := newconfigmap[key]; exists {
			currentconfigmap[key] = newconfigmap[key]
		} else {
			delete(currentconfigmap, key)
		}
	}

	if !WriteConfig(ConfigPath, currentconfigmap) {
		ErrorLog("Failed to write the config file")
		return
	} else {
		TraceLog("Wrote to the config file")
	}
}
func WriteConfig(path string, configmap map[string]string) bool {
	var buffer bytes.Buffer
	for key, value := range configmap {
		buffer.WriteString(key + "=" + value + "\n")
		if ErrExists(os.WriteFile(path, buffer.Bytes(), 0644)) {
			ErrorLog("Failed to write the config file")
			return false
		}
	}
	return true
}
