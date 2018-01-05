package config

import (
	"sort"
	"errors"
)

// NamespaceConfig is a struct describing single metric namespaces
type NamespaceConfig struct {
	Name           string            `hcl:",key"`
	SourceFiles    []string          `hcl:"source_files" yaml:"source_files"`
	Format         string            `hcl:"format"`
	Labels         map[string]string `hcl:"labels"`
	RelabelConfigs []RelabelConfig   `hcl:"relabel" yaml:"relabel_configs"`

	OrderedLabelNames  []string
	OrderedLabelValues []string
}

func (c *NamespaceConfig) StabilityWarnings() error {
	if len(c.RelabelConfigs) > 0 {
		return errors.New("you are using the 'relabel' configuration parameter")
	}

	return nil
}

func (c *NamespaceConfig) MustCompile() {
	err := c.Compile()
	if err != nil {
		panic(err)
	}
}

func (c *NamespaceConfig) Compile() error {
	for i := range c.RelabelConfigs {
		if err := c.RelabelConfigs[i].Compile(); err != nil {
			return nil
		}
	}

	c.OrderLabels()

	return nil
}

// OrderLabels builds two lists of label keys and values, ordered by label name
func (c *NamespaceConfig) OrderLabels() {
	keys := make([]string, 0, len(c.Labels))
	values := make([]string, len(c.Labels))

	for k := range c.Labels {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for i, k := range keys {
		values[i] = c.Labels[k]
	}

	c.OrderedLabelNames = keys
	c.OrderedLabelValues = values
}