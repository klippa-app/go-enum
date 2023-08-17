package values

import (
	"go/ast"
	"go/types"

	"github.com/klippa-app/go-enum/internal/options"
	"github.com/klippa-app/go-enum/internal/util"
)

type EnumValue struct {
	Name    string
	Type    string
	Options []string
	Value   string
}

func ExtractEnumValues(typeInfo *types.Info, enumType string) (enums []EnumValue, underlyingType string, enumDefault string) {
	for scope := range typeInfo.Scopes {
		file, ok := scope.(*ast.File)
		if !ok {
			continue
		}

		genDecls := util.Only[*ast.GenDecl](file.Decls)
		for i := range genDecls {
			valueSpecs := util.Only[*ast.ValueSpec](genDecls[i].Specs)
			for j := range valueSpecs {
				value := valueSpecs[j]
				object, ok := typeInfo.ObjectOf(value.Names[0]).(*types.Const) // TODO: only const?
				if !ok || object.Type().String() != enumType {
					continue
				}

				if underlyingType == "" {
					underlyingType = object.Type().Underlying().String()
				} else if underlyingType != object.Type().Underlying().String() {
					panic("differing underlying types for enum")
				}

				enums = append(enums, EnumValue{
					Name:    object.Name(),
					Type:    underlyingType,
					Options: options.Parse(object.Name(), value.Comment, &enumDefault),
					Value:   object.Val().ExactString(),
				})
			}
		}
	}
	return
}
