package istioparser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// IstioParser parse a JSON string.
type IstioParser struct{}

// NewIstioParser returns a new json parser.
func NewIstioParser() *IstioParser {
	return &IstioParser{}
}

// ParseString implements the Parser interface.
// The value in the map is not necessarily a string, so it needs to be converted.
func (j *IstioParser) ParseString(line string) (map[string]string, error) {
	actualLogLine := gjson.Get(line, "log").String()
	actualLogLine = strings.TrimSpace(actualLogLine)

	var parsed map[string]interface{}
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

	requestTime, err := strconv.ParseFloat(j.GetValue(fields, "duration", ""), 64)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse request_time, value: %+v", j.GetValue(fields, "duration", ""))
	}

	upstreamRequestTime, err := strconv.ParseFloat(j.GetValue(fields, "upstream_service_time", ""), 64)
	if err != nil {
		upstreamRequestTime = 0
	}

	result := map[string]string{
		"body_bytes_sent": j.GetValue(fields, "bytes_sent", ""),
		"request": fmt.Sprintf(
			"%s %s %s",
			j.GetValue(fields, "method", ""),
			j.GetValue(fields, "path", ""),
			j.GetValue(fields, "protocol", ""),
		),
		"request_length":         j.GetValue(fields, "bytes_received", ""),
		"request_method":         j.GetValue(fields, "method", ""),
		"request_time":           strconv.FormatFloat(requestTime/1000, 'f', 3, 64),         // divide by 1000, convert from ms to s
		"upstream_response_time": strconv.FormatFloat(upstreamRequestTime/1000, 'f', 3, 64), // divide by 1000, convert from ms to s
		"status":                 j.GetValue(fields, "response_code", ""),
		"time_local":             j.GetValue(fields, "start_time", ""),
		"upstream_cluster":       j.GetUpstreamCluster(fields),
	}
	return result, nil
}

func (j *IstioParser) GetUpstreamCluster(fields map[string]string) string {
	if v, ok := fields["upstream_cluster"]; ok {
		s := strings.Split(v, "|")
		if len(s) > 0 {
			return s[len(s)-1]
		} else {
			return ""
		}
	} else {
		return ""
	}
}

func (j *IstioParser) GetValue(fields map[string]string, key, defaultValue string) string {
	if v, ok := fields[key]; ok {
		return v
	} else {
		return defaultValue
	}
}
