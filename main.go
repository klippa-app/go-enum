package main

import (
	"bufio"
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
	stringerCase = flag.String("case", "snake", "camel, pascal, snake, upper_snake, kebab, upper_kebab, upper, lower")
	generateGql = flag.String("gql", "none", "'go': only generate marshaller, 'gql': only generate gql enum, 'full' generate both the marshaller and enum")
	generateBson = flag.Bool("bson", false, "generate functions for Bson")
	generateJson = flag.Bool("json", false, "generate functions for Json")
	generateXml = flag.Bool("xml", false, "generate functions for Xml")
	generateSql = flag.Bool("sql", false, "generate functions for sql")
	generateEnt = flag.Bool("ent", false, "generate functions for ent")
}

var prints = []string{
	`Add a go generate line with in your go code as followed:
	
	//go:generate go run --mod=mod github.com/klippa-app/go-enum`,
	`	To this line add any of the defined flags:`,
	`	-v						enables verbose logging
							Verbose logging prints additional logging used for debugging`,
	`	-case=[camel|pascal|snake|upper_snake
		|kebab|upper_kebab|upper|lower]
							case changes the casing for the stringer
							The stringer is used to convert between external systems`,

	`	-gpl=[none|full]				generate gql enum values`,
	`	-bson						generate functions for Bson`,
	`	-json						generate functions for Json`,
	`	-xml						generate functions for Xml`,
	`	-sql						generate functions for sql`,
	`	-ent						generate functions for ent`,
	`Change how your enums behave
	Go-enums supports flags per value to modify your enums behavior
	Behind your values add a comment as followed
	//enum:
	after the colon add your flags with a comma separator`,
	`	default						default set the default value of the enum to this value
							The default value is used when no other enum value is valid`,
	`	invalid						invalid set a specific value as invalid
							An invalid enum value can be used in code but can not be marshalled and is not in the valid list of enum values`,
	`	Example
	const Unknown Day = 0 //enum:default,invalid`,

	`
You are ready to generate your code:
run go generate`,
}

func Usage() {
	fmt.Fprintln(os.Stdout, `Usage of klippa/go-enum:

	go-enum is a code generator for Golang.
	Golang by default does not support enums
	Golang defines enums using uints and custom marshallers however if you have a lot of these there can be a lot of boilerplate code
	go-enum generate custom marshallers for different marshallers next to your existing code

	Make sure to read the README and checkout the examples

	press "enter" to continue to implementation`)

	s := bufio.NewScanner(os.Stdin)
	for i := 0; s.Scan(); i++ {
		if i < len(prints) {
			fmt.Fprintln(os.Stdout, prints[i])
		}
	}
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
	case "upper":
		return coerce.Upper(s)
	case "lower":
		return coerce.Lower(s)
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
	case "upper":
		return "coerce.Upper"
	case "lower":
		return "coerce.Lower"
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
	"upper":          strings.ToUpper,
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
