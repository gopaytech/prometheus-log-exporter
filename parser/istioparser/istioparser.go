package istioparser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// IstioParser parse a JSON string.
type IstioParser struct {
	fnGetKubeLog func(string) string
}

// NewIstioParser returns a new json parser.
func NewIstioParser() *IstioParser {
	return &IstioParser{}
}

func NewIstioCRIParser() *IstioParser {
	return &IstioParser{
		fnGetKubeLog: func(s string) string {
			// timestamp stdout [FP] actualLog
			return strings.SplitN(s, " ", 4)[3]
		},
	}
}

// kubernetes wrap the actual log line, need to unwrap it to get the actual
func (p *IstioParser) getKubeLogLine(line string) string {
	if p.fnGetKubeLog == nil {
		p.fnGetKubeLog = func(s string) string {
			line := gjson.Get(s, "log").String()
			return strings.TrimSpace(line)
		}
	}

	return p.fnGetKubeLog(line)
}

// ParseString implements the Parser interface.
// The value in the map is not necessarily a string, so it needs to be converted.
func (p *IstioParser) ParseString(line string) (map[string]string, error) {
	actualLogLine := p.getKubeLogLine(line)

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

	requestTime, err := strconv.ParseFloat(p.GetValue(fields, "duration", ""), 64)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse request_time, value: %+v", p.GetValue(fields, "duration", ""))
	}

	upstreamRequestTime, err := strconv.ParseFloat(p.GetValue(fields, "upstream_service_time", ""), 64)
	if err != nil {
		upstreamRequestTime = 0
	}

	result := map[string]string{
		"body_bytes_sent": p.GetValue(fields, "bytes_sent", ""),
		"request": fmt.Sprintf(
			"%s %s %s",
			p.GetValue(fields, "method", ""),
			p.GetValue(fields, "path", ""),
			p.GetValue(fields, "protocol", ""),
		),
		"request_length":         p.GetValue(fields, "bytes_received", ""),
		"request_method":         p.GetValue(fields, "method", ""),
		"request_time":           strconv.FormatFloat(requestTime/1000, 'f', 3, 64),         // divide by 1000, convert from ms to s
		"upstream_response_time": strconv.FormatFloat(upstreamRequestTime/1000, 'f', 3, 64), // divide by 1000, convert from ms to s
		"status":                 p.GetValue(fields, "response_code", ""),
		"time_local":             p.GetValue(fields, "start_time", ""),
		"upstream_cluster":       p.GetUpstreamCluster(fields),
		"authority":              p.GetValue(fields, "authority", ""),
	}
	return result, nil
}

func (p *IstioParser) GetUpstreamCluster(fields map[string]string) string {
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

func (p *IstioParser) GetValue(fields map[string]string, key, defaultValue string) string {
	if v, ok := fields[key]; ok {
		return v
	} else {
		return defaultValue
	}
}
