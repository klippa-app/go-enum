package main

import (
	"fmt"
	"go/token"
	"log"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/klippa-app/go-enum/coerce"
	"github.com/klippa-app/go-enum/internal/config"
	"github.com/klippa-app/go-enum/internal/values"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("go-enum: ")

	cfg := config.Instance()

	// determine current directory
	dir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	// parse enum file into AST
	fset := token.NewFileSet()
	pkgs, err := packages.Load(&packages.Config{
		Fset: fset,
		Mode: packages.NeedSyntax | packages.NeedName | packages.NeedModule | packages.NeedTypes | packages.NeedTypesInfo,
	}, fmt.Sprintf("file=%s.go", cfg.FileName))
	if err != nil {
		panic(err)
	}

	// extract package and type Info
	packageName := pkgs[0].Name
	packagePath := pkgs[0].PkgPath
	typeInfo := pkgs[0].TypesInfo

	// Extract the enum values, underlying types, and default from AST
	enumValues, underlyingType, enumDefault := values.ExtractEnumValues(typeInfo, fmt.Sprint(packagePath, ".", cfg.EnumName))
	if len(enumValues) == 0 {
		panic("no enum values found")
	}

	if underlyingType == "" {
		panic("could not determine underlying type for enum")
	}

	// Package context for templates.
	data := TemplateData{
		Pkg:              packageName,
		PkgPath:          packagePath,
		EnumName:         cfg.EnumName,
		BaseType:         underlyingType,
		EnumValues:       enumValues,
		EnumDefaultValue: enumDefault,
	}

	// alias for brevity
	execTemplate := func(name string, extension string) {
		ExecuteTemplate(name, data, fullPath(dir, cfg.FileName, cfg.EnumName, extension))
	}

	// execute all enabled templates.
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

// fullPath joins the `dir`, `fullPath`, `enumName`, and `suffix` to produce an absolute path for a (new) generated go-enum file.
func fullPath(dir string, fileName string, enumName string, suffix string) string {
	snakeEnum := coerce.SnakeCase(enumName)

	newFileNameParts := []string{fileName}

	// If the fileName does not match the enum name, we append the enum name to avoid conflicts.
	if fileName != snakeEnum {
		newFileNameParts = append(newFileNameParts, snakeEnum)
	}

	// Join all the name parts with '_' (might produce `path_.go`)
	newFileName := strings.Join(append(newFileNameParts, "enum", suffix), "_")
	newFileName = strings.Replace(newFileName, "_.", ".", 1)

	return path.Join(dir, newFileName)
}
