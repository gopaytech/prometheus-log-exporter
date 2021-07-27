package istioparser

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIstioParse(t *testing.T) {
	parser := NewIstioParser()

	jsonLine := `{"bytes_sent":"37","upstream_cluster":"outbound|80||payment-service","downstream_remote_address":"127.0.0.1:40421","authority":"payment.example.com","path":"/v1/order/2145","protocol":"HTTP/1.1","upstream_service_time":"26","upstream_local_address":"127.0.0.1:53232","duration":"27","upstream_transport_failure_reason":"-","route_name":"-","downstream_local_address":"127.0.0.1:8443","user_agent":"my-user-agent","response_code":"203","response_flags":"-","start_time":"2021-02-03T11:22:33.033Z","method":"PUT","request_id":"046CD4F5-7E0F-4EB1-B0C7-7BF404B08F10","upstream_host":"127.0.0.1:80","x_forwarded_for":"127.0.0.1","requested_server_name":"payment.example.com","bytes_received":"123","istio_policy_status":"-"}`
	line := fmt.Sprintf(`{"log":"%s\n","stream":"stdout","time":"2021-07-27T03:48:15.952919705Z"}`, strings.ReplaceAll(jsonLine, `"`, `\"`))

	got, err := parser.ParseString(line)
	require.NoError(t, err)

	want := map[string]string{
		"time_local":             "2021-02-03T11:22:33.033Z",
		"request_time":           "0.027",
		"request_length":         "123",
		"upstream_response_time": "0.026",
		"status":                 "203",
		"body_bytes_sent":        "37",
		"request":                "PUT /v1/order/2145 HTTP/1.1",
		"request_method":         "PUT",
		"upstream_cluster":       "payment-service",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("IstioParser.Parse(), got:\n%v\nwant\n%v", got, want)
	}
}

func BenchmarkParseIstio(b *testing.B) {
	parser := NewIstioParser()
	jsonLine := `{"bytes_sent":"37","upstream_cluster":"outbound|80||payment-service","downstream_remote_address":"127.0.0.1:40421","authority":"payment.example.com","path":"/v1/order/2145","protocol":"HTTP/1.1","upstream_service_time":"26","upstream_local_address":"127.0.0.1:53232","duration":"27","upstream_transport_failure_reason":"-","route_name":"-","downstream_local_address":"127.0.0.1:8443","user_agent":"my-user-agent","response_code":"203","response_flags":"-","start_time":"2021-02-03T11:22:33.033Z","method":"PUT","request_id":"046CD4F5-7E0F-4EB1-B0C7-7BF404B08F10","upstream_host":"127.0.0.1:80","x_forwarded_for":"127.0.0.1","requested_server_name":"payment.example.com","bytes_received":"123","istio_policy_status":"-"}`
	line := fmt.Sprintf(`{"log":"%s\n","stream":"stdout","time":"2021-07-27T03:48:15.952919705Z"}`, strings.ReplaceAll(jsonLine, `"`, `\"`))

	for i := 0; i < b.N; i++ {
		res, err := parser.ParseString(line)
		if err != nil {
			b.Error(err)
		}
		_ = fmt.Sprintf("%v", res)
	}
}
