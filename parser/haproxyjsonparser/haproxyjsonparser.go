package haproxyjsonparser

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// JsonParser parse a JSON string.
type HaproxyJsonParser struct {
	regexPrefix *regexp.Regexp
}

// NewJsonParser returns a new json parser.
func NewHaproxyJsonParser() *HaproxyJsonParser {
	regexPrefix := regexp.MustCompile("^.+{")
	return &HaproxyJsonParser{
		regexPrefix: regexPrefix,
	}
}

// ParseString implements the Parser interface.
// The value in the map is not necessarily a string, so it needs to be converted.
func (j *HaproxyJsonParser) ParseString(line string) (map[string]string, error) {
	result := map[string]string{
		"request_length": "",
		"time_local":     "",
	}

	line = j.regexPrefix.ReplaceAllString(line, "{")
	parsed, err := j.parseJSON(line)
	if err != nil {
		return nil, err
	}

	// request
	if v, ok := parsed["header"]; ok {
		result["request"] = v
	} else {
		return nil, errors.New("Error parse request")
	}

	// status code
	if v, ok := parsed["status"]; ok {
		result["status"] = v
	} else {
		return nil, errors.New("Error parse status code")
	}

	// body bytes sent
	if v, ok := parsed["bytes"]; ok {
		result["body_bytes_sent"] = v
	} else {
		return nil, errors.New("Error parse body bytes sent")
	}

	// request method
	splitRequestStr := strings.Split(result["request"], " ")
	if len(splitRequestStr) == 0 {
		return nil, errors.New("Error parse request method")
	}
	result["request_method"] = splitRequestStr[0]

	// request_time
	if v, ok := parsed["TR/Tw/Tc/Tr/Ta"]; ok {
		timeSplit := strings.Split(v, "/")
		if len(timeSplit) != 5 {
			return nil, errors.New("Error parse request time")
		}
		intUpstreamTime, err := strconv.Atoi(timeSplit[3])
		if err != nil {
			return nil, errors.New("Error parse request time")
		}
		intRequestTime, err := strconv.Atoi(timeSplit[4])
		if err != nil {
			return nil, errors.New("Error parse request time")
		}
		// convert to second format
		result["upstream_response_time"] = fmt.Sprintf("%.3f", float64(intUpstreamTime)/1000)
		result["request_time"] = fmt.Sprintf("%.3f", float64(intRequestTime)/1000)
	} else {
		return nil, errors.New("Error parse request time")
	}

	return result, nil
}

func (j *HaproxyJsonParser) parseJSON(line string) (map[string]string, error) {
	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(line), &parsed)
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
