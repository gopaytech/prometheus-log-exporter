package jsonparser

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

// JsonParser parse a JSON string.
type JsonParser struct {
	fnGetActualLog func(string) string
}

func NewJsonParser() *JsonParser {
	return &JsonParser{}
}

func NewKubeJsonParser() *JsonParser {
	return &JsonParser{
		fnGetActualLog: func(line string) string {
			// the actual k8s log is wrapped in log field
			line = gjson.Get(line, "log").String()
			return strings.TrimSpace(line)
		},
	}
}

// ParseString implements the Parser interface.
// The value in the map is not necessarily a string, so it needs to be converted.
func (j *JsonParser) ParseString(line string) (map[string]string, error) {
	var parsed map[string]interface{}
	actualLogLine := line
	if j.fnGetActualLog != nil {
		actualLogLine = j.fnGetActualLog(line)
	}
	err := json.Unmarshal([]byte(actualLogLine), &parsed)
	if err != nil {
		return nil, fmt.Errorf("json log parsing err: %w", err)
	}

	fields := make(map[string]string, len(parsed))
	for k, v := range parsed {
		if s, ok := v.(string); ok {
			fields[k] = s
		} else {
			fields[k] = fmt.Sprintf("%v", v)
		}
	}
	return fields, nil
}
