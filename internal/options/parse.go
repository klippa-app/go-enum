package options

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"
)

var optionsRegex = regexp.MustCompile(`//enum:([\w,]+)`)

// Parse extracts go-enum options from an inline comment group.
func Parse(name string, cgroup *ast.CommentGroup, enumDefault *string) []string {
	if cgroup == nil {
		return []string{}
	}

	// extract the comma seperated list after the enum directive
	directive := optionsRegex.FindStringSubmatch(cgroup.List[0].Text)[1]
	options := strings.Split(directive, ",")

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
