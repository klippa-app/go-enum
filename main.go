package main

import (
	"embed"
	"fmt"
	"go/token"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gertd/go-pluralize"
	"golang.org/x/tools/go/packages"

	"github.com/klippa-app/go-enum/coerce"
	"github.com/klippa-app/go-enum/internal/config"
	"github.com/klippa-app/go-enum/internal/util"
	"github.com/klippa-app/go-enum/internal/values"
)

var (
	//go:embed templates/*
	templates embed.FS
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("go-enum: ")

	cfg := config.Instance()

	dir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	pkgs, err := packages.Load(&packages.Config{
		Fset: fset,
		Mode: packages.NeedSyntax | packages.NeedName | packages.NeedModule | packages.NeedTypes | packages.NeedTypesInfo,
	}, fmt.Sprintf("file=%s.go", cfg.FileName))
	if err != nil {
		panic(err)
	}

	packageName := pkgs[0].Name
	packagePath := pkgs[0].PkgPath

	typeInfo := pkgs[0].TypesInfo

	enumValues, underlyingType, enumDefault := values.ExtractEnumValues(typeInfo, fmt.Sprint(packagePath, ".", cfg.EnumName))
	if len(enumValues) == 0 {
		panic("no enum values found")
	}

	if underlyingType == "" {
		panic("could not determine underlying type for enum")
	}

	visited := make(map[string]bool, 0)
	for i := range enumValues{
		if visited[enumValues[i].Value] == true{
			panic("no duplicate values allowed in enum")
		} else {
			visited[enumValues[i].Value] = true
		}
	}

	templates, err := template.New("").
		Funcs(TemplateFunctions). // Custom functions
		ParseFS(templates, "templates/*.tmpl")
	if err != nil {
		panic(err)
	}

	data := TemplateData{
		Pkg:              packageName,
		PkgPath:          packagePath,
		EnumName:         cfg.EnumName,
		BaseType:         underlyingType,
		EnumValues:       enumValues,
		EnumDefaultValue: enumDefault,
		Config:           cfg,
	}

	execTemplate := func(name string, extension string) {
		ExecuteTemplate(templates, name, fullPath(dir, cfg.FileName, cfg.EnumName, extension), data)
	}

	execTemplate("enum.tmpl", ".go")
	if cfg.Generate.Bson {
		execTemplate("bson.tmpl", "marshal_bson.go")
	}
	if cfg.Generate.Json {
		execTemplate("json.tmpl", "marshal_json.go")
	}
	if cfg.Generate.Xml {
		execTemplate("xml.tmpl", "marshal_xml.go")
	}
	if cfg.Generate.Sql || cfg.Generate.Ent {
		execTemplate("sql.tmpl", "marshal_sql.go")
	}
	if cfg.Generate.Text {
		execTemplate("text.tmpl", "marshal_text.go")
	}
	if cfg.Generate.Ent {
		execTemplate("ent.tmpl", "marshal_ent.go")
	}
	switch cfg.Generate.Gql {
	case "go":
		execTemplate("gql.go.tmpl", "marshal_gql.go")
	case "gql":
		execTemplate("gql.graphql.tmpl", ".graphql")
	case "full":
		execTemplate("gql.go.tmpl", "marshal_gql.go")
		execTemplate("gql.graphql.tmpl", ".graphql")
	}
}

func fullPath(dir string, fileName string, enumName string, suffix string) string {
	filePathBaseParts := []string{coerce.CamelCase(fileName)}
	if coerce.CamelCase(fileName) != coerce.CamelCase(enumName) {
		filePathBaseParts = append(filePathBaseParts, coerce.CamelCase(enumName))
	}

	suf := fmt.Sprint(strings.Join(filePathBaseParts, "_"), "Enum", coerce.PascalCase(suffix))

	return path.Join(dir, coerce.SnakeCase(suf))
}

func ExecuteTemplate(tmpl *template.Template, name string, path string, data TemplateData) {
	writer, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer writer.Close()

	err = tmpl.ExecuteTemplate(writer, name, data)
	if err != nil {
		panic(err)
	}
}

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

func receiver(s string) string {
	return fmt.Sprintf("%s_enum", strings.ToLower(s))
}

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
	Config           *config.Config
}
