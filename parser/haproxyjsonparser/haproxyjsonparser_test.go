package haproxyjsonparser

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJsonParse(t *testing.T) {
	parser := NewHaproxyJsonParser()
	line := `Nov 16 05:01:27 somehost haproxy[24334]: {"time":"16/Nov/2021:05:01:27.794", "client":"192.168.10.3:17334", "frontend":"some-frontend", "backend":"backend/some-backend", "path":"/api/v1/users", "status":200, "bytes":872, "TR/Tw/Tc/Tr/Ta":"0/0/3/31/34", "CC":"-", "CS":"-", "tsc":"----", "concurrent_connections": "7/6/0/0/0", "queues":"0/0", "header":"POST /api/v1/users HTTP/1.1"}`

	got, err := parser.ParseString(line)
	require.NoError(t, err)

	want := map[string]string{
		"time_local":             "",
		"request_time":           "0.034",
		"request_length":         "",
		"upstream_response_time": "0.031",
		"status":                 "200",
		"body_bytes_sent":        "872",
		"request":                "POST /api/v1/users HTTP/1.1",
		"request_method":         "POST",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("JsonParser.Parse() = %v, want %v", got, want)
	}
}
