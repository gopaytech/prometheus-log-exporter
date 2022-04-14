package parser

import (
	"github.com/martin-helmich/prometheus-nginxlog-exporter/config"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/parser/haproxyjsonparser"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/parser/istioparser"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/parser/jsonparser"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/parser/kubeparser"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/parser/textparser"
)

// Parser parses a line of log to a map[string]string.
type Parser interface {
	ParseString(line string) (map[string]string, error)
}

// NewParser returns a Parser with the given config.NamespaceConfig.
func NewParser(nsCfg config.NamespaceConfig) Parser {
	switch nsCfg.Parser {
	case "text":
		return textparser.NewTextParser(nsCfg.Format)
	case "haproxy-json":
		return haproxyjsonparser.NewHaproxyJsonParser()
	case "json":
		return jsonparser.NewJsonParser()
	case "kube":
		return kubeparser.NewKubeParser(nsCfg.Format)
	case "kube-cri":
		return kubeparser.NewKubeCRIParser(nsCfg.Format)
	case "istio":
		return istioparser.NewIstioParser()
	case "istio-cri":
		return istioparser.NewIstioCRIParser()
	default:
		return textparser.NewTextParser(nsCfg.Format)
	}
}
