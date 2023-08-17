package main

import (
	"embed"
	"flag"
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
	"github.com/klippa-app/go-enum/internal/util"
	"github.com/klippa-app/go-enum/internal/values"
)

var (
	//go:embed templates/*
	templates embed.FS
)
var (
	packagePath string
	fileName    string
)

var (
	verbose      *bool
	stringerCase *string
	enumName     *string
	prefix       *string

	generateGql  *string
	generateBson *bool
	generateJson *bool
	generateXml  *bool
	generateSql  *bool
	generateEnt  *bool
)

func init() {
	fileName = strings.TrimSuffix(os.Getenv("GOFILE"), ".go")
	enumName = flag.String(
		"name",
		coerce.PascalCase(fileName),
		"the name of the enum (defaults to the name of the file)",
	)
	prefix = flag.String("prefix", *enumName, "the prefix of the enum to strip (defaults to the name of the enum)")

	verbose = flag.Bool("v", false, "enable verbose logging")
	stringerCase = flag.String("case", "snake", "camel, pascal, snake, upper_snake, kebab, upper_kebab")
	generateGql = flag.String("gql", "none", "'go': only generate marshaller, 'gql': only generate gql enum, 'full' generate both the marshaller and enum")
	generateBson = flag.Bool("bson", false, "generate functions for Bson")
	generateJson = flag.Bool("json", false, "generate functions for Json")
	generateXml = flag.Bool("xml", false, "generate functions for Xml")
	generateSql = flag.Bool("sql", false, "generate functions for sql")
	generateEnt = flag.Bool("ent", false, "generate functions for ent")
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of klippa/go-enum:\n")
	fmt.Fprintf(os.Stderr, "TODO\n")
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("go-enum: ")

	flag.Usage = Usage
	flag.Parse()

	dir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	pkgs, err := packages.Load(&packages.Config{
		Fset: fset,
		Mode: packages.NeedSyntax | packages.NeedName | packages.NeedModule | packages.NeedTypes | packages.NeedTypesInfo,
	}, fmt.Sprintf("file=%s.go", fileName))
	if err != nil {
		panic(err)
	}
	packageName := pkgs[0].Name
	packagePath := pkgs[0].PkgPath

	typeInfo := pkgs[0].TypesInfo

	enumValues, underlyingType, enumDefault := values.ExtractEnumValues(typeInfo, fmt.Sprint(packagePath, ".", *enumName))

	templates, err := template.New("").
		Funcs(TemplateFunctions). // Custom functions
		ParseFS(templates, "templates/*.tmpl")
	if err != nil {
		panic(err)
	}

	data := TemplateData{
		Pkg:              packageName,
		PkgPath:          packagePath,
		EnumName:         *enumName,
		BaseType:         underlyingType,
		EnumValues:       enumValues,
		EnumDefaultValue: enumDefault,
	}

	ExecuteTemplate(templates, "enum.tmpl", fullPath(dir, fileName, ".go"), data)
	if util.DereferenceOrNew(generateBson) {
		ExecuteTemplate(templates, "bson.tmpl", fullPath(dir, fileName, "marshal_bson.go"), data)
	}
	if util.DereferenceOrNew(generateJson) {
		ExecuteTemplate(templates, "json.tmpl", fullPath(dir, fileName, "marshal_json.go"), data)
	}
	if util.DereferenceOrNew(generateXml) {
		ExecuteTemplate(templates, "xml.tmpl", fullPath(dir, fileName, "marshal_xml.go"), data)
	}
	if util.DereferenceOrNew(generateSql) || util.DereferenceOrNew(generateEnt) {
		ExecuteTemplate(templates, "sql.tmpl", fullPath(dir, fileName, "marshal_sql.go"), data)
	}
	if util.DereferenceOrNew(generateEnt) {
		ExecuteTemplate(templates, "ent.tmpl", fullPath(dir, fileName, "marshal_ent.go"), data)
	}
	switch util.DereferenceOrNew(generateGql) {
	case "go":
		ExecuteTemplate(templates, "gql.go.tmpl", fullPath(dir, fileName, "marshal_gql.go"), data)
	case "gql":
		ExecuteTemplate(templates, "gql.graphql.tmpl", fullPath(dir, fileName, ".graphql"), data)
	case "full":
		ExecuteTemplate(templates, "gql.go.tmpl", fullPath(dir, fileName, "marshal_gql.go"), data)
		ExecuteTemplate(templates, "gql.graphql.tmpl", fullPath(dir, fileName, ".graphql"), data)
	}
}

func fullPath(dir string, fileName string, suffix string) string {
	suf := fmt.Sprint(coerce.CamelCase(fileName), "Enum", coerce.PascalCase(suffix))

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
	if stringerCase == nil {
		panic("no stringerCase set. (how did you do that?)")
	}

	s = strings.TrimPrefix(coerce.SnakeCase(s), fmt.Sprint(coerce.SnakeCase(*prefix), "_"))

	switch *stringerCase {
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

	panic(fmt.Sprintf("unknown stringerCase: %s", *stringerCase))
}

func stringerFn() string {
	if stringerCase == nil {
		panic("no stringerCase set. (how did you do that?)")
	}

	switch *stringerCase {
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

	panic(fmt.Sprintf("unknown stringerCase: %s", *stringerCase))
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
}
