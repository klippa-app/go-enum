package main

import (
	"embed"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/gertd/go-pluralize"
	"github.com/klippa-app/go-enum/coerce"
	"github.com/klippa-app/go-enum/internal/config"
	"github.com/klippa-app/go-enum/internal/util"
	"github.com/klippa-app/go-enum/internal/values"
)

var (
	//go:embed templates/*
	templatesFS embed.FS
)

// The parsed templates from the file system.
var templates *template.Template

func init() {
	var err error

	// Initialise marshaller templates
	templates, err = template.New("").
		Funcs(TemplateFunctions).
		ParseFS(templatesFS, "templates/*.tmpl")
	if err != nil {
		panic(err)
	}
}

// ExecuteTemplate executes the template with the name `name`, with the data `data`, and writes the output to `path`.
func ExecuteTemplate(name string, data TemplateData, path string) {
	writer, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer writer.Close()

	err = templates.ExecuteTemplate(writer, name, data)
	if err != nil {
		panic(err)
	}
}

// TemplateData is the data provided to a template.
type TemplateData struct {
	Pkg              string
	PkgPath          string
	EnumName         string
	BaseType         string
	EnumDefaultValue string
	Gqlgen           bool
	Bson             bool
	Json             bool
	Xml              bool
	EnumValues       []values.EnumValue
}

// TemplateFunctions are the functions that are available for templates to use.
var TemplateFunctions = template.FuncMap{
	"containsString": util.Contains[string],
	"lower":          strings.ToLower,
	"camel":          coerce.CamelCase,
	"pascal":         coerce.PascalCase,
	"upperSnake":     coerce.UpperSnakeCase,
	"plural":         pluralize.NewClient().Plural,
	"stringer":       stringer,
	"stringerFn":     stringerFn,
	"receiver":       receiver,
}

// stringer coerces a string to match the case configured for the Stringer interface.
func stringer(s string) string {
	cfg := config.Instance()

	s = strings.TrimPrefix(coerce.SnakeCase(s), fmt.Sprint(coerce.SnakeCase(cfg.Prefix), "_"))

	switch cfg.StringerCase {
	case "camel":
		return coerce.CamelCase(s)
	case "pascal":
		return coerce.PascalCase(s)
	case "snake":
		return coerce.SnakeCase(s)
	case "upper_snake":
		return coerce.UpperSnakeCase(s)
	case "kebab":
		return coerce.KebabCase(s)
	case "upper_kebab":
		return coerce.UpperKebabCase(s)
	}

	panic(fmt.Sprintf("unknown stringerCase: %s", cfg.StringerCase))
}

// stringerFn returns the function used to coerce a string by the stringer.
func stringerFn() string {
	cfg := config.Instance()

	switch cfg.StringerCase {
	case "camel":
		return "coerce.CamelCase"
	case "pascal":
		return "coerce.PascalCase"
	case "snake":
		return "coerce.SnakeCase"
	case "upper_snake":
		return "coerce.UpperSnakeCase"
	case "kebab":
		return "coerce.KebabCase"
	case "upper_kebab":
		return "coerce.UpperKebabCase"
	}

	panic(fmt.Sprintf("unknown stringerCase: %s", cfg.StringerCase))
}

// receiver returns a safe variable name to use for the method recievers.
func receiver(s string) string {
	return fmt.Sprintf("%s_enum", strings.ToLower(s))
}
