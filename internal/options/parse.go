package options

import (
	"fmt"
	"go/ast"
	"strings"
)

func Parse(name string, cgroup *ast.CommentGroup, enumDefault *string) []string {
	if cgroup == nil {
		return []string{}
	}

	// "// enum:default,invalid some other text   "
	comment := strings.TrimPrefix(cgroup.List[0].Text, "//")
	// " enum:default,invalid some other text   "
	comment = strings.TrimSpace(comment)
	// "enum:default,invalid some other text"

	cmd := strings.Split(comment, " ")[0]
	// "enum:default,invalid"
	if !strings.HasPrefix(cmd, "enum:") {
		return []string{}
	}

	cmd = strings.TrimPrefix(cmd, "enum:")
	// "default,invalid"
	options := strings.Split(strings.ReplaceAll(cmd, " ", ""), ",")
	// [default, invalid]

	for i := range options {
		option := Option(options[i])
		if !option.isValid() {
			panic(fmt.Sprintf("unknown option: '%s'\n", option))
		}

		if name != "" && Option(option) == DefaultOption {
			if *enumDefault != "" {
				panic(fmt.Sprintf("Multiple defaults defined: %s, %s\n", *enumDefault, name))
			}
			*enumDefault = name
		}
	}

	return options
}
