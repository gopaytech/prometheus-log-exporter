package kubeparser

import (
	"fmt"

	"github.com/satyrius/gonx"
	"github.com/tidwall/gjson"
)

// KubeParser parses variables patterns using config.NamespaceConfig.Format.
type KubeParser struct {
	parser *gonx.Parser
}

// NewKubeParser returns a new text parser.
func NewKubeParser(format string) *KubeParser {
	return &KubeParser{
		parser: gonx.NewParser(format),
	}
}

// ParseString implements the Parser interface.
func (t *KubeParser) ParseString(line string) (map[string]string, error) {
	actualLogLine := gjson.Get(line, "log").String()
	entry, err := t.parser.ParseString(actualLogLine)
	if err != nil {
		return nil, fmt.Errorf("text log parsing err: %w", err)
	}

	return entry.Fields(), nil
}
