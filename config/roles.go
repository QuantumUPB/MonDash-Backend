package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Roles maps role names to the permissions granted for that role.
type Roles struct {
	Roles map[string][]string `yaml:"roles"`
}

// LoadRoles reads a YAML file and unmarshals it into a Roles struct.
func LoadRoles(path string) (Roles, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Roles{}, err
	}
	var r Roles
	if err := yaml.Unmarshal(b, &r); err != nil {
		return Roles{}, err
	}
	return r, nil
}

// LoadRolesFromEnv loads roles configuration from the ROLES_FILE environment
// variable. If ROLES_FILE is unset it defaults to roles.yaml in the working
// directory.
func LoadRolesFromEnv() (Roles, error) {
	path := os.Getenv("ROLES_FILE")
	if path == "" {
		path = "roles.yaml"
	}
	return LoadRoles(path)
}
