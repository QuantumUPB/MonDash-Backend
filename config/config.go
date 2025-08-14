package config

import (
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

type Coordinates struct {
	Lat  float64 `yaml:"lat"`
	Long float64 `yaml:"long"`
}

type Config struct {
	Names       []string                         `yaml:"names"`
	URLs        map[string]string                `yaml:"urls"`
	Consumers   []string                         `yaml:"consumers"`
	Geolocation map[string]Coordinates           `yaml:"geolocation"`
	Links       [][]string                       `yaml:"links"`
	Paths       map[string]map[string][][]string `yaml:"paths"`
	// Additional fields are ignored
}

func Load(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func LoadFromEnv() (Config, error) {
	path := os.Getenv("CONFIG_FILE")
	if path == "" {
		path = "config.yaml"
	}
	return Load(path)
}

// NodesByConsumer returns a map of consumer names to the nodes that reference
// them in the Paths configuration. If no Paths are defined, an empty map is
// returned.
func (c Config) NodesByConsumer() map[string][]string {
	m := make(map[string][]string)
	seen := make(map[string]map[string]struct{})
	for node, consumerMap := range c.Paths {
		base := baseName(node)
		for cons := range consumerMap {
			if seen[cons] == nil {
				seen[cons] = make(map[string]struct{})
			}
			if _, ok := seen[cons][base]; !ok {
				seen[cons][base] = struct{}{}
				m[cons] = append(m[cons], base)
			}
		}
	}
	for cons := range m {
		sort.Strings(m[cons])
	}
	return m
}

// ConsumersByNode returns a map of node names to the consumers that can use
// that node based on the Paths configuration. If no Paths are defined, an
// empty map is returned.
func (c Config) ConsumersByNode() map[string][]string {
	m := make(map[string][]string)
	for node, consumerMap := range c.Paths {
		base := baseName(node)
		for cons := range consumerMap {
			m[base] = append(m[base], cons)
		}
	}
	for node := range m {
		sort.Strings(m[node])
	}
	return m
}

func baseName(s string) string {
	if len(s) == 0 {
		return s
	}
	r := s[len(s)-1]
	if r >= 'A' && r <= 'Z' {
		return s[:len(s)-1]
	}
	return s
}
