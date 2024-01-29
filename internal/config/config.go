package config

import (
	"flag"
	"os"
	"strings"

	"github.com/klippa-app/go-enum/coerce"
)

type Config struct {
	FileName    string
	PackagePath string
	EnumName    string
	Prefix      string

	Verbose      bool
	StringerCase string

	Generate struct {
		Gql  string
		Bson bool
		Json bool
		Xml  bool
		Sql  bool
		Ent  bool
		Text bool
	}
}

var config *Config

func init() {
	config = loadConfig()
}

// Instance returns the singleton instance of the loaded config.
func Instance() *Config {
	return config
}

// loadConfig initializes a default config and merges in the configuration from the CLI flags.
// TODO: also merge in a project configuration file.
func loadConfig() *Config {
	config := &Config{
		Verbose:      false,
		StringerCase: "snake",
		FileName:     strings.TrimSuffix(os.Getenv("GOFILE"), ".go"),
	}

	config.EnumName = coerce.PascalCase(config.FileName)

	overrideWithFlags(config)

	return config
}

// overrideWithFlags overrides the config with the settings provided by CLI flags.
func overrideWithFlags(config *Config) {
	bindBool("v", &config.Verbose, "enable verbose logging")
	bindString("case", &config.StringerCase, "camel, pascal, snake, upper_snake, kebab, upper_kebab")
	bindString("prefix", &config.Prefix, "the prefix of the enum to strip (defaults to the name of the enum)")
	bindString("name", &config.EnumName, "the name of the enum (defaults to the name of the file)")
	bindString("gql", &config.Generate.Gql, "'go': only generate marshaller, 'gql' only generate gql enum, 'full' generate both the marshaller and enum")
	bindBool("bson", &config.Generate.Bson, "generate functions for Bson")
	bindBool("json", &config.Generate.Json, "generate functions for Json")
	bindBool("xml", &config.Generate.Xml, "generate functions for Xml")
	bindBool("sql", &config.Generate.Sql, "generate functions for sql")
	bindBool("ent", &config.Generate.Ent, "generate functions for ent")
	bindBool("text", &config.Generate.Text, "generate functions for text")
	flag.Parse()
}

// bindString is a helper for binding the flag name, to the dest string, and using the initial value of dest as the default value.
func bindString(name string, dest *string, usage string) {
	flag.CommandLine.StringVar(dest, name, *dest, usage)
}

// bindString is a helper for binding the flag name, to the dest bool, and using the initial value of dest as the default value.
func bindBool(name string, dest *bool, usage string) {
	flag.CommandLine.BoolVar(dest, name, *dest, usage)
}
