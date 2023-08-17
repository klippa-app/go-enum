package coerce

import (
	"github.com/stoewer/go-strcase"
)

func CamelCase(str string) string {
	return strcase.LowerCamelCase(str)
}

func PascalCase(str string) string {
	return strcase.UpperCamelCase(str)
}

func SnakeCase(str string) string {
	return strcase.SnakeCase(str)
}

func UpperSnakeCase(str string) string {
	return strcase.UpperSnakeCase(str)
}

func KebabCase(str string) string {
	return strcase.KebabCase(str)
}

func UpperKebabCase(str string) string {
	return strcase.UpperKebabCase(str)
}
