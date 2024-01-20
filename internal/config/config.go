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

func Instance() *Config {
	return config
}

func loadConfig() *Config {
	config := &Config{
		Verbose:      false,
		StringerCase: "snake",
	}

	config.FileName = strings.TrimSuffix(os.Getenv("GOFILE"), ".go")
	config.EnumName = coerce.PascalCase(config.FileName)
	config.Prefix = config.EnumName

	overrideWithFlags(config)

	return config
}

func overrideWithFlags(config *Config) {
	bindBool("v", &config.Verbose, "enable verbose logging")
	bindString("case", &config.StringerCase, "camel, pascal, snake, upper_snake, kebab, upper_kebab")

	bindString("name", &config.EnumName, "the name of the enum (defaults to the name of the file)")
	config.Prefix = config.EnumName // Prefix should default to whatever EnumName is generated as, or set to.
	bindString("prefix", &config.Prefix, "the prefix of the enum to strip (defaults to the name of the enum)")

	bindString("gql", &config.Generate.Gql, "'go': only generate marshaller, 'gql' only generate gql enum, 'full' generate both the marshaller and enum")
	bindBool("bson", &config.Generate.Bson, "generate functions for Bson")
	bindBool("json", &config.Generate.Json, "generate functions for Json")
	bindBool("xml", &config.Generate.Xml, "generate functions for Xml")
	bindBool("sql", &config.Generate.Sql, "generate functions for sql")
	bindBool("ent", &config.Generate.Ent, "generate functions for ent")
	bindBool("text", &config.Generate.Text, "generate functions for text")
	flag.Parse()
}

func bindString(name string, dest *string, usage string) {
	flag.CommandLine.StringVar(dest, name, *dest, usage)
}

func bindBool(name string, dest *bool, usage string) {
	flag.CommandLine.BoolVar(dest, name, *dest, usage)
}
