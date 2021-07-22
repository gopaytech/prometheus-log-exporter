package kubeparser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKubeParse(t *testing.T) {
	parser := NewKubeParser(`[$time_local] $request_method "$request" $request_length $body_bytes_sent $status $request_time $upstream_response_time`)
	line := `{"log": "[03/Feb/2021:11:22:33 +0800] GET \"GET /order/2145 HTTP/1.1\" 123 518 200 0.544 0.543", "stream":"stdout","time":"2021-07-21T07:26:39.102952491Z"}`

	got, err := parser.ParseString(line)
	require.NoError(t, err)

	want := map[string]string{
		"time_local":             "03/Feb/2021:11:22:33 +0800",
		"request_time":           "0.544",
		"request_length":         "123",
		"upstream_response_time": "0.543",
		"status":                 "200",
		"body_bytes_sent":        "518",
		"request":                "GET /order/2145 HTTP/1.1",
		"request_method":         "GET",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("KubeParser.Parse() = \n%v\n%v", got, want)
	}
}

func BenchmarkParseKube(b *testing.B) {
	parser := NewKubeParser(`[$time_local] $request_method "$request" $request_length $body_bytes_sent $status $request_time $upstream_response_time`)
	line := `{"log": "[03/Feb/2021:11:22:33 +0800] GET \"GET /order/2145 HTTP/1.1\" 123 518 200 0.544 0.543", "stream":"stdout","time":"2021-07-21T07:26:39.102952491Z"}`

	for i := 0; i < b.N; i++ {
		res, err := parser.ParseString(line)
		if err != nil {
			b.Error(err)
		}
		_ = fmt.Sprintf("%v", res)
	}
}
