package kubeparser

import (
	"fmt"
	"strings"

	"github.com/satyrius/gonx"
	"github.com/tidwall/gjson"
)

// KubeParser parses variables patterns using config.NamespaceConfig.Format.
type KubeParser struct {
	parser       *gonx.Parser
	fnGetKubeLog func(string) string
}

// NewKubeParser returns a new text parser.
func NewKubeParser(format string) *KubeParser {
	return &KubeParser{
		parser: gonx.NewParser(format),
	}
}

func NewKubeCRIParser(format string) *KubeParser {
	return &KubeParser{
		parser: gonx.NewParser(format),
		fnGetKubeLog: func(s string) string {
			// timestamp stdout [FP] actualLog
			line := strings.SplitN(s, " ", 4)[3]
			line = gjson.Get(line, "log").String()
			return strings.TrimSpace(line)
		},
	}
}

// kubernetes wrap the actual log line, need to unwrap it to get the actual
func (p *KubeParser) getKubeLogLine(line string) string {
	if p.fnGetKubeLog == nil {
		p.fnGetKubeLog = func(s string) string {
			line := gjson.Get(s, "log").String()
			return strings.TrimSpace(line)
		}
	}

	return p.fnGetKubeLog(line)
}

// ParseString implements the Parser interface.
func (p *KubeParser) ParseString(line string) (map[string]string, error) {
	actualLogLine := p.getKubeLogLine(line)

	entry, err := p.parser.ParseString(actualLogLine)
	if err != nil {
		return nil, fmt.Errorf("text log parsing err: %w", err)
	}

	return entry.Fields(), nil
}
